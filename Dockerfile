# syntax=docker/dockerfile:1

FROM golang:1.16-alpine
WORKDIR $GOPATH/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /scraper
CMD [ "/scraper" ]
