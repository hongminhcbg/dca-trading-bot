
##
## Build
##
FROM golang:1.16-alpine AS build

RUN apk add build-base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dca_bot main.go

##
## Deploy
##
FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/dca_bot ./dca_bot
COPY --from=build /app/config.yaml ./config.yaml
EXPOSE 80
ENTRYPOINT ["/app/dca_bot"]