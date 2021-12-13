FROM ubuntu

ARG BINARY_NAME

RUN apt-get update; \
    apt-get install -y ca-certificates

ENTRYPOINT [ "/usr/local/bin/entrypoint" ]

COPY $BINARY_NAME /usr/local/bin/entrypoint
