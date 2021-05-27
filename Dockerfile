FROM golang:1.15 AS builder
RUN mkdir /app 
ADD . /app 
WORKDIR /app 

RUN CGO_ENABLED=0 GOS=linux go build -o app cmd/server/main.go

FROM alpine:latest as production 
COPY --from=builder /app .
CMD ["./app"]