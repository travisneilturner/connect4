# build stage
FROM golang:alpine AS build-env
WORKDIR /app
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/connect4

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/connect4 .
ENTRYPOINT ./connect4