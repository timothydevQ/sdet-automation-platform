FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 go build -o /out/notification-service .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/notification-service /notification-service
ENTRYPOINT ["/notification-service"]
