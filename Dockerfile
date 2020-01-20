FROM golang:1.13 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/app cmd/stardust/main.go

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /

CMD ["/app"]