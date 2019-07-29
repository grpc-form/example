FROM node:12 AS build-env
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package.json /usr/src/app
RUN npm install
COPY . /usr/src/app
RUN $(npm bin)/ng build --prod

FROM nginx:alpine
COPY --from=build-env /usr/src/app/dist/frontend /usr/share/nginx/html
COPY nginx-custom.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
