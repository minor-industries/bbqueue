package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/minor-industries/bbqueue/db"
	"github.com/minor-industries/bbqueue/schema"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
	"time"
)

func pollUsbSerial() (string, error) {
	polling := false
	for {
		dev, err := filepath.Glob("/dev/tty.usb*")
		if err != nil {
			return "", errors.Wrap(err, "glob")
		}

		switch len(dev) {
		case 1:
			fmt.Println("found", dev[0])
			return dev[0], nil
		case 0:
			if !polling {
				fmt.Println("polling for device")
			}
			polling = true
			time.Sleep(100 * time.Millisecond)
		default:
			return "", errors.New("found more than one serial device")
		}
	}
}

func poll(db *gorm.DB) error {
	dev, err := pollUsbSerial()
	if err != nil {
		return errors.Wrap(err, "poll usb serial")
	}

	device, err := serial.OpenPort(&serial.Config{
		Name:        dev,
		Baud:        115200,
		ReadTimeout: 0,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    1,
	})
	if err != nil {
		return errors.Wrap(err, "open serial port")
	}

	scanner := bufio.NewScanner(device)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		fmt.Println(parts)

		switch parts[0] {
		case "LOG":
			log := strings.TrimPrefix(line, "LOG ")
			_ = log // fmt.Println(log)
		case "RPC":
			args := parts[1:]
			err := rpc(args)
			if err != nil {
				fmt.Println("rpc error", err)
			}
		}

	}

	return scanner.Err()
}

func rpc(args []string) error {
	if len(args) == 0 {
		return errors.New("missing rpc method")
	}

	method := args[0]
	switch method {
	case "RADIO-TX":
		return handleRadioTx(args[1:])
	}

	fmt.Println("rpc", args)
	return nil
}

func handleRadioTx(args []string) error {
	if len(args) == 0 {
		return errors.New("missing payload")
	}
	payload := args[0]
	fmt.Println(payload)

	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return errors.Wrap(err, "base64 decode payload")
	}

	fmt.Println(hex.Dump(raw))

	p := &rfm69.Packet{}
	_, err = p.UnmarshalMsg(raw)
	if err != nil {
		return errors.Wrap(err, "unmarshal packet")
	}

	return processPacket(p)
}

func processPacket(p *rfm69.Packet) error {
	fmt.Println(spew.Sdump(p))

	cmd := p.Payload[0]
	data := p.Payload[1:]

	switch cmd {
	case 0x02:
		fmt.Println(cmd, hex.Dump(data))
		tcData := &schema.ThermocoupleData{}
		err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, tcData)
		if err != nil {
			return errors.Wrap(err, "parse temperature data")
		}
		fmt.Println(
			nullTerminatedBytesToString(tcData.Description[:]),
			tcData.Temperature,
		)
	}

	return nil
}

func run() error {
	db, err := db.Get("sqlite3.db")
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	for {
		err := poll(db)
		fmt.Println("error:", err)
	}
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
