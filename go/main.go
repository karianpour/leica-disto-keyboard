package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/micmonay/keybd_event"
	"tinygo.org/x/bluetooth"
)

var (
	adapter             = bluetooth.DefaultAdapter
	SERVICE_UUID        = bluetooth.NewUUID([16]byte{0x3a, 0xb1, 0x01, 0x00, 0xf8, 0x31, 0x43, 0x95, 0xb2, 0x9d, 0x57, 0x09, 0x77, 0xd5, 0xbf, 0x94})
	CHARACTERISTIC_UUID = bluetooth.NewUUID([16]byte{0x3a, 0xb1, 0x01, 0x01, 0xf8, 0x31, 0x43, 0x95, 0xb2, 0x9d, 0x57, 0x09, 0x77, 0xd5, 0xbf, 0x94})
)

func main() {
	must("enable adapter", adapter.Enable())

	fmt.Println("🔍 Scanning for Leica DISTO...")
	var found bluetooth.ScanResult
	ch := make(chan bluetooth.ScanResult)

	go func() {
		adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			name := result.LocalName()
			if name != "" {
				fmt.Printf("Found: %s [%s]\n", result.Address.String(), name)
			}
			if len(name) >= 5 && name[:5] == "DISTO" {
				ch <- result
			}
		})
	}()

	select {
	case found = <-ch:
		fmt.Printf("✅ Found device: %s [%s]\n", found.Address.String(), found.LocalName())
	case <-time.After(20 * time.Second):
		log.Fatal("Timeout: Disto not found")
	}

	adapter.StopScan()

	fmt.Println("Connecting to Leica DISTO...")
	device, err := adapter.Connect(found.Address, bluetooth.ConnectionParams{})
	must("connect", err)
	defer device.Disconnect()
	fmt.Println("🔗 Connected!")

	fmt.Println("Discovering the services...")
	services, err := device.DiscoverServices(nil)
	must("discover services", err)

	var characteristic *bluetooth.DeviceCharacteristic

	for _, s := range services {
		if s.UUID() == SERVICE_UUID {
			fmt.Println("🧩 Found service:", s.UUID())
			chars, err := s.DiscoverCharacteristics(nil)
			must("discover characteristics", err)
			for _, c := range chars {
				fmt.Println("Characteristic:", c.UUID())
				if c.UUID() == CHARACTERISTIC_UUID {
					characteristic = &c
					fmt.Println("📡 Found measurement characteristic")
				}
			}
		}
	}

	if characteristic == nil {
		log.Fatal("❌ Measurement characteristic not found")
	}

	// Subscribe for notifications
	err = characteristic.EnableNotifications(handleMeasurement)

	// err = characteristic.EnableNotifications(func(buf []byte) {
	// 	if len(buf) >= 4 {
	// 		bits := binary.LittleEndian.Uint32(buf[:4])
	// 		value := math.Float32frombits(bits)
	// 		fmt.Printf("📏 %.3f m\n", value)
	// 	} else {
	// 		fmt.Printf("Raw data: %x\n", buf)
	// 	}
	// })

	// err = characteristic.EnableNotifications(func(buf []byte) {
	// 	if len(buf) >= 4 {
	// 		val := binary.LittleEndian.Uint32(buf[:4])
	// 		fmt.Printf("📏 %.3f m\n", float64(val)/1000.0)
	// 	} else {
	// 		fmt.Printf("Data: %x\n", buf)
	// 	}
	// })
	must("subscribe", err)

	fmt.Println("🕓 Waiting for measurements... (press Ctrl+C to stop)")
	select {}
}

func must(action string, err error) {
	if err != nil {
		log.Fatalf("❌ Failed to %s: %v", action, err)
	}
}

func handleMeasurement(buf []byte) {
	if len(buf) < 4 {
		fmt.Printf("Raw data: %x\n", buf)
		return
	}

	// Decode as little-endian float32
	bits := binary.LittleEndian.Uint32(buf[:4])
	value := math.Float32frombits(bits)

	fmt.Printf("📏 %.3f m\n", value)

	// Simulate typing the value and press Enter
	typeValue(fmt.Sprintf("%.3f", value))
}

func typeValue(s string) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// Required on some systems (esp. Mac)
	time.Sleep(100 * time.Millisecond)

	for _, ch := range s {
		switch ch {
		case '.':
			kb.SetKeys(keybd_event.VK_KeypadDecimal) // Decimal point
		case '-':
			kb.SetKeys(keybd_event.VK_MINUS)
		default:
			if ch == '0' {
				kb.SetKeys(keybd_event.VK_0)
			} else if ch == '1' {
				kb.SetKeys(keybd_event.VK_1)
			} else if ch == '2' {
				kb.SetKeys(keybd_event.VK_2)
			} else if ch == '3' {
				kb.SetKeys(keybd_event.VK_3)
			} else if ch == '4' {
				kb.SetKeys(keybd_event.VK_4)
			} else if ch == '5' {
				kb.SetKeys(keybd_event.VK_5)
			} else if ch == '6' {
				kb.SetKeys(keybd_event.VK_6)
			} else if ch == '7' {
				kb.SetKeys(keybd_event.VK_7)
			} else if ch == '8' {
				kb.SetKeys(keybd_event.VK_8)
			} else if ch == '8' {
				kb.SetKeys(keybd_event.VK_8)
			} else if ch == '9' {
				kb.SetKeys(keybd_event.VK_9)
			} else {
				fmt.Printf("Unsupported char: %c\n", ch)
				continue
			}
		}
		kb.Launching()
		time.Sleep(30 * time.Millisecond)
	}
	// Press Enter
	kb.SetKeys(keybd_event.VK_ENTER)
	kb.Launching()
}
