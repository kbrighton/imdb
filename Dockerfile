FROM golang:1.15-alpine AS builder
RUN apk add --no-cache make git
RUN apk add --no-cache ca-certificates

WORKDIR /go/src/imdb
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build /go/src/imdb
RUN wget https://datasets.imdbws.com/title.basics.tsv.gz
RUN gunzip title.basics.tsv.gz

ENV HOST 0.0.0.0
ENV PORT 8080

EXPOSE 8080

ENTRYPOINT ["./imdb"]
