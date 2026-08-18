FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . /go/src/app

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o stremio-film-festivals .

FROM gcr.io/distroless/static:nonroot

# Copy compiled binary
COPY --from=builder /go/src/app/stremio-film-festivals /stremio-film-festivals

# Copy pre-packaged festival datasets so the container works out-of-the-box
COPY --from=builder /go/src/app/data /data

VOLUME ["/data"]
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/stremio-film-festivals"]
CMD ["-bindAddr", "0.0.0.0", "-dataDir", "/data"]
