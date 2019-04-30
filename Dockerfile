FROM golang:1.12 AS builder

WORKDIR /opt/transponder

# For caching purposes, download modules first
ADD go.mod go.sum /opt/transponder/
RUN go mod download

# Now build the source statically
ADD . /opt/transponder/
RUN CGO_ENABLED=0   \
    GOOS=linux      \
    go build -o transponder

# ---
# Now the second lightweigth stage
FROM scratch

# CA certificates are often needed
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Here is our binary
COPY --from=builder /opt/transponder/transponder /usr/local/bin/transponder

USER transponder
EXPOSE 1492

VOLUME /var/local/tfstate

CMD [ "/usr/local/bin/transponder", "-data", "/var/local/tfstate" ]