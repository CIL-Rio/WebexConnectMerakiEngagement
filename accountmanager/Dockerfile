FROM golang:1.21-alpine3.18 as builder 
WORKDIR /build/
COPY . .
RUN go build -o main

FROM alpine:3.18 as final
WORKDIR /app
COPY --from=builder /build/main .
CMD [ "./main" ]
