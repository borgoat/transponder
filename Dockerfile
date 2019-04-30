FROM golang:1.12 AS builder

WORKDIR /opt/transponder

# For caching purposes, download modules first
ADD go.mod go.sum /opt/transponder
RUN go mod download

# Now build the source statically
ADD . /opt/transponder
RUN CGO_ENABLED=0   \
    GOOS=linux      \
    go build -o transponder

FROM scratch

COPY --from=builder /opt/transponder/transponder /usr/local/bin/transponder

EXPOSE 1492

CMD [ "/usr/local/bin/transponder" ]