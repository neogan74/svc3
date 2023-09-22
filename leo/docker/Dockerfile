FROM golang:1.21 AS build_svc3
ENV CGO_ENABLED 0
ARG VCS_REF
ARG BUILD_REF

# Copy source code
COPY . /service

# Build the svc3 binary
WORKDIR /service
RUN go build -ldflags "-X main.build=${BUILD_REF}"
CMD ["ls -al ./"]
 

# Run the Go binary in Alpine 
FROM alpine:latest
ARG BUILD_DATE
ARG BUILD_REF
ARG VCS_REF

COPY --from=build_svc3 /service/svc3 /service/svc3 
WORKDIR /service
CMD ["./svc3"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="svc3" \
      org.opencontainers.image.authors="Andrei Gushchin <neogan74@gmail.com>" \
      org.opencontainers.image.source="${VCS_REF}" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendors="neogan74" 