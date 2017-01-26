FROM golang

ENV PORT "8080"
ADD ./coinage /go/bin/coinage
ADD ./config.json /go/bin/config.json

ENTRYPOINT /go/bin/coinage

EXPOSE 8081
