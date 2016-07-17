package main

import (
	"fmt"
	"time"
	"github.com/tarm/serial"
	"log"
	"errors"
	"strconv"
	"b00lduck/raspberry-datalogger-serial/parser"
	"b00lduck/raspberry-datalogger-serial/sensor"
)

func main() {

    c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, ReadTimeout: time.Second * 2}
    s, err := serial.OpenPort(c)

    if err != nil {
        log.Fatal(err)
    }
    
    for {
		err = requestDatagram(s)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(5 * time.Second)
		}
    }
}

func requestDatagram(s *serial.Port) error {
    
    n, err := s.Write([]byte("DUMP\r"))
    if err != nil {
        return err
    }

    buf := make([]byte, 38) // Datagram length is 38 bytes incl. \n
    n, err = s.Read(buf)
    if err != nil {
        return err
    }

    // check for correct size
    if n != 38 {
	return errors.New("received datagram with invalid size (must: 38, was: " + strconv.Itoa(n) + ")")
    }

    return processDatagram(buf)
}

var thermBrauchwasser = sensor.NewThermometer("HEIZ_BRAUCHW", 0.25)
var thermAussen = sensor.NewThermometer("HEIZ_AUSSEN", 0.50)
var thermKessel = sensor.NewThermometer("HEIZ_KESSEL", 0.50)

var flagZirkulationspumpe = sensor.NewFlag("HEIZ_BWZP")
var flagLadepumpe = sensor.NewFlag("HEIZ_BWLP")
var flagWinter = sensor.NewFlag("HEIZ_WINT")
var flagUmwaelzpumpe = sensor.NewFlag("HEIZ_UWP")
var flagBrenner = sensor.NewFlag("HEIZ_BRENN")
var flagTag = sensor.NewFlag("HEIZ_TAG")

func processDatagram(data []byte) error {

    if data[37] != 10 {
		return errors.New("last char in datagram must be newline (0x0a)")
    }

    // Check ADC values
    for i:=0;i<8;i++ {
	index := 4 * i
	// check first ADC digit (0,1,2,3)
	if err := parser.IsSmallHexDigit(data,index); err != nil {
	    return err	
	}

	// check second and third ADC digit (0-9,a-f)
	for i:=0;i<2;i++ {
	    index += 1
	    if err := parser.IsHexDigit(data,index); err != nil {
			return err
	    }
	}

	// check spaces between values
	index += 1
        if data[index] != 32 {
    	    return errors.New("char at index " + strconv.Itoa(index) + " must be a space (0x20)")
		}
    }

    // check DIGITAL values
    if data[31] != 32 {
        return errors.New("char at index 31 must be a space (0x20)")
    }

    if err := parser.IsHexDigit(data,32); err != nil {
		return err
    }

    if err := parser.IsHexDigit(data,33); err != nil {
		return err
    }

    // check CRC value (format only, real CRC check later)
    if data[34] != 32 {
        return errors.New("char at index 34 must be a space (0x20)")
    }

    if err := parser.IsHexDigit(data,35); err != nil {
		return err
    }

    if err := parser.IsHexDigit(data,36); err != nil {
		return err
    }

	// TODO: CRC check

	thermBrauchwasser.SetNewReading(parser.ParseADCSensorC(5, data))
    thermAussen.SetNewReading(parser.ParseADCSensorB(6, data))
	thermKessel.SetNewReading(parser.ParseADCSensorA(7, data))

	flagZirkulationspumpe.SetNewState(parser.ParseDigitalIn(0, data))
	flagLadepumpe.SetNewState(parser.ParseDigitalIn(1, data))
	flagWinter.SetNewState(parser.ParseDigitalIn(2, data))
	flagUmwaelzpumpe.SetNewState(parser.ParseDigitalIn(3, data))
	flagBrenner.SetNewState(parser.ParseDigitalIn(4, data))
	flagTag.SetNewState(parser.ParseDigitalIn(5, data))

    return nil
}


