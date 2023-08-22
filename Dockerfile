FROM node:18-alpine as builder
WORKDIR /app
COPY . .
RUN npm i -g pkg
RUN npm install
RUN pkg prlint.js -o bin/prlint -C Brotli

FROM alpine:3.18
RUN apk add --no-cache git
COPY --from=builder /app/bin/prlint /usr/local/bin/prlint
RUN chmod +x /usr/local/bin/prlint

CMD [ "prlint" ]