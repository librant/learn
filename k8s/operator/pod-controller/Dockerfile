FROM golang:1.18 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLE=0 go build -o pod-controller pod-controller.go

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/pod-controller .

CMD ["./pod-controller"]