# Kirk
Kirk is a Slack bot to integrate Slack and JIRA.  Kirk can:
- Automatically post links to JIRA when a JIRA issue ID is detected in a channel Kirk is a member of

## Deployment
Kirk requires the following environment variables to deploy:
- `SLACK_BOT_ACCESS_TOKEN`: Provided by Slack
- `SLACK_CLIENT_ID`: Provided by Slack
- `SLACK_CLIENT_SECRET`: Provided by Slack
- `PUBLIC_HOST`: The publicly accessible hostname and port Kirk is accessible at (e.g. `localhost:8080` for local development)

Environment variables can be configured by add a `.env` file at the app root (alongside `main.go`)

For local development, set `APP_ENV=development` to have kirk automatically restart when deployed with docker-compose