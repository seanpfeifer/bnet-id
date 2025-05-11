#---------------
#--- Builder ---
#---------------
FROM docker.io/library/golang:1.24-alpine AS builder

COPY [".", "/bnet"]
WORKDIR /bnet
# Make sure you turn off cgo here, because the distroless/static image doesn't have glibc
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

#-----------------------
#--- Resulting image ---
#-----------------------
FROM gcr.io/distroless/static:nonroot

# Mount your secret credentials file in here
VOLUME /secrets
USER nonroot
WORKDIR /

# Copy the actual built binary
COPY --from=builder /server /bnet/
# Also copy our static files that will be served
COPY ["./static", "/static"]

ENTRYPOINT ["/bnet/server"]
