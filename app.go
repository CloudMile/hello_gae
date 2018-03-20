package app

import (
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/cron/snapshot", cronHandle)
	http.HandleFunc("/cron/worker", workHandle)
	appengine.Main()
}
