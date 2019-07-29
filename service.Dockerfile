FROM golang:alpine AS build-env
RUN apk update && apk upgrade && apk add --no-cache bash git openssh
COPY ./backend/api/ /src/backend/api/
COPY ./backend/service/ /src/backend/app/
COPY ./go.* /src/
ENV GO111MODULE=on
RUN cd /src && go mod download
RUN cd /src/backend/app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM centurylink/ca-certs
COPY --from=build-env /src/backend/app /
ENTRYPOINT ["/app"]
EXPOSE 50051
