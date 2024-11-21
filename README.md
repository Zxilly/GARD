# GitHub Actions Runner Daemon

Automatically add and remove GitHub Actions self-hosted runners in the Docker container.

## Environment Variables

The following environment variables are required for the application to function correctly:

| Variable Name       | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| `GITHUB_TOKEN`      | The GitHub token used for authentication.                                   |
| `GITHUB_REPOSITORY` | The GitHub repository in the format `owner/repo` where the runner is added. |
| `RUNNER_LOCATION`   | (Optional) The directory where the runner is located.                       |

