FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . /go/src/app

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o stremio-film-festivals .

FROM gcr.io/distroless/static:nonroot

# Copy compiled binary (includes embedded branding logo and static UI)
COPY --from=builder /go/src/app/stremio-film-festivals /stremio-film-festivals

# Copy pre-packaged festival datasets for 100% stateless multi-replica operation
COPY --from=builder /go/src/app/data /data

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/stremio-film-festivals"]
CMD ["-bindAddr", "0.0.0.0", "-port", "8080", "-dataDir", "/data"]
