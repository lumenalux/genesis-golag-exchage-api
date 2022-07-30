FROM golang:1.18-alpine

WORKDIR /api

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod tidy

COPY . .

EXPOSE 8081
EXPOSE 587
EXPOSE 80

RUN go build -o exchange-api .

CMD [ "./exchange-api" ]
