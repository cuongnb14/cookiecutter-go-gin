FROM golang:1.22 AS builder

WORKDIR /app

ENV REFRESHED_AT 2023-12-25

# RUN apt-get update -qq
# RUN DEBIAN_FRONTEND=noninteractive apt-get install -y upx git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/http cmd/http/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/worker cmd/worker/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/devcli cmd/devcli/main.go


FROM alpine:latest
WORKDIR /app

RUN wget https://release.ariga.io/atlas/atlas-linux-amd64-latest -O /usr/local/bin/atlas \
        && chmod +x /usr/local/bin/atlas

COPY --from=builder /build/ /app/
COPY env/ env/
COPY migrations/ migrations/

COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh
EXPOSE 8080

CMD ["./entrypoint.sh"]
