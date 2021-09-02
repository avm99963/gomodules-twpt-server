FROM node:16-alpine as build-stage

WORKDIR /app

COPY package*.json /app

RUN npm install

COPY ./ /app/

RUN ./node_modules/webpack-cli/bin/cli.js --mode production

FROM nginx:1.21-alpine
LABEL maintainer "me@avm99963.com"

COPY --from=build-stage /app/dist /usr/share/nginx/html
