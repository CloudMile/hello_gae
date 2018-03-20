# hello_GAE

## Clone This Project
- Part 1 HelloWorld
  - git checkout origin hello-world
- Part 2 Service and Version
  - git checkout origin services_and_versions
- Part 3 Custom Domail and HTTPs
  - git checkout origin custom_domain_and_ssl
- Part 4 Datastore and Memcache
  - git checkout origin datastore_memcache
- Part 5 Queue and Cron Job
  - git checkout origin queue_and_cronjob

## Deploy to GAE
install [gcloud SDK](https://cloud.google.com/sdk/downloads) before
```
$ gcloud app deploy
```
### deploy with version
```
$ gcloud app deploy -v <YOUR_VERSION>
```

## For Part3
### deploy dispatch
edit your Domain into the file before you deploy
```
$ gcloud app deploy dispatch.yaml
```

## For Part4
### Set up Datastore
```
$ vim app.go
```

```
const (
	GAEEntityKind      = "<DATASTORE_KIND>"
	HelloKey           = "<SET_YOUR_KEY>"
	IsUsingMemcache    = false
	MemcacheExprieTime = time.Duration(3600) * time.Second
)

```
change <DATASTORE_KIND> and <SET_YOUR_KEY> to you want
if want to use GAE memcache, please enable `IsUsingMemcache` and also change `3600` to you need

## For Part5
### Snapshot GCE Instance Disks with GAE Cron Job
using GAE cron job and queue to call the GCP apis to create snapshot for all of instances.
