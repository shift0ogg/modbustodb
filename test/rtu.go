package test

import (
	"fmt"
	"time"
	"xxnet/utils"

	"github.com/goburrow/modbus"
)

var handler *modbus.RTUClientHandler = nil

//Rtu exported
func Rtu(addr byte) ([]byte, error) {
	handler = modbus.NewRTUClientHandler("COM1")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	handler.SlaveId = 8
	adu, err := client.ReadHoldingRegisters(1501, 4)

	for i := 0; i < 2; i++ {
		handler.SlaveId = 8
		adu, err = client.ReadHoldingRegisters(1501, 4)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(utils.BytesToShort(adu[0:2]))
			fmt.Println(utils.BytesToShort(adu[2:4]))
			fmt.Println(utils.BytesToShort(adu[4:6]))
			fmt.Println(utils.BytesToShort(adu[6:8]))
		}

		handler.SlaveId = 1
		adu, err = client.ReadHoldingRegisters(1501, 4)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(utils.BytesToShort(adu[0:2]))
			fmt.Println(utils.BytesToShort(adu[2:4]))
			fmt.Println(utils.BytesToShort(adu[4:6]))
			fmt.Println(utils.BytesToShort(adu[6:8]))
		}
	}

	return adu, err
}

//Rtu1 exported
func Rtu1() ([]byte, error) {
	handler := modbus.NewRTUClientHandler("COM1")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 8
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	adu, err := client.ReadHoldingRegisters(1501, 4)

	if err != nil {
	} else {
		fmt.Println(utils.BytesToShort(adu[0:2]))
		fmt.Println(utils.BytesToShort(adu[2:4]))
		fmt.Println(utils.BytesToShort(adu[4:6]))
		fmt.Println(utils.BytesToShort(adu[6:8]))
	}

	handler.Close()

	return adu, err
}
