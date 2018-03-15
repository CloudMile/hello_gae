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

## Set up Datastore
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

## Deploy to GAE
install [gcloud SDK](https://cloud.google.com/sdk/downloads) before
```
$ gcloud app deploy
```
### deploy with version
```
$ gcloud app deploy -v <YOUR_VERSION>
```
