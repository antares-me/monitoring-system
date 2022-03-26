FROM golang:1.17-alpine3.15 AS builder
RUN mkdir /build
ADD ./ /build/
WORKDIR /build
RUN go build -o ./monitoring ./cmd/monitoring-system && go build -o ./generator ./cmd/test-generator

FROM alpine
RUN adduser -S -D -H -h /app appuser
COPY --from=builder /build/generator /build/monitoring /app/
COPY web /app/web/
COPY configs /app/configs/
COPY wrapper.sh /app/
RUN chown -R appuser /app
USER appuser
WORKDIR /app
CMD ["./wrapper.sh", "&"]
EXPOSE 80
