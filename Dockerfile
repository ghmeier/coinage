FROM golang

ENV PORT "8080"
ADD ./coinage /go/bin/coinage
ADD ./config-dev.json /go/bin/config.json

ENTRYPOINT /go/bin/coinage

EXPOSE 8081
