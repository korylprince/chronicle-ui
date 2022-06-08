FROM golang:1-alpine as builder

ARG VERSION

RUN go install github.com/korylprince/fileenv@v1.1.0
RUN go install "github.com/korylprince/chronicle-ui@$VERSION"

FROM alpine:3.16

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/chronicle-ui /chronicle-ui

CMD ["/fileenv", "sh", "/setenv.sh", "/chronicle-ui"]
