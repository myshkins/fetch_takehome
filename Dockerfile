# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23.4 AS build-stage

WORKDIR /app

# COPY go.mod go.sum ./

COPY . .

RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/health_check

# Final stage
FROM gcr.io/distroless/base-debian11

# Set the working directory
WORKDIR /

# Copy the binary from the build stage
COPY --from=build-stage /app/cmd/health_check/health_check .
COPY --from=build-stage /app/sample_input.yaml .

ENTRYPOINT [ "./health_check" ]

CMD ["-config-file=./sample_input.yaml"]

