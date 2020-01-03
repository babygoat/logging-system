ARG server_user=poc-be

FROM golang:1.12.9-alpine3.10 As build

RUN apk add --update --no-cache tzdata ca-certificates git && \
    update-ca-certificates

WORKDIR /go/cache

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

WORKDIR /src

COPY go.mod go.sum server.go logger.go ./

RUN go install

ARG server_user

RUN adduser -D -g '' ${server_user}

FROM alpine:3.10

COPY --from=build /go/bin /usr/local/bin /usr/local/bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd

ARG server_user

USER ${server_user}

EXPOSE 8080 

CMD poc-backend
