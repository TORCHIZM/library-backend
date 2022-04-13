FROM golang:1.16-alpine AS builder

ENV GO111MODULE=on

# Download dependencies
WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code from the host and build it
COPY . /usr/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/api/main.go

# Copy binary into second stage and run app
FROM scratch
COPY --from=builder /app ./
EXPOSE 3000
CMD [ "./app" ]