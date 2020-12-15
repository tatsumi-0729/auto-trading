FROM golang:latest

# sqlite3のインストール
RUN apt-get -y update \
&& apt-get -y upgrade \
&& apt-get install -y sqlite3 libsqlite3-dev

# gitの諸々のインストール
RUN go get github.com/gorilla/websocket \
&& go get github.com/mattn/go-sqlite3 \
&& go get github.com/smartystreets/goconvey \
&& go get gopkg.in/ini.v1

# RUN go mod init app

# ENV GO111MODULE=on
# ENV GOPATH=

WORKDIR /go/src/auto-trading/app