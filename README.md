# Codewars Bot

Have your hard work on [codewars](https://www.codewars.com) automagically populate your history on GitHub

Codewars bot is a program written in Go, meant to run at a daily frequency and commit a record of your work solving katas on Codewars to GitHub.

It comes with a `Dockerfile` to be built into an image for cloud deployment, and can be hosted anywhere that runs Docker containers.

Alternatively, you could just run it on your machine with Docker, or simply avoid using Docker and run it manually with `go run main.go`.

## Deployment

It can be run remotely on any cloud platform that knows what to do with a GitHub repo with a Dockerfile.

The easiest deployment strategy I would recommend is [Railway.](https://railway.com) Just fork this repo, create a project, create a service, browse the GitHub repo option for your fork. Then, add the ENV vars, and choose a cron schedule to run it.

## ENV

You will need the following:

- `GITHUB_PUBLIC_NAME` | Your human readable name on GitHub eg "John Doe."
- `GITHUB_EMAIL`       | Your email that GitHub associates with your commits, found in [email settings.](https://github.com/settings/emails) Often takes the format of `01234567+<USERNAME>@users.noreply.github.com.`
- `GITHUB_USERNAME`    | Your username on GitHub, can be found in the url when you visit your profile or repos.
- `GITHUB_REPO_NAME`   | Name of a separate GitHub repo you want the program to commit to.
- `GITHUB_TOKEN`       | Access token allowing commits to your repo of codewars history. Create one in [developer token settings.](https://github.com/settings/tokens)
- `CODEWARS_USERNAME`  | Your username on Codewars. Can be found in the url when visiting your profile.


