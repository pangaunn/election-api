FROM golang:1.17.2-buster as builder

WORKDIR /app

COPY go.* /
RUN go mod download
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o webapp ./main.go
RUN ls

FROM alpine:3.12.0
RUN apk add bash
WORKDIR /app
COPY --from=builder app/webapp .
COPY --from=builder app/wait-for-it.sh .
RUN chmod +x wait-for-it.sh
EXPOSE 3000
ENTRYPOINT [ "/bin/bash", "-c" ]
CMD ["./wait-for-it.sh election-redis:6379 && ./webapp"]