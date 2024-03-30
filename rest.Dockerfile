FROM golang:alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD ./go.mod .

COPY . .

RUN go build -trimpath -o rest ./restservice/rest.go

COPY ./.env ./build/.env

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /build/rest /build/rest

COPY --from=builder ./build/.env /build/.env

EXPOSE 8000

CMD ["/build/rest"]