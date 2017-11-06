package main

import (
	"flag"
	"log"
	"time"

	usbrelay "github.com/mengzhuo/go-usb-relay"
)

func main() {
	flag.Parse()
	all := usbrelay.ListAll()
	for id, r := range all {
		log.Printf("%s has %d slots", id, r.SlotNum())
		log.Println(id, "....on")
		r.TurnAllOn()
		log.Print(r.GetAllStatus())
		time.Sleep(time.Second)
		log.Println(id, "....off")
		r.TurnAllOff()
		log.Print(r.GetAllStatus())
		r.Close()
	}
}
