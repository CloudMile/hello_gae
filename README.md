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

## Snapshot GCE Instance Disks with GAE Cron Job
using GAE cron job and queue to call the GCP apis to create snapshot for all of instances.

## Deploy
```
$ gcloud app deploy app.yaml cron.yaml queue.yaml
```
or
```
$ gcloud app deploy app.yaml
$ gcloud app deploy cron.yaml
$ gcloud app deploy queue.yaml
```
