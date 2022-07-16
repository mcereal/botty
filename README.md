# botty [![build](https://github.com/mcereal/botty/actions/workflows/build.yml/badge.svg)](https://github.com/mcereal/botty/actions/workflows/build.yml) ![Go Report Card](https://goreportcard.com/badge/github.com/mcereal/botty) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**botty** is a simple Go application that posts open pull requests to a team communication app, and checks for stale PRs.

Currently supported apps are:

[Discord](https://discord.com)

[Slack](https://slack.com)

## Table of Contents

- [Requirements](https://github.com/mcereal/go-api-server#requirements)

  - [Install Go](https://github.com/mcereal/go-api-server#install-go)

  - [Setup Private Repositories](https://github.com/mcereal/go-api-server#setup-private-repositories)

  - [Install nodemon(optional)](https://github.com/mcereal/go-api-server#install-nodemon-globally)

- [Environment](https://github.com/mcereal/go-api-server#environment)

  - [Adding the .env](https://github.com/mcereal/go-api-server#adding-the-.env)

  - [Editing the config](https://github.com/mcereal/go-api-server#editing-the-config)

- [How to run](https://github.com/mcereal/go-api-server#how-to-run)

- [Test it out](https://github.com/mcereal/go-api-server#test-it-out)

- [Supplimental](https://github.com/mcereal/go-api-server#supplimental)

  - [Unit tests](https://github.com/mcereal/go-api-server#unit-tests)

  - [linting](https://github.com/mcereal/go-api-server#linting)

## Requirements

**botty** is currently running on Go version 1.18. You must have this installed to run. Optionally you can [install nodemon](https://www.npmjs.com/package/nodemon) **globally** if you want to hot reload the project on file save. You will have to have [Node.js](https://nodejs.org/en/download/) installed if you want to do this.

_note: Using nodemon is optional, it just makes it easy while developing._

### Install Go

- [Download Go Installer](https://go.dev/doc/install).

  _note: You can use [homebrew](https://brew.sh) to install with `brew install go`._

- After install you will need to add few environment variables to your path. Open your `.zshrc` or `.bashrc` file which should be in your `$HOME` directory.
- You can navigate to your `$HOME` directory with `cd $HOME`
- Add `GOPATH` and `GOROOT` to your `.zshrc` or `.bashrc`:
  ```
  export GOPATH=$HOME/go
  export GOROOT=/usr/local/go
  export PATH=$PATH:$GOPATH/bin
  export PATH=$PATH:$GOROOT/bin
  ```
- Save your `.zshrc` or `.bashrc` file.
- In your `$HOME` directory reload your `.zshrc` or `.bashrc` file.
  - For `.zshrc` run:
    ```
    $ source .zshrc
    ```
  - For `.bashrc` run:
    ```
    $ source .bashrc
    ```
- Make sure go is installed by running:

```
$ go version
```

### Setup Private Repositories

You need to tell Go to use ssh when downloading modules and also that the git repository is private.

#### Use ssh

- cd to your `$HOME` directory and open your `.gitconfig` file. Paste the following:

```
[url "git@github.com:"]
	insteadOf = https://github.com/
```

- Save
- In your `$HOME` directory create a `.netrc` file if it does not already exist. Paste the following:

```
machine github.com login YOUR_GITHUB_USERNAME_HERE password YOUR_PERSONAL_ACCESS_TOKEN_HERE
```

_note: Your GitHub username is not the same as your email. Open on your profile in GitHub to check yours. Your personal access token can be created following [these instructions](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)._

#### Private repo

You need to tell go that it needs to look for private repos.

- cd to your `$HOME` directory and open your `.zshrc` or `.bashrc` file. Paste the following:

```
$ export GOPRIVATE=*
```

- Save then run:

  - For `.zshrc` run:
    ```
    $ source .zshrc
    ```
  - For `.bashrc` run:
    ```
    $ source .bashrc
    ```

### Install nodemon globally:

```
$ npm i -g nodemon
```

The `-g` flag here is important. This is not a Node project, we're just using this package running in the background to ease development.

_note: You must have [Node.js](https://nodejs.org/en/download/) installed already for this._

## Environment

Environment variables and application configuration can be be managed with a `.env` file and the `config.yml`

### Adding the .env

The `.env` file is included in `.gitignore` to prevent sharing sensitive information. You must create a new `.env` file and add it to the project root.

You must provide the slack webhook URL:

```
SLACK_WEBHOOK_URL="https://exampleslackwebhookurl.com"
```

You can specify a specific port and environment. If not specified these will default to `PORT=8080` and `ENVIRONMENT="development"`.

```
PORT=8080
ENVIRONMENT="development"
```

### Editing the config

Edit the [config.yml](config.yml) with the desired settings. **botty** can send slack messages to any channel, monitor any repo, and tag any team. To set up a new team provide the team name, the Slack group Id, a unique name for the webhook environment variable, and any repos you would like to monitor.

```
  - name: MyCoolTeamName
    channel: MY_COOL_CHANNEL_SLACK_WEBHOOK_URL
    enable_cron: true
    cron_elapsed_duration: 14400000000000 #4 Hours
    org:  "your GitHub org or Username here"
    ignore_users:
    - BotThatMakesPRs
    repos:
    - "cool-repo1"
    - "cool-repo2"
    - "hello-world"
    - "cool-repo2-rework"
```

## How to run

```
$ go run main.go
```

If you have nodemon installed gloablly you can just run:

```
$ nodemon
```

nodemon just lets the code run automatically when you save a file so you dont have to restart the program every time.

## Test it out

Try making a **GET** request to the _/healthz_ endpoint.

```
$ curl http://localhost:8080/healthz
```

You should get a `200` status code and an "OK" response.

You can also make a **POST** request to the _/api/v1/githubwebhook/payload_ endpoint.

Provide a JSON body in the request with an [example payload](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request) from the github pull request event.

## Supplimental

Some additional useful information for developing in Go.

### Unit Tests

Run all unit tests. These are always suffixed with _test.go_

```
$ go test ./...
```

### linting

[golangcli-lint](https://golangci-lint.run) is linter aggregator for Go. Install golanglint-cli with:

```
$ brew install golangci-lint
```

#### linters used

You will also need to install [goimports](golang.org/x/tools/cmd/goimports), [revive](https://revive.run/docs), and [staticcheck](https://staticcheck.io/docs/). These packages are enabled in the [.golangci.yml](.golangci.yml).

```
$ go install golang.org/x/tools/cmd/goimports@latest
$ go install github.com/mgechev/revive@latest
$ go install honnef.co/go/tools/cmd/staticcheck@latest
```

#### Run golangcli-lint

```
$ golangci-lint run
```
