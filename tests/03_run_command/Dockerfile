FROM docker.io/library/golang:latest AS builder

WORKDIR /proj
COPY ./go.mod ./main.go ./
RUN go build -o prog

FROM scratch

WORKDIR /proj
COPY --from=builder /proj/prog /proj/prog

ENTRYPOINT ["/proj/prog"]
