FROM hypriot/rpi-golang
WORKDIR /gopath1.5/src/b00lduck/raspberry-datalogger-serial
CMD ["raspberry-datalogger-serial"]

ADD . /gopath1.5/src/b00lduck/raspberry-datalogger-serial
RUN go get
RUN go build
USER root
