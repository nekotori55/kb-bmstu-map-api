# BUILD STAGE
FROM golang:alpine AS builder

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

######################################################
# PROD BUILD
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=dockertest Version=0.0.1
EXPOSE 3000


#####################################################
# DEVELOPMENT BUILD (AUTO-REBUILD ON CHANGE)
FROM golang:alpine AS dev
WORKDIR /go/src/app
RUN mkdir "/go/build"
COPY . .
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -build="go build -o /go/build/app" -command="/go/build/app"