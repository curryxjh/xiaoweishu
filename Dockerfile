FROM ubuntu:22.04

COPY webook /app/webook

WORKDIR /app

ENTRYPOINT ["/app/webook"]