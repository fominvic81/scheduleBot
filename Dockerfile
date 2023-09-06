FROM node:16.18.0

WORKDIR /app

COPY ./package.json /app/

RUN yarn

COPY ./prisma/ /app/

RUN yarn prisma generate

COPY . /app/

CMD [ "yarn", "serve" ]