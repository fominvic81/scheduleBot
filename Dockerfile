FROM oven/bun:latest

WORKDIR /app

COPY ./package.json /app/

RUN bun install

COPY . /app/

CMD [ "bun", "run", "--watch", "src/main.ts" ]