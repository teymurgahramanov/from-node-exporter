FROM golang:1.21.6 AS build
WORKDIR /app
COPY . .
RUN golint .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch as final
COPY --from=build /app/from-node-exporter /
ENTRYPOINT [ "/from-node-exporter" ]