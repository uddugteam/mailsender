# build container
FROM golang:1.15 AS builder

WORKDIR /

RUN apt-get update \
    && apt-get -y install make openssh-client ca-certificates && update-ca-certificates \

ADD . / service/
WORKDIR /service

RUN make build

# live container
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /service/mailsender /

ENTRYPOINT ["/mailsender"]