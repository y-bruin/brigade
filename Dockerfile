FROM golang:1.21-alpine AS builder

WORKDIR /workspace

RUN apk add --update --no-cache git && rm -rf /var/cache/apk/*
COPY go.mod go.sum /workspace/
RUN go mod download
COPY cmd /workspace/cmd
COPY internal /workspace/internal
RUN go build -o server ./cmd/cloud/main.go

FROM alpine
RUN apk add --update --no-cache ca-certificates tzdata && rm -rf /var/cache/apk/*
COPY --from=builder /workspace/server /usr/local/bin/server
CMD [ "/usr/local/bin/server" ]
