## Sample Dockerfile to package Backend-API as a docker image.
# Sample build command:
# docker build --force-rm --squash -t $shortname$_be:$version$ .

FROM golang:1.17-alpine AS builder
LABEL maintainer="$author$"
RUN apk add git build-base \
    && mkdir /build
COPY . /build
RUN cd /build && go build -o main

FROM alpine:3.12
RUN mkdir /app
COPY --from=builder /build/main /app/main
COPY --from=builder /build/config /app/config
COPY --from=builder /build/frontend /app/frontend
RUN apk add --no-cache -U tzdata bash ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/$timezone$ /etc/localtime \
    && chmod 711 /app/main \
    && rm -rf /var/cache/apk/*
WORKDIR /app
EXPOSE 8000
CMD ["/app/main"]
#ENTRYPOINT /app/main
