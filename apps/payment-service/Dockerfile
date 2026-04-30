FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/payment-service .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/payment-service /payment-service
EXPOSE 8083
ENTRYPOINT ["/payment-service"]
