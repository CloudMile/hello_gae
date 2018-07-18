package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	gceZoneList := getGCEZone(ctx, *client)
	rangeDiskZone(ctx, *client, gceZoneList, filterParams)
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

func getGCEZone(ctx context.Context, client http.Client) (gceZoneList GCEZoneList) {
	url := "https://www.googleapis.com/compute/v1/projects/" + getProjectID(ctx) + "/zones"
	log.Infof(ctx, "url: %v", url)
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf(ctx, "resp error: %s", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	gceZoneList = GCEZoneList{}
	json.Unmarshal(bodyBytes, &gceZoneList)
	return
}

func rangeDiskZone(ctx context.Context, client http.Client, gceZoneList GCEZoneList, filterParams string) {
	for _, zone := range gceZoneList.Items {
		url := ""
		if filterParams != "" {
			url = "https://www.googleapis.com/compute/v1/projects/" + getProjectID(ctx) + "/zones/" + zone.Name + "/disks?filter=" + filterParams
		} else {
			url = "https://www.googleapis.com/compute/v1/projects/" + getProjectID(ctx) + "/zones/" + zone.Name + "/disks"
		}

		log.Infof(ctx, "url: %v", url)
		resp, err := client.Get(url)
		if err != nil {
			log.Errorf(ctx, "resp error: %s", err)
		}
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		gceDisk := GCEDisk{}
		json.Unmarshal(bodyBytes, &gceDisk)

		rangeCreateSnapshot(ctx, client, zone.Name, gceDisk)
	}
}

func rangeCreateSnapshot(ctx context.Context, client http.Client, zone string, gceDisk GCEDisk) {
	t := time.Now()
	tString := strings.ToLower(t.Format("06010215MST"))

	for _, diskItem := range gceDisk.Items {
		log.Infof(ctx, "name: %v", diskItem.Name)
		if len(diskItem.Users) > 0 {
			createSnapshotURL := "https://www.googleapis.com/compute/v1/projects/" + getProjectID(ctx) + "/zones/" + zone + "/disks/" + diskItem.Name + "/createSnapshot"
			log.Infof(ctx, "createSnapshotURL: %v", createSnapshotURL)
			jsonString := `{"name":"` + diskItem.Name + "-" + tString + `"}`
			log.Infof(ctx, "jsonString: %v", jsonString)
			_, err := client.Post(createSnapshotURL, "application/json", bytes.NewBuffer([]byte(jsonString)))
			if err != nil {
				log.Errorf(ctx, "resp error: %s", err)
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
