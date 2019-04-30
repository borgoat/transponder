FROM golang:1.12 AS builder

ADD . /opt/transponder
WORKDIR /opt/transponder

RUN go build -o transponder


FROM scratch

COPY --from=builder /opt/transponder/transponder /usr/local/bin/transponder

CMD [ "/usr/local/bin/transponder" ]