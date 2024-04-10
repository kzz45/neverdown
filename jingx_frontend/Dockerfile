FROM node:18.15.0 as builder

WORKDIR /temp
COPY . .

RUN npm i

RUN npm run build:prod

FROM nginx

LABEL maintainer="kongzz"

COPY --from=builder /temp/dist /usr/share/nginx/html
COPY --from=builder /temp/dist/default.conf /etc/nginx/conf.d/
