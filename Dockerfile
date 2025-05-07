FROM golang:1.24.2-bookworm AS build

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin ./bin

ENTRYPOINT ["app"]
