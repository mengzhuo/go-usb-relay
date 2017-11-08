package usbrelay

import (
	"errors"
	"fmt"
	"sync"

	"github.com/zserge/hid"
)

type Relay struct {
	dev       hid.Device
	idCache   string
	slotCache int
}

var (
	SlotNumberOverflow = errors.New("slot should be in 1 to 8")
)

func (r *Relay) Id() (id string, err error) {

	if r.idCache != "" {
		return r.idCache, nil
	}

	buf := make([]byte, 8)
	_, err = r.dev.Ctrl(0xa0, 0x01, 3<<8, 0x0, buf, 1000)
	if err != nil || buf[4] >= '8' || buf[4] <= '0' {
		return "", err
	}
	return string(buf[:5]), err
}

func (r *Relay) SlotNum() int {
	if r.slotCache != 0 {
		return r.slotCache
	}
	id, err := r.Id()
	if err != nil || id == "" {
		return 0
	}
	r.slotCache = int(id[len(id)-1] - '0')
	return r.slotCache
}

func (r *Relay) Device() hid.Device {
	return r.dev
}

func (r *Relay) GetAllStatus() (bitmask byte, err error) {
	buf := make([]byte, 8)
	_, err = r.dev.Ctrl(0xa0, 0x01, 3<<8, 0x0, buf, 1000)
	if err != nil {
		return 0, err
	}
	return buf[7], err
}

func (r *Relay) Status(n int) (on bool, err error) {

	if n < 1 || n > 8 {
		return false, SlotNumberOverflow
	}
	mask, err := r.GetAllStatus()
	return (mask&1<<(uint(n)-1) != 0), err
}

func (r *Relay) TurnOn(n int) (err error) {
	if n < 1 || n > 8 {
		return SlotNumberOverflow
	}
	return r.onOff(n, true)
}

func (r *Relay) TurnOff(n int) (err error) {
	if n < 1 || n > 8 {
		return SlotNumberOverflow
	}
	return r.onOff(n, false)
}

func (r *Relay) TurnAllOn() (err error) {
	return r.onOff(0, true)
}

func (r *Relay) TurnAllOff() (err error) {
	return r.onOff(0, false)
}

func (r *Relay) onOff(n int, isOn bool) (err error) {

	buf := make([]byte, 8)
	if n <= 0 {
		if isOn {
			buf[0] = 0xfe
		} else {
			buf[0] = 0xfc
		}
	} else {
		buf[1] = byte(n)
		if isOn {
			buf[0] = 0xff
		} else {
			buf[0] = 0xfd
		}
	}
	err = r.dev.SetReport(0, buf)
	return
}

func (r *Relay) Toggle(n int) (err error) {
	var on bool
	on, err = r.Status(n)
	if err != nil {
		return
	}
	return r.onOff(n, !on)
}

func (r *Relay) Close() {
	r.dev.Close()
}

func ListAll() (ret map[string]*Relay) {
	ret = map[string]*Relay{}
	var lock sync.Mutex

	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		pid := fmt.Sprintf("%04x:%04x", info.Vendor, info.Product)
		if pid != "16c0:05df" {
			return
		}
		if err := device.Open(); err != nil {
			return
		}
		var (
			id  string
			err error
		)
		r := &Relay{dev: device}
		if id, err = r.Id(); err != nil || id == "" {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		ret[id] = r
		return
	})
	return ret
}
