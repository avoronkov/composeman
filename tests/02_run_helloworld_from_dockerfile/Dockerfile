FROM docker.io/library/golang:latest AS builder

WORKDIR /proj
COPY ./go.mod ./main.go ./
RUN go build -o helloworld

FROM scratch

WORKDIR /proj
COPY --from=builder /proj/helloworld /proj/helloworld

CMD ["./helloworld"]
