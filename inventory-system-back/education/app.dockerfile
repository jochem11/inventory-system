FROM golang:1.24-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/jochem11/inventory-system-back
COPY go.mod go.sum ./
COPY vendor vendor
COPY education education
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./education/cmd/education

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
