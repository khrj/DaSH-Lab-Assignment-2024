FROM alpine:3.19

WORKDIR /usr/src/app
RUN apk add --no-cache npm nodejs

# npm dependencies
COPY server/package*.json ./
RUN npm config set registry https://registry.npmjs.com/
RUN npm ci

# Bundle app source
COPY server .
RUN npm run build

CMD [ "node", "build/server.js" ]
