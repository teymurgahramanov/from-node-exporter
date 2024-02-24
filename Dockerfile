FROM golang:1.21.6 AS build
WORKDIR /app
COPY . .
RUN \
    if ! test -f bin/from-node-exporter; then \
      CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bin/from-node-exporter; \
    fi
FROM scratch as final
COPY --from=build /app/bin/from-node-exporter /
ENTRYPOINT [ "/from-node-exporter" ]