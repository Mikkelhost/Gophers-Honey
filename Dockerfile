FROM golang:1.16.6-alpine AS builder
WORKDIR /

RUN apk add git gcc g++ --no-cache
COPY ./ ./

RUN go mod download

RUN go build -o app ./main.go

FROM golang:1.16.6-alpine

COPY --from=builder /app /go/

EXPOSE 8000

CMD ["/go/app"]

