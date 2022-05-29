FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache musl-dev gcc build-base ca-certificates

WORKDIR /src

# Naively copy everything. The final binary will be copied into a scratch container.
COPY . .

RUN go build -ldflags "-linkmode external -extldflags \"-static\" -s -w $LDFLAGS" -o the-binary cmd/bot/main.go

# Copy the binary from the "builder" docker target into a scratch container
# to vastly reduce the overall size of the image
FROM scratch AS final

EXPOSE 8000
ENTRYPOINT ["/the-binary"]
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /src/the-binary /the-binary
