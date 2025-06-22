FROM node:24-alpine

WORKDIR /usr/src/app

EXPOSE 3001

RUN npm install express

COPY . .

CMD [ "node", "index.js" ]