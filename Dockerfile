FROM golang:1.4.2

RUN go get github.com/constabulary/gb/...

ADD . /usr/src/sink

RUN cd /usr/src/sink && gb build

VOLUME /usr/src/sink/bin
