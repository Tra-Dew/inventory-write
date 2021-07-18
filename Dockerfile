FROM golang:1.14.2-alpine as builder

RUN grep nobody /etc/passwd > /etc/passwd.nobody \
    && grep nobody /etc/group > /etc/group.nobody \
    && apk --no-cache update \
    && apk add --no-cache git wget


COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory-write

FROM scratch
WORKDIR /

COPY --from=builder /etc/group.nobody /etc/group
COPY --from=builder /etc/passwd.nobody /etc/passwd

USER nobody

COPY --from=builder /app/inventory-write .


EXPOSE 9001
ENTRYPOINT ["/inventory-write"]
