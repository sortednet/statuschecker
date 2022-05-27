FROM alpine

COPY bin/statuschecker /usr/local/bin/statuschecker

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/statuschecker"]



