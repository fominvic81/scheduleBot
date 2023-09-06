FROM node:16.18.0

WORKDIR /app

COPY ./package.json /app/

RUN yarn

RUN yarn prisma generate

COPY . /app/

CMD [ "yarn", "serve" ]