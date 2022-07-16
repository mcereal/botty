##
## Build
##

FROM golang:1.18.2-buster AS build

## How to enable for a private repo:
## This is the Github username and associated API Key. 
## Pass these in from a pipeline Environment Variable or just have them exported on your local machine if building a local image.
## If running locally you need to create a .netrc file like the one shown below.
# ARG GITHUB_USERNAME
# ARG GITHUB_KEY

## Set the GOPRIVATE environment variable to the name of the private repo so Go knows about it
# ENV GOPRIVATE=github.com/mcereal/go-api-server

## Creates a .netrc file so go has access to Username and API Key
# RUN echo "machine private.git.local login $GITHUB_USERNAME password $GITHUB_KEY" > ~/.netrc

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /botty

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /botty /botty

COPY --from=build /app/config.yml config.yml

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/botty"]