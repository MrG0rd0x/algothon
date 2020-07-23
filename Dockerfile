FROM golang:1.13-alpine3.12 AS build
RUN mkdir /app
ADD . /app
WORKDIR /app/web
RUN apk add --no-cache git
RUN go mod download
RUN go build -o webserver ./cmd

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
COPY web/files/templates/* ./templates/
COPY web/files/static/* ./static/
COPY --from=build /app/web/webserver .
ENTRYPOINT [ "./webserver" ]