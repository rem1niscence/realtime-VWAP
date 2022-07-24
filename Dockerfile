FROM golang:1.18-alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/github.com/rem1niscence/

COPY . /go/src/github.com/rem1niscence/realtime-VWAP/

WORKDIR /go/src/github.com/rem1niscence/realtime-VWAP
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/realtime_vwap ./cmd/*.go

FROM alpine:3.16.0
WORKDIR /app
COPY --from=builder /go/src/github.com/rem1niscence/realtime-VWAP/bin/realtime_vwap ./
CMD ["/app/realtime_vwap"]
