FROM golang:latest as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV PORT=8080
#ENV MONGO_SERVER=host.docker.internal
WORKDIR /
COPY . .
RUN go build main.go && go build app/controllers/model_test.go

FROM alpine
COPY --from=builder /main /main
COPY --from=builder /app /app
COPY --from=builder /css /css
ENTRYPOINT ["/main"]