package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zserge/hid"
)

var (
	onOff = flag.Bool("s", false, "relay status")
	rNum  = flag.Int("n", 1, "relay number 1-8")
)

func main() {
	flag.Parse()
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		id := fmt.Sprintf("%04x:%04x:%04x:%02x", info.Vendor, info.Product, info.Revision, info.Interface)
		if !strings.HasPrefix(id, "16c0:05df") {
			return
		}
		log.Print("id: ", id)
		if err := device.Open(); err != nil {
			log.Println("Open error: ", err)
			return
		}

		go func() {
			for {
				if buf, err := device.Read(-1, 1*time.Second); err == nil {
					log.Println("Input report:  ", hex.EncodeToString(buf))
				}
			}
		}()

		relOnOff(device, *rNum, *onOff)
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
			cmd1 = 0xfd
		}
	}

	buf[0] = cmd1
	buf[1] = cmd2
	log.Printf("num [%d] isOn=%v buf=%x", num, isOn, buf[:8])
	err := dev.SetReport(0, buf[:8])
	if err != nil {
		log.Print(err)
	}
}
