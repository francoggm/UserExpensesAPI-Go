FROM golang:1.21.0 as builder
COPY go.mod go.sum /go/src/github.com/francoggm/go_expenses_api/
WORKDIR /go/src/github.com/francoggm/go_expenses_api
RUN go mod download
COPY . /go/src/github.com/francoggm/go_expenses_api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/service github.com/francoggm/go_expenses_api
EXPOSE 8080 8080
ENTRYPOINT ["build/service"]

# No-cache. Needs adjust in paths

# FROM alpine
# RUN apk add --no-cache ca-certificates && update-ca-certificates
# COPY --from=builder /go/src/github.com/francoggm/go_expenses_api/build/go_expenses_api /usr/bin/go_expenses_api
# EXPOSE 8080 8080
# ENTRYPOINT ["/usr/bin/go_expenses_api"]