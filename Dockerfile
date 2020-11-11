FROM golang:1.15.1
RUN mkdir /go/mail-srvc
WORKDIR /go/mail-srvc/cmd/server
ADD . /go/mail-srvc
##COPY . /test/

ENV GO111MODULE=on

RUN go mod download
RUN go get -v
## -o создаст exe файл в текущей дериктории
##RUN go build -o main

EXPOSE 8081 5432 6379