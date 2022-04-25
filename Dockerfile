FROM golang:1.18 AS build

WORKDIR /app
COPY . ./
WORKDIR /app/server
RUN go build

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app /app
EXPOSE 8080
USER nonroot:nonroot
WORKDIR /app/server
ENTRYPOINT ["./server"]