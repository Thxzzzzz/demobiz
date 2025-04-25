FROM golang:1.23 as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace
COPY ./ /workspace
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o demobiz main.go

FROM alpine:3.21.3
WORKDIR /
COPY --from=builder workspace/demobiz ./demobiz

ENTRYPOINT ["/demobiz"]