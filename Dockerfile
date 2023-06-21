#FROM golang:1.17

#WORKDIR /app

#COPY go.mod go.sum ./

#RUN go mod download

#COPY . .

#RUN go build -o gateway ./gateway-test
#RUN go build -o vivasoft-main ./vivasoft
#RUN go build -o tigerit-main ./tigerit
#RUN go build -o cefalo-main ./cefalo
#RUN go build -o enosis-main ./enosis

#EXPOSE 8080
#EXPOSE 8000
#EXPOSE 8001
#EXPOSE 8002
#EXPOSE 8003

#CMD ["sh", "-c", "./gateway & ./vivasoft-main & ./tigerit-main & ./cefalo-main & ./enosis-main"]

# Default to Go 1.17.1
ARG GO_VERSION=1.17.1

# Start from golang v1.17.1 base image
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable support for modules.
WORKDIR /app

# Import the code from the context.
COPY ./ ./

# Build the Go apps
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o gateway ./gateway-test
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o vivasoft-main ./vivasoft
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o tigerit-main ./tigerit
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o cefalo-main ./cefalo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o enosis-main ./enosis

# Start a new stage from scratch
FROM alpine:3.14 AS final

# Create the user and group files from the builder stage.
COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import curl.
RUN apk add --no-cache curl

# Import the compiled executables from the builder stage.
#COPY --from=builder /app/gateway /app/vivasoft-main /app/tigerit-main /app/cefalo-main /app/enosis-main /app/
COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/vivasoft-main /app/vivasoft-main
COPY --from=builder /app/tigerit-main /app/tigerit-main
COPY --from=builder /app/cefalo-main /app/cefalo-main
COPY --from=builder /app/enosis-main /app/enosis-main


# Set the user to nobody
USER nobody

# Expose the necessary ports for the servers
EXPOSE 8080
EXPOSE 8000
EXPOSE 8001
EXPOSE 8002
EXPOSE 8003

# Use ENTRYPOINT to specify the entrypoint command
ENTRYPOINT ["sh", "-c", "/app/gateway & /app/vivasoft-main & /app/tigerit-main & /app/cefalo-main & /app/enosis-main"]


