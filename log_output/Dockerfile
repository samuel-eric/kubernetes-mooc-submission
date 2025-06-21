FROM node:24-alpine

WORKDIR /usr/src/app

RUN npm install uuid

COPY . .

CMD [ "node", "index.js" ]