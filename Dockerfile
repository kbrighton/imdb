FROM golang:1.15-alpine AS builder
RUN apk add --no-cache make git

WORKDIR /go/src/imdb
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build /go/src/imdb

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/imdb .

ENV HOST 0.0.0.0
ENV PORT 8080

EXPOSE 8080

ENTRYPOINT ["./imdb"]
