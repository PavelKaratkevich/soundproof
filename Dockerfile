FROM golang:1.20

WORKDIR /app

COPY . ./
COPY db/migration ./db/migration

RUN go mod download && go mod verify
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go
RUN git clone https://github.com/vishnubob/wait-for-it.git

CMD ["./app"]