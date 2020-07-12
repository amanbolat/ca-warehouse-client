# Build
FROM golang:alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
RUN go build -o main ./cmd/main.go

# Create actual image
FROM alpine
RUN adduser -S -D -H -h /app appuser
RUN mkdir -p /app
RUN chown -R appuser /app
USER appuser
COPY --from=builder /build/main /app/
COPY --from=builder /build/fonts/ /fonts/
ENV FONT_PATH /fonts/notosanssc-mono-regular.ttf
ENV PORT 8999
WORKDIR /app
CMD ["./main"]