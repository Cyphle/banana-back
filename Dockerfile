FROM alpine:latest

RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app

WORKDIR /app

COPY config ./config
COPY target/release/banana-back ./banana-back

RUN chown -R app:app /app

USER app

ENTRYPOINT ["./banana-back"]