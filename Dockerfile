FROM golang:1.24.5-alpine AS builder

WORKDIR /go/src/app

RUN apk add --no-cache upx ca-certificates tzdata git

ARG VERSION=main
ARG BUILD="N/A"

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux

COPY go.mod go.sum /go/src/app/

RUN go mod download \
  && go mod tidy

COPY . /go/src/app/

RUN go build -a -installsuffix cgo -ldflags="-w -s -X github.com/AkashRajpurohit/git-sync/pkg/version.Version=${VERSION} -X github.com/AkashRajpurohit/git-sync/pkg/version.Build=${BUILD}" -o git-sync . \
  && upx -q git-sync

# Application image
FROM alpine:latest
WORKDIR /opt/go

LABEL maintainer="AkashRajpurohit <me@akashrajpurohit.com>"

# Install git since it's required for the application
RUN apk add --no-cache git su-exec

RUN mkdir -p /git-sync /backups

COPY --from=builder /go/src/app/git-sync /opt/go/git-sync
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh", "/opt/go/git-sync"]
CMD ["--config", "/git-sync/config.yaml", "--backup-dir", "/backups"]
