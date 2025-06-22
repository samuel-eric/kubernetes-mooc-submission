FROM node:24-alpine

WORKDIR /usr/src/app

EXPOSE 3000

RUN npm install express uuid

COPY . .

CMD [ "node", "index.js" ]