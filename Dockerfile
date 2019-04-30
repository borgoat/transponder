FROM golang:1.12

ADD . /opt/transponder
WORKDIR /opt/transponder

RUN go build -o /usr/local/bin/transponder

CMD [ "/usr/local/bin/transponder" ]