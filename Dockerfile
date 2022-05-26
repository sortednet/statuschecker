FROM alpine

COPY bin/statuschecker /usr/local/bin/statuschecker

ENTRYPOINT ["/usr/local/bin/statuschecker"]



