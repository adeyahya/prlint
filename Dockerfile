# Multi build stages Dockerfile

# This is for compiling, don't publish this
FROM golang:1.19-alpine AS builder
RUN echo "RUNNING AS BUILD"
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o /app/prlint .

# Image with exec only, do publish this image
FROM alpine:3.19 AS executable
RUN echo "RUNNING AS EXEC"
COPY --from=builder /app/prlint /usr/local/bin/prlint

# This is only for internal testing
FROM executable AS test
RUN echo "RUNNING AS TEST"
WORKDIR /go/src/app
COPY . .
ENTRYPOINT ["/usr/local/bin/prlint"]
