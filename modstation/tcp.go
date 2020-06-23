package modstation

import (
	"fmt"

	"github.com/goburrow/modbus"
)

//TCP exporeted
func TCP() {
	handler := modbus.NewTCPClientHandler("localhost:502")
	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()
	client := modbus.NewClient(handler)

	_, err = client.WriteMultipleRegisters(0, 4, []byte{0, 10, 0, 255, 1, 5, 0, 3})
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	results, err := client.ReadHoldingRegisters(0, 3)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("results %v\n", results)
}
