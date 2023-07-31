FROM golang:1.20

WORKDIR /app

COPY . ./

RUN go mod download && go mod verify
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

CMD ["./app"]