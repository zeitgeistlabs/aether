FROM golang:1.12

WORKDIR /go/src/app
COPY . .

ENV GO111MODULE=on
ENV GOFLAGS="-mod=vendor"

RUN go install -v ./...

CMD ["container"]
