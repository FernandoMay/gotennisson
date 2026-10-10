FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/gotennisson .

FROM scratch
COPY --from=build /out/gotennisson /gotennisson
EXPOSE 8080
ENTRYPOINT ["/gotennisson"]
