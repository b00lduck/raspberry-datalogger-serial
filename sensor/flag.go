package sensor
import (
	"github.com/b00lduck/raspberry-datalogger-dataservice-client"
	"fmt"
)

type Flag struct {
	oldValue uint8
	code string
}

func NewFlag(code string) Flag {
	return Flag{
		oldValue: 255,
		code: code,
	}
}

func (f *Flag) SetNewState(state uint8) {

	if f.oldValue != state {
		if err := client.SendFlagState(f.code, state); err != nil {
			fmt.Println(err)
		} else {
			f.oldValue = state
		}
	}
}

