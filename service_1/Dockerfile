FROM golang:1.23 AS build
WORKDIR /app
COPY main.go .
ENV GO111MODULE=off
RUN go build -o main main.go

FROM gcr.io/distroless/base-debian12
COPY --from=build /app/main /usr/local/bin/main
EXPOSE 8001
CMD ["/usr/local/bin/main"]