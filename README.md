# hello_GAE

## Clone This Project
- Part 1 HelloWorld
  - git checkout origin hello-world
- Part 2 Service and Version
  - git checkout origin services_and_versions
- Part 3 Custom Domail and HTTPs
  - git checkout origin custom_domain_and_ssl

## Deploy to GAE
install [gcloud SDK](https://cloud.google.com/sdk/downloads) before
```
$ gcloud app deploy
```
### deploy with version
```
$ gcloud app deploy -v <YOUR_VERSION>
```

### deploy dispatch
edit your Domain into the file before you deploy
```
$ gcloud app deploy dispatch.yaml
```
