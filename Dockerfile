# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Dynasti Madani Faz <fazdynastimadani@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN apt-get update -y && apt-get install dh-autoreconf nasm cmake libpng-tools libpng-dev -y

RUN deps/build-deps-linux.sh

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build  -ldflags '-w -s -extldflags "-static"' -a -installsuffix cgo -o main


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy ENV 
COPY --from=builder /app/.env .

RUN mkdir -p uploads/

# Command to run the executable
CMD ["./main"] 