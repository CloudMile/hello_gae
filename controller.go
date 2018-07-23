package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

func cronHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cron/snapshot" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "query: %v", r.URL.Query())

	t := taskqueue.NewPOSTTask("/cron/worker", r.URL.Query())
	if _, err := taskqueue.Add(ctx, t, "create-snapshot"); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func workHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cron/worker" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	filterParams := r.PostFormValue("filter")
	log.Infof(ctx, "filter: %v", filterParams)

	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, "https://www.googleapis.com/auth/compute"),
			Base: &urlfetch.Transport{
				Context: ctx,
			},
		},
	}
	computeService, err := compute.New(client)
	if err != nil {
		log.Errorf(ctx, "compute error: %s", err)
	}

	gceZoneList := getGCEZone(ctx, computeService)
	gceDiksMap := rangeDiskZone(ctx, computeService, gceZoneList, filterParams)
	rangeCreateSnapshot(ctx, computeService, gceDiksMap)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	switch status {
	case http.StatusNotFound:
		fmt.Fprint(w, "404 Not Found")
	case http.StatusMethodNotAllowed:
		fmt.Fprint(w, "405 Method Not Allow")
	default:
		fmt.Fprint(w, "Bad Request")
	}
}

func getGCEZone(ctx context.Context, computeService *compute.Service) (gceZoneList []string) {
	req := computeService.Zones.List(getProjectID(ctx))

	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			gceZoneList = append(gceZoneList, zone.Name)
		}
		return nil
	}); err != nil {
		log.Errorf(ctx, "compute.ZoneList error: %s", err)
	}
	return
}

func rangeDiskZone(ctx context.Context, computeService *compute.Service, gceZoneList []string, filterParams string) (gceDiksMap map[string][]string) {
	projectID := getProjectID(ctx)
	gceDiksMapTemp := make(map[string][]string)

	for _, zoneName := range gceZoneList {
		req := computeService.Disks.List(projectID, zoneName)
		if filterParams != "" {
			req = req.Filter(filterParams)
		}
		if err := req.Pages(ctx, func(page *compute.DiskList) error {
			for _, disk := range page.Items {
				if len(disk.Users) > 0 {
					gceDiksMapTemp[zoneName] = append(gceDiksMapTemp[zoneName], disk.Name)
				}
			}
			return nil
		}); err != nil {
			log.Errorf(ctx, "compute.DiskList error: %s", err)
		}
	}
	return gceDiksMapTemp
}

func rangeCreateSnapshot(ctx context.Context, computeService *compute.Service, gceDiksMap map[string][]string) {
	projectID := getProjectID(ctx)
	t := time.Now()
	tString := strings.ToLower(t.Format("06010215MST"))

	for zoneName, gceDiskList := range gceDiksMap {
		for _, diskName := range gceDiskList {
			rb := &compute.Snapshot{
				Name: diskName + `-` + tString,
			}
			log.Infof(ctx, "create snapshot: %s => %s", zoneName, diskName)
			_, err := computeService.Disks.CreateSnapshot(projectID, zoneName, diskName, rb).Context(ctx).Do()
			if err != nil {
				log.Errorf(ctx, "CreateSnapshot error: %s", err)
			}
		}
	}
}

func getProjectID(ctx context.Context) (projectID string) {
	if appengine.AppID(ctx) == "None" {
		projectID = os.Getenv("PROJECT_ID")
	} else {
		projectID = appengine.AppID(ctx)
	}
	return
}
