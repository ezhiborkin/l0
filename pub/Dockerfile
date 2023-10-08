# Use an official Golang runtime as a parent image
FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

# Install any needed dependencies
RUN go mod download
RUN go build -o pub ./cmd/main.go 

# Command to run the executable
CMD ["./pub"]