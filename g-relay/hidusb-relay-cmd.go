package main

import (
	"flag"
	"log"

	usbrelay "github.com/mengzhuo/go-usb-relay"
)

var (
	numS  = flag.Int("n", 1, "switch on slot number")
	onOff = flag.Bool("o", false, "on or off")
)

func main() {
	flag.Parse()
	log.Printf("n=%d, o=%v", *numS, *onOff)
	all := usbrelay.ListAll()
	for id, r := range all {
		log.Printf("%s has %d slots", id, r.SlotNum())
		if *onOff {
			r.TurnOn(*numS)
		} else {
			r.TurnOff(*numS)
		}
		log.Print(r.GetAllStatus())
	}
}
