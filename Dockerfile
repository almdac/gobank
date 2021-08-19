FROM golang:alpine AS dev

RUN apk update && apk update && apk add --no-cache build-base
WORKDIR /server
COPY . .

FROM dev AS test
WORKDIR /server
CMD [ "go", "test", "-v" ]

FROM dev AS build
WORKDIR /server
RUN go build -o bin/bank

FROM alpine AS prod
COPY --from=build /server/bin/bank .
CMD ["./bank"]