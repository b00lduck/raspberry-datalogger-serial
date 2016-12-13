FROM rem/rpi-golang-1.7:latest

WORKDIR /gopath/src/github.com/b00lduck/raspberry-datalogger-serial
CMD ["raspberry-datalogger-serial"]

ADD . /gopath/src/github.com/b00lduck/raspberry-datalogger-serial
RUN go get
RUN go build
