# stardust

[![actions-workflow-test][actions-workflow-test-badge]][actions-workflow-test]
[![docker-build][docker-build-badge]][docker-build]
[![pkg.go.dev][pkg.go.dev-badge]][pkg.go.dev]
[![release][release-badge]][release]
[![license][license-badge]][license]

Report a summary what GitHub repositories you recently have starred.

![screenshot](docs/assets/screenshot.png)

For example, you can get a report about what repositories you have starred in the week on every Sunday.

## Run locally

```
$ docker run -d -p 8080:8080 --env-file env micnncim/stardust:latest -local
$ curl localhost:8080
```

## Run on Cloud Run

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run/?git_repo=https://github.com/micnncim/stardust.git)

Using it with Cloud Scheduler and [Berglas](https://github.com/GoogleCloudPlatform/berglas) is recommended.

Before push the above button, you need to set up the `app.json`.
The example is [app.example.json](app.example.json).
The detailed document is [here](https://github.com/GoogleCloudPlatform/cloud-run-button#customizing-deployment-parameters).

### CLI

1. Deploy stardust to Cloud Run

```
$ PROJECT_ID=my-project
$ BUCKET_ID=my-bucket
$ KMS_KEY=projects/${PROJECT_ID}/locations/global/keyRings/berglas/cryptoKeys/berglas-key
$ berglas create ${BUCKET_ID}/github-token "<GITHUB_TOKEN>" --key $KMS_KEY
$ berglas create ${BUCKET_ID}/slack-token "<SLACK_TOKEN>" --key $KMS_KEY
$ berglas create ${BUCKET_ID}/slack-channel-id "<SLACK_CHANNEL_ID>" --key $KMS_KEY
$ SERVICE_ACCOUNT=my-service-account
$ gcloud iam service-accounts create $SERVICE_ACCOUNT --project $PROJECT_ID
$ SERVICE_ACCOUNT_EMAIL=${SERVICE_ACCOUNT}@${PROJECT_ID}.iam.gserviceaccount.com
$ berglas grant ${BUCKET_ID}/github-token --member serviceAccount:${SERVICE_ACCOUNT_EMAIL}
$ berglas grant ${BUCKET_ID}/slack-token --member serviceAccount:${SERVICE_ACCOUNT_EMAIL}
$ berglas grant ${BUCKET_ID}/slack-channel-id --member serviceAccount:${SERVICE_ACCOUNT_EMAIL}
$ SERVICE=stardust
$ USERNAME=micnncim
$ gcloud run deploy $SERVICE \
    --project $PROJECT_ID \
    --platform managed \
    --image micnncim/stardust:latest \
    --set-env-vars GITHUB_TOKEN=berglas://${BUCKET_ID}/github-token,GITHUB_USERNAME=${USERNAME},ENABLE_SLACK=true,SLACK_TOKEN=berglas://${BUCKET_ID}/slack-token,SLACK_CHANNEL_ID=berglas://${BUCKET_ID}/slack-channel-id,INTERVAL=168h \
    --service-account ${SERVICE_ACCOUNT_EMAIL}
```

2. Set up Cloud Scheduler

```
$ SERVICE_ACCOUNT=my-service-account
$ gcloud iam service-accounts create $SERVICE_ACCOUNT \
    --display-name "<DISPLAYED_SERVICE_ACCOUNT_NAME>"
$ SERVICE=stardust
$ gcloud run services add-iam-policy-binding $SERVICE \
    --member=serviceAccount:${SERVICE_ACCOUNT}@{PROJECT_ID}.iam.gserviceaccount.com \
    --role=roles/run.invoker
$ SERVICE_ACCOUNT_EMAIL=${SERVICE_ACCOUNT}@${PROJECT_ID}.iam.gserviceaccount.com
$ SERVICE_URL=$(gcloud run services describe $SERVICE --format 'value(status.url)')
$ JOB=my-job
$ CRON="0 9 * * 0"
$ gcloud beta scheduler jobs create http $JOB --schedule $CRON
    --http-method GET \
    --uri $SERVICE_URL \
    --oidc-service-account-email $SERVICE_ACCOUNT_EMAIL   \
    --oidc-token-audience $SERVICE_URL
```

## Report Platforms

The supported report platforms are below.

- Slack

## Reference

- [Running services on a schedule | Cloud Run Documentation | Google Cloud](https://cloud.google.com/run/docs/triggering/using-scheduler)
- [Berglas Cloud Run Example - Go](https://github.com/GoogleCloudPlatform/berglas/blob/master/examples/cloudrun/go/README.md)

<!-- badge links -->

[actions-workflow-test]: https://github.com/micnncim/stardust/actions?query=workflow%3ATest
[actions-workflow-test-badge]: https://img.shields.io/github/workflow/status/micnncim/stardust/Test?label=Test&style=for-the-badge&logo=github

[docker-build]: https://hub.docker.com/r/micnncim/stardust
[docker-build-badge]: https://img.shields.io/docker/cloud/build/micnncim/stardust?logo=docker&style=for-the-badge

[pkg.go.dev]: https://pkg.go.dev/github.com/micnncim/stardust?tab=overview
[pkg.go.dev-badge]: https://img.shields.io/badge/pkg.go.dev-reference-02ABD7?style=for-the-badge&logoWidth=25&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9Ijg1IDU1IDEyMCAxMjAiPjxwYXRoIGZpbGw9IiMwMEFERDgiIGQ9Ik00MC4yIDEwMS4xYy0uNCAwLS41LS4yLS4zLS41bDIuMS0yLjdjLjItLjMuNy0uNSAxLjEtLjVoMzUuN2MuNCAwIC41LjMuMy42bC0xLjcgMi42Yy0uMi4zLS43LjYtMSAuNmwtMzYuMi0uMXptLTE1LjEgOS4yYy0uNCAwLS41LS4yLS4zLS41bDIuMS0yLjdjLjItLjMuNy0uNSAxLjEtLjVoNDUuNmMuNCAwIC42LjMuNS42bC0uOCAyLjRjLS4xLjQtLjUuNi0uOS42bC00Ny4zLjF6bTI0LjIgOS4yYy0uNCAwLS41LS4zLS4zLS42bDEuNC0yLjVjLjItLjMuNi0uNiAxLS42aDIwYy40IDAgLjYuMy42LjdsLS4yIDIuNGMwIC40LS40LjctLjcuN2wtMjEuOC0uMXptMTAzLjgtMjAuMmMtNi4zIDEuNi0xMC42IDIuOC0xNi44IDQuNC0xLjUuNC0xLjYuNS0yLjktMS0xLjUtMS43LTIuNi0yLjgtNC43LTMuOC02LjMtMy4xLTEyLjQtMi4yLTE4LjEgMS41LTYuOCA0LjQtMTAuMyAxMC45LTEwLjIgMTkgLjEgOCA1LjYgMTQuNiAxMy41IDE1LjcgNi44LjkgMTIuNS0xLjUgMTctNi42LjktMS4xIDEuNy0yLjMgMi43LTMuN2gtMTkuM2MtMi4xIDAtMi42LTEuMy0xLjktMyAxLjMtMy4xIDMuNy04LjMgNS4xLTEwLjkuMy0uNiAxLTEuNiAyLjUtMS42aDM2LjRjLS4yIDIuNy0uMiA1LjQtLjYgOC4xLTEuMSA3LjItMy44IDEzLjgtOC4yIDE5LjYtNy4yIDkuNS0xNi42IDE1LjQtMjguNSAxNy05LjggMS4zLTE4LjktLjYtMjYuOS02LjYtNy40LTUuNi0xMS42LTEzLTEyLjctMjIuMi0xLjMtMTAuOSAxLjktMjAuNyA4LjUtMjkuMyA3LjEtOS4zIDE2LjUtMTUuMiAyOC0xNy4zIDkuNC0xLjcgMTguNC0uNiAyNi41IDQuOSA1LjMgMy41IDkuMSA4LjMgMTEuNiAxNC4xLjYuOS4yIDEuNC0xIDEuN3oiLz48cGF0aCBmaWxsPSIjMDBBREQ4IiBkPSJNMTg2LjIgMTU0LjZjLTkuMS0uMi0xNy40LTIuOC0yNC40LTguOC01LjktNS4xLTkuNi0xMS42LTEwLjgtMTkuMy0xLjgtMTEuMyAxLjMtMjEuMyA4LjEtMzAuMiA3LjMtOS42IDE2LjEtMTQuNiAyOC0xNi43IDEwLjItMS44IDE5LjgtLjggMjguNSA1LjEgNy45IDUuNCAxMi44IDEyLjcgMTQuMSAyMi4zIDEuNyAxMy41LTIuMiAyNC41LTExLjUgMzMuOS02LjYgNi43LTE0LjcgMTAuOS0yNCAxMi44LTIuNy41LTUuNC42LTggLjl6bTIzLjgtNDAuNGMtLjEtMS4zLS4xLTIuMy0uMy0zLjMtMS44LTkuOS0xMC45LTE1LjUtMjAuNC0xMy4zLTkuMyAyLjEtMTUuMyA4LTE3LjUgMTcuNC0xLjggNy44IDIgMTUuNyA5LjIgMTguOSA1LjUgMi40IDExIDIuMSAxNi4zLS42IDcuOS00LjEgMTIuMi0xMC41IDEyLjctMTkuMXoiLz48L3N2Zz4=

[release]: https://github.com/micnncim/stardust/releases
[release-badge]: https://img.shields.io/github/v/release/micnncim/stardust?style=for-the-badge&logo=github

[license]: LICENSE
[license-badge]: https://img.shields.io/github/license/micnncim/stardust?style=for-the-badge

