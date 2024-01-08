package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"path/filepath"
	"strings"
	"time"
)

func pollUsbSerial() string {
	polling := false
	for {
		dev, err := filepath.Glob("/dev/tty.usb*")
		noErr(err)

		switch len(dev) {
		case 1:
			fmt.Println("found", dev[0])
			return dev[0]
		case 0:
			if !polling {
				fmt.Println("polling for device")
			}
			polling = true
			time.Sleep(100 * time.Millisecond)
		default:
			panic(errors.New("found more than one serial device"))
		}
	}
}

func run() error {
	dev := pollUsbSerial()

	device, err := serial.OpenPort(&serial.Config{
		Name:        dev,
		Baud:        115200,
		ReadTimeout: 0,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    1,
	})
	noErr(err)

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

	return nil
}

func main() {
	for {
		err := run()
		fmt.Println("error:", err)
	}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
