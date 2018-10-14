# Hello GAE

This is a project for _Getting Start for GAE_

## Clone This Project

```shell
git clone git@github.com:CloudMile/hello_gae.git
```

You can checkout the different branch to different chapter.

- Part 1 HelloWorld
  - `git checkout origin hello-world`
- Part 2 Service and Version
  - `git checkout origin services_and_versions`
- Part 3 Custom Domail and HTTPs
  - `git checkout origin custom_domain_and_ssl`
- Part 4 Datastore and Memcache
  - `git checkout origin datastore_memcache`
- Part 5 Queue and Cron Job
  - `git checkout origin queue_and_cronjob`

## Deploy to GAE

Install [gcloud SDK](https://cloud.google.com/sdk/downloads) first, and login.

Deploy project:

```shell
$ gcloud app deploy
```

### Deploy with version

```shell
$ gcloud app deploy -v <YOUR_VERSION>
```

## For Part3

### Deploy dispatch

Edit your __Domain__ in the file `dispatch.yaml` before you deploy

```shell
$ gcloud app deploy dispatch.yaml
```

## For Part4

### Set up Datastore

```shell
$ vim app.go
```

```shell
const (
	GAEEntityKind      = "<DATASTORE_KIND>"
	HelloKey           = "<SET_YOUR_KEY>"
	IsUsingMemcache    = false
	MemcacheExprieTime = time.Duration(3600) * time.Second
)

```

Replace the `<DATASTORE_KIND>` and `<SET_YOUR_KEY>`.

If you want to use GAE memcache, please enable `IsUsingMemcache` and also change `3600` to you need.

## For Part5

### Snapshot GCE Instance Disks with GAE Cron Job

Use the GAE cron job and queue to call the GCP apis to create snapshot for all of instances.
