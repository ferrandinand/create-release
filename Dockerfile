FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o create-release .

FROM alpine
COPY --from=builder /build/create-release /app/
WORKDIR /app
CMD ["./create-release"]