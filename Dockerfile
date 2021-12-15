FROM golang:1.16.6-alpine AS builder
WORKDIR /

RUN apk add git gcc g++ --no-cache
COPY ./ ./

RUN go mod download

RUN go build -o app ./main.go

FROM golang:1.16.6-alpine

RUN adduser -S gophershoney

COPY --from=builder /app /home/gophershoney/

RUN mkdir /home/gophershoney/images



USER gophershoney

CMD ["/home/gophershoney/app"]

