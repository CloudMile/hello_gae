package app

import (
	"context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func getFromDatastore(ctx context.Context) (result string) {
	log.Infof(ctx, "get from datastore")
	oobj := OutputObject{}
	datastoreKey := datastore.NewKey(ctx, GAEEntityKind, HelloKey, 0, nil)

	if err := datastore.Get(ctx, datastoreKey, &oobj); err != nil {
		log.Errorf(ctx, "%v", err)
		result = ""
	}
	result = oobj.Key
	return
}

func setToDatastore(ctx context.Context, oobj OutputObject) {
	log.Infof(ctx, "set into datastore")
	datastoreKey := datastore.NewKey(ctx, GAEEntityKind, HelloKey, 0, nil)

	if _, err := datastore.Put(ctx, datastoreKey, &oobj); err != nil {
		log.Errorf(ctx, "%v", err)
		return
	}
}
