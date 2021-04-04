FROM golang
WORKDIR /go/src/app
#RUN go get -u github.com/calebwilliams-datastax/papertrader-api
COPY . /go/src/app
RUN go mod download
CMD [ "go","run","main.go"]