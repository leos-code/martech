FROM node:alpine3.12 as builder1
WORKDIR /app
ADD . .
RUN cd web && yarn && yarn build

FROM golang:alpine as builder2
RUN apk add --no-cache git make build-base ffmpeg
WORKDIR /app
ADD . .
RUN go install -v cmd/web/web.go && \
    go install -v cmd/rta/server/rta_server.go

ENTRYPOINT ["./web"]

FROM alpine:latest
ENV TZ "PRC"
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder1 /app/web/dist dist
COPY --from=builder2 /go/bin/.  .
COPY --from=builder2 /app/configs/.  .
COPY --from=builder2 /app/docs docs
ENTRYPOINT ["./web"]