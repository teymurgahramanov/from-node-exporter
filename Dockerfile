FROM golang:1.21.6 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.18.4 as final
WORKDIR /app
COPY --from=build /app/from-node-exporter .
ENTRYPOINT [ "./from-node-exporter" ]