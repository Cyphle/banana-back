FROM alpine:latest

# Installer les dépendances runtime
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Créer un utilisateur non-root
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Copier le binaire depuis le build
COPY target/x86_64-unknown-linux-musl/release/banana-back /usr/local/bin/banana-back

# S'assurer que le binaire est exécutable
RUN chmod +x /usr/local/bin/banana-back

# Changer vers l'utilisateur non-root
USER appuser

# Exposer le port (ajustez selon votre app)
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD /usr/local/bin/banana-back --health || exit 1

CMD ["banana-back"]