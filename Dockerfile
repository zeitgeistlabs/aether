FROM golang:1.12

WORKDIR /go/src/compute-aether
COPY . .

ENV GO111MODULE=on
#ENV GOFLAGS="-mod=vendor"

#RUN go get -d -v ./...
RUN go install -v ./...

CMD ["compute-aether"]
