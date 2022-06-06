## Sample Dockerfile to package Frontend-GUI as a docker image.
# Sample build command:
# docker build --force-rm --squash -t $shortname$_fe:$version$ .

#FROM node:13.6-alpine3.11 AS builder
#FROM node:12.22-alpine3.15 AS builder
FROM node:12-alpine AS builder
MAINTAINER $author$
RUN mkdir /build
COPY . /build
RUN cd /build && rm -rf .env dist node_modules package-lock.json && npm install && npm run build

FROM nginx:1.17-alpine
COPY nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /usr/share/nginx/html/app
COPY --from=builder /build/dist /usr/share/nginx/html/app
COPY --from=builder /build/dist/favicon*.* /usr/share/nginx/html/
COPY --from=builder /build/dist/manifest.json /usr/share/nginx/html
COPY index.html /usr/share/nginx/html
EXPOSE 80
RUN chown nginx.nginx /usr/share/nginx/html/ -R
