package main

import (
	"encoding/binary"
	"encoding/json"
	"github.com/goburrow/modbus"
	"log"
	"time"
)

type RoverData struct {
	PVVolts     float32
	PVAmps      float32
	ChargeVolts float32
	ChargeAmps  float32
	ChargeMode  string
	Timestamp   int64
}

var chargeModes = []string{"deactivated", "activated", "mppt", "equalizing", "boost", "floating", "current limiting"}

func main() {
	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second
	var d RoverData

	err := handler.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)
	for {

		results, err := client.ReadHoldingRegisters(288, 1)
		if err != nil {
			log.Fatal(err)
		}
		d.ChargeMode = chargeModes[int(results[1])]

		results, err = client.ReadHoldingRegisters(264, 2)
		if err != nil {
			log.Fatal(err)
		}
		tmp := binary.BigEndian.Uint16(results[0:])
		d.PVAmps = float32(tmp) / 100

		results, err = client.ReadHoldingRegisters(263, 2)
		if err != nil {
			log.Fatal(err)
		}
		tmp = binary.BigEndian.Uint16(results[0:])
		d.PVVolts = float32(tmp) / 10

		results, err = client.ReadHoldingRegisters(0x102, 2)
		if err != nil {
			log.Fatal(err)
		}
		tmp = binary.BigEndian.Uint16(results[0:])
		d.ChargeAmps = float32(tmp) / 100

		results, err = client.ReadHoldingRegisters(0x101, 2)
		if err != nil {
			log.Fatal(err)
		}
		tmp = binary.BigEndian.Uint16(results[0:])
		d.ChargeVolts = float32(tmp) / 10

		d.Timestamp = time.Now().Unix()
		j, err := json.Marshal(d)
		log.Printf("%s\n", j)
	}
}
