FROM docker.io/golang:1.20-alpine as builder

WORKDIR /workspace

COPY main.go go.mod go.sum ./

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-w' -o latency-monitor main.go

FROM alpine

COPY --from=builder /workspace/latency-monitor latency-monitor

ENTRYPOINT ["./latency-monitor"]