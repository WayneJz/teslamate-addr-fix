FROM golang:alpine

ENV TZ=Asia/Shanghai
ENV TESLAMATE_ADDR_FIX_ENV=docker

WORKDIR $GOPATH/src/github.com/WayneJz/teslamate-addr-fix
COPY teslamate-addr-fix $GOPATH/src/github.com/WayneJz/teslamate-addr-fix

ENTRYPOINT ["./teslamate-addr-fix"]
