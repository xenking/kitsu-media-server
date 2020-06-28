FROM golang:alpine
RUN mkdir /kitsu-media-server
WORKDIR /kitsu-media-server

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd/
COPY pkg ./pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/kitsu-media-server cmd/main.go
ENTRYPOINT ["/go/bin/kitsu-media-server"]

