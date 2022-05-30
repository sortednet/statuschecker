FROM alpine

COPY bin/statuschecker /usr/local/bin/statuschecker
COPY config/config.yaml /usr/local/status-checker-config.yaml

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/statuschecker", "--config", "/usr/local/status-checker-config.yaml"]



