FROM golang:1.21-alpine3.19 AS builder
WORKDIR /
ADD . .
RUN go install -v github.com/owasp-amass/amass/v4/...@master
RUN go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
RUN go build -o bin/live-sub /cmd/main.go

FROM golang:1.19-alpine3.18 
WORKDIR /
COPY --from=builder /bin/live-sub .
COPY --from=builder /go/bin/amass /usr/bin/amass
COPY --from=builder /go/bin/subfinder /usr/bin/subfinder

ENTRYPOINT ["./live-sub"]
