package parser
import (
	"errors"
	"strconv"
)

type interpol struct {
	t float64
	u float64
}

func ParseDigitalIn(ch uint8, data []byte) uint8 {
	hex := parseHexDigit(32, data) * 16 + parseHexDigit(33, data)
	if (hex & (1 << ch)) > 0 {
		return 1
	}
	return 0
}

// Sensor type A
func ParseADCSensorA(ch int, data []byte) (ret float64) {
	points := []interpol {
		interpol { 20, 2.47 },
		interpol { 25, 2.37 },
		interpol { 30, 2.27 },
		interpol { 35, 2.17 },
		interpol { 40, 2.06 },
		interpol { 45, 1.94 },
		interpol { 50, 1.82 },
		interpol { 55, 1.70 },
		interpol { 60, 1.57 },
		interpol { 65, 1.46 },
		interpol { 70, 1.34 },
		interpol { 75, 1.23 },
		interpol { 80, 1.13 },
	}
	volt := parseADCVolt(ch, data)
	return interpolate(volt, points)
}



func interpolate(u float64, table []interpol) float64 {

	if u >= table[0].u {
		return table[0].t
	}

	for index := 0; index < len(table); index++ {
		if (table[index].u < u) {
			interval := table[index-1].u - table[index].u
			a := table[index-1].u - u
			frac := a / interval

			t_interval := table[index-1].t - table[index].t

			return table[index-1].t - t_interval * frac
		}
	}

	return table[len(table)-1].t

}

// Sensor type B
func ParseADCSensorB(ch int, data []byte) (ret float64) {
	points := []interpol {
		interpol { -20, 4.54 },
		interpol { -15, 4.42 },
		interpol { -10, 4.29 },
		interpol { -5, 4.13 },
		interpol { 0, 3.96 },
		interpol { 5, 3.77 },
		interpol { 10, 3.56 },
		interpol { 15, 3.34 },
		interpol { 20, 3.05 },
	}
	volt := parseADCVolt(ch, data)
	return interpolate(volt, points)
}

// Sensor type C
func ParseADCSensorC(ch int, data []byte) (ret float64) {
	points := []interpol {
		interpol { 20, 2.60 },
		interpol { 25, 2.47 },
		interpol { 30, 2.34 },
		interpol { 35, 2.20 },
		interpol { 40, 2.06 },
		interpol { 45, 1.91 },
		interpol { 50, 1.77 },
		interpol { 55, 1.63 },
		interpol { 60, 1.49 },
		interpol { 65, 1.36 },
		interpol { 70, 1.23 },
		interpol { 75, 1.12 },
		interpol { 80, 1.01 },
	}
	volt := parseADCVolt(ch, data)
	return interpolate(volt, points)
}

func parseADCVolt(ch int, data []byte) (ret float64) {
	return float64(parseADC(ch, data)) * (4.97 / 1024.0)
}

func parseADC(ch int, data []byte) (ret uint16) {
	index := ch * 4
	ret = parseHexDigit(index, data) * 256
	ret += parseHexDigit(index + 1, data) * 16
	ret += parseHexDigit(index + 2, data)
	return
}

func parseHexDigit(index int, data []byte) uint16 {
	val := uint16(data[index])
	if val > 57 {
		return val - 87
	}
	return val - 48
}

func IsSmallHexDigit(data []byte, index int) error {
	c := data[index]
	if c < 48 || c > 52 {
		return errors.New("char at index " + strconv.Itoa(index) + " must be a valid small hex digit (0-3) but was " + string(c))
	}
	return nil
}

func IsHexDigit(data []byte, index int) error {
	c := data[index]
	if c < 48 || (c > 57 && (c < 97 || c > 102)) {
		return errors.New("char at index " + strconv.Itoa(index) + " must be a valid lowercase hex digit (0-9,a-f) but was " + string(c))
	}
	return nil
}
