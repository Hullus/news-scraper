FROM golang:1.23
LABEL authors="None of your beez wax"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

ENTRYPOINT ["/main"]