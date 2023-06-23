FROM golang:1.20-alpine as BUILD

COPY . /build
WORKDIR /build

RUN go build sample-app/main.go

FROM scratch

COPY --from=BUILD /build/main /main
COPY --from=BUILD /build/sample-app/config /config
COPY --from=BUILD /build/sample-app/infrastructure/_crdb/node_1/certs /dbcerts

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "/main" ]