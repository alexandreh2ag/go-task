FROM golang:1.21.0-alpine3.17 as build
RUN apk update && apk add git make
RUN mkdir /build
WORKDIR /build
COPY . .
RUN export CGO_ENABLED=0 && make build

FROM alpine:3.17.2
COPY --from=build /build/bin/gtask /bin/gtask
EXPOSE 9469
ENTRYPOINT  [ "/bin/gtask" ]
