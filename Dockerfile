FROM golang:1.18.3-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY svc/catalog svc/catalog/.

RUN CGO_ENABLED=0 go build -o /target/app svc/catalog/cmd/main.go


FROM scratch

WORKDIR /
COPY --from=builder /target/app /app

EXPOSE 8080
USER 1101:1101

ENTRYPOINT ["/app"]