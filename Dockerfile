FROM golang:1.14 as builder

WORKDIR /src
COPY . .

RUN go build -o /bin/stardust cmd/stardust/main.go

FROM gcr.io/distroless/base
COPY --from=builder /bin/stardust /bin/stardust

ENTRYPOINT ["/bin/stardust"]
