# build stage

FROM golang:bullseye AS BuildStage

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /gobuild

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /teslamate-addr-fix .

# deploy stage

FROM debian:bullseye-slim

ENV TZ=Asia/Shanghai
ENV TESLAMATE_ADDR_FIX_ENV=docker

WORKDIR /
COPY --from=BuildStage /teslamate-addr-fix /teslamate-addr-fix

ENTRYPOINT ["./teslamate-addr-fix"]
