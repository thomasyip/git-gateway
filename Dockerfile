FROM golang:1.11.2-stretch

ADD . /go/src/github.com/netlify/git-gateway

RUN go get github.com/Masterminds/glide

RUN useradd -m netlify && cd /go/src/github.com/netlify/git-gateway && make deps build && mv git-gateway /usr/local/bin/

USER netlify
CMD ["git-gateway"]
