FROM golang:1.12

WORKDIR /go/src/aether
COPY . .

ENV GO111MODULE=on
ENV GOFLAGS="-mod=vendor"

RUN go install -v ./...



FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/aether .

CMD ["./aether"]
