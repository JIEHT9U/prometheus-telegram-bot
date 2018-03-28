FROM golang:1.10.0-alpine3.7 AS build

WORKDIR /build
COPY ./src/bot src/prometheus_bot/
COPY ./vendor vendor

RUN apk add --no-cache git ;\
    go get github.com/constabulary/gb/... ; \
    gb build 


FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=build /build/bin/prometheus_bot /opt/prometheus_bot 
COPY ./template template 
COPY ./default  default 
EXPOSE 9087
ENTRYPOINT ["/opt/prometheus_bot"]



#RUN apk add --no-cache ca-certificates gcc musl-dev
#RUN /bin/sh -c "apk add --update gcc musl-dev; go build -ldflags \"-s -w\" -a -o bot"
