FROM scratch

ARG BINARY_NAME

ENTRYPOINT [ "/bin" ]

COPY $BINARY_NAME /bin