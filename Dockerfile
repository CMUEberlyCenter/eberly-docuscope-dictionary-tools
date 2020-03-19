FROM alpine:latest
WORKDIR /app
COPY ./docuscope-rules /app/
ENTRYPOINT ["/app/docuscope-rules"]