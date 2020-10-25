FROM golang:latest as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /
COPY . .
RUN go build main.go

FROM alpine
COPY --from=builder /main /main
COPY --from=builder /app /app
COPY --from=builder /css /css
ENTRYPOINT ["/main"]