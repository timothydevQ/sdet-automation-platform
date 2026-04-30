FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 go build -o /out/order-service .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/order-service /order-service
EXPOSE 8084
ENTRYPOINT ["/order-service"]
