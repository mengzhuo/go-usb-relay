package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zserge/hid"
)

func main() {
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		id := fmt.Sprintf("%04x:%04x:%04x:%02x", info.Vendor, info.Product, info.Revision, info.Interface)
		if !strings.HasPrefix(id, "16c0:05df") {
			return
		}
		if err := device.Open(); err != nil {
			log.Println("Open error: ", err)
			return
		}
		defer device.Close()

		for i := 1; i <= 8; i++ {
			relOnOff(device, i, true)
			time.Sleep(1 * time.Second)
			relOnOff(device, i, false)
			time.Sleep(1 * time.Second)
		}
	})
}

func relOnOff(dev hid.Device, num int, isOn bool) {

	buf := make([]byte, 10)

	var cmd1, cmd2 byte
	if num < 0 && -num <= 8 {
		cmd2 = 0

		if isOn {
			cmd1 = 0xfe
			//maskVal = (1 << uint(-num)) - 1
		} else {
			cmd1 = 0xfc
		}
	} else {
		if num <= 0 || num > 8 {
			log.Print("relay num must be 1-8")
			return
		}

		cmd2 = byte(num)

		if isOn {
			cmd1 = 0xff
		} else {
			cmd1 = 0xf0
		}
	}

	buf[0] = 0
	buf[1] = cmd1
	buf[2] = cmd2
	log.Printf("num [%d] isOn=%v buf=%x", num, isOn, buf[:9])
	err := dev.SetReport(0, buf[:9])
	if err != nil {
		log.Print(err)
	}
}
