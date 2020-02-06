# stardust

![[actions-workflow-test][actions-workflow-test-badge]][actions-workflow-test]
![[docker-build][docker-build-badge]][docker-build]
![[release][release-badge]][release]
![[license][license-badge]][license]

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
[docker-build]: https://hub.docker.com/r/micnncim/stardust
[release]: https://github.com/micnncim/stardust/releases
[license]: LICENSE

[actions-workflow-test-badge]: https://img.shields.io/github/workflow/status/micnncim/stardust/Test?label=Test&style=for-the-badge&logo=github
[docker-build-badge]: https://img.shields.io/docker/cloud/build/micnncim/stardust?logo=docker&style=for-the-badge
[release-badge]: https://img.shields.io/github/v/release/micnncim/stardust?style=for-the-badge
[license-badge]: https://img.shields.io/github/license/micnncim/stardust?style=for-the-badge
