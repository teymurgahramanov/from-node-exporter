FROM golang:1.21.6 AS build
WORKDIR /app
COPY . .
RUN go get golang.org/x/lint/golint
RUN ${GOPATH}/bin/golint -set_exit_status ./...
RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch as final
COPY --from=build /app/from-node-exporter /
ENTRYPOINT [ "/from-node-exporter" ]