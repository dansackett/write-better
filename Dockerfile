FROM golang:1.5

ADD app/ /go/src/write-better
ADD templates/ /go/templates

RUN go get github.com/dansackett/go-text-processors
RUN go get github.com/rainycape/unidecode
RUN go install write-better
RUN echo "export GOPATH=/etc/gopath" >> /etc/profile

ENTRYPOINT /go/bin/write-better
EXPOSE 8080
