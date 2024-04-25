FROM golang:1.21 AS build

WORKDIR /app
COPY . .
RUN make build

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/build/bin/api-go-financial-accounting /app

USER 1000

ENTRYPOINT [ "/app/api-go-financial-accounting" ][]