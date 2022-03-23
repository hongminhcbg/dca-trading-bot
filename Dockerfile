FROM golang:1.16 as builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o server *.go

FROM scratch
WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /app/config.yaml /app/config.yaml
COPY --from=builder --chown=nonroot:nonroot /app/server /app/server
ENTRYPOINT ["/app/server"]
