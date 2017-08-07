# Kirk

[![Build Status](https://travis-ci.org/ContainX/kirk.svg)](https://travis-ci.org/ContainX/kirk)

Kirk is a Slack bot to integrate Slack and JIRA.  Kirk can:
- Automatically post links to JIRA when a JIRA issue ID is detected in a channel Kirk is a member of

## Deployment
Kirk requires the following environment variables to deploy:
- `SLACK_CLIENT_ID`: Provided by Slack
- `SLACK_CLIENT_SECRET`: Provided by Slack
- `PUBLIC_HOST`: The publicly accessible hostname and port Kirk is accessible at (e.g. `localhost:8080` for local development)
- `HOST`: The host where the datadog agent is running

Environment variables can be configured by add a `.env` file at the app root (alongside `main.go`)

For local development, set `APP_ENV=development` to have kirk automatically restart when deployed with docker-compose

## DataDog

Notes for running datadog agent:
```
docker run -d --name dd-agent \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v /proc/:/host/proc/:ro \
    -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
    -e API_KEY={apikey} \
    -e SD_BACKEND=docker \
    -p 8125:8125/udp \
    -p 8126:8126/tcp \
    -e DD_APM_ENABLED=true \
    -e DD_HOSTNAME=local \
    -e NON_LOCAL_TRAFFIC=true \
    datadog/docker-dd-agent
```
