FROM golang:1.26

WORKDIR /backend

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/main.go ./

RUN go build -v -o app .

CMD ["./app"]