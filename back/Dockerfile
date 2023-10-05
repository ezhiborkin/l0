FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get install -y postgresql-client

# make wait-for-it.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o back ./cmd/main.go

CMD ["./back"]