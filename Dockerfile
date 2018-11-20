FROM golang:1.11.2-alpine3.8

ADD . /go/src/github.com/netlify/git-gateway

RUN apk update && apk add wget && rm -rf /var/cache/apk/*

RUN apk add --update alpine-sdk

RUN go get github.com/Masterminds/glide

RUN adduser -D -u 1000 netlify && cd /go/src/github.com/netlify/git-gateway && make deps build && mv git-gateway /usr/local/bin/

USER netlify
CMD ["git-gateway"]
