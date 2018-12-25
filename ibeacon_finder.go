package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type iBeacon struct {
	uuid  string
	major uint16
	minor uint16
}

func NewiBeacon(data []byte) (*iBeacon, error) {
	if len(data) < 25 || binary.BigEndian.Uint32(data) != 0x4c000215 {
		return nil, errors.New("Not an iBeacon")
	}
	fmt.Println("Head: ", binary.BigEndian.Uint32(data))
	beacon := new(iBeacon)
	beacon.uuid = strings.ToUpper(hex.EncodeToString(data[4:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:12]) + "-" + hex.EncodeToString(data[12:14]) + "-" + hex.EncodeToString(data[14:20]))
	beacon.major = binary.BigEndian.Uint16(data[20:22])
	beacon.minor = binary.BigEndian.Uint16(data[22:24])
	return beacon, nil
}

func onPerhipheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	// fmt.Println(string(a.ManufacturerData))
	b, err := NewiBeacon(a.ManufacturerData)
	if err == nil {
		fmt.Println("UUID: ", b.uuid)
		fmt.Println("Major: ", b.major)
		fmt.Println("Minor: ", b.minor)
		fmt.Println("RSSI: ", rssi)
	}
}

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	// unless StatePowerdOn is true, we wont scan.
	case gatt.StatePoweredOn:
		// scan all the uuids
		// you can specify uuid if you want.
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func main() {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
	}
	device.Handle(gatt.PeripheralDiscovered(onPerhipheralDiscovered))
	device.Init(onStateChanged)
	select {}
}
