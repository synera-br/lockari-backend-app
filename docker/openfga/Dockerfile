# Dockerfile para OpenFGA Server
FROM openfga/openfga:v1.5.0

# Definir variáveis de ambiente padrão
ENV OPENFGA_HTTP_ADDR=0.0.0.0:8080
ENV OPENFGA_GRPC_ADDR=0.0.0.0:8081
ENV OPENFGA_LOG_LEVEL=info
ENV OPENFGA_AUTHN_METHOD=none
ENV OPENFGA_PLAYGROUND_ENABLED=false
ENV OPENFGA_DATASTORE_ENGINE=postgres

# Configurar usuário não-root para segurança
USER nobody

# Expor portas
EXPOSE 8080 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

# Comando padrão
CMD ["run"]
