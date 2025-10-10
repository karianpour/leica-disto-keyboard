package main

import (
	"fmt"
	"log"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	// Enable Bluetooth
	must("enable adapter", adapter.Enable())

	fmt.Println("Scanning for devices...")
	adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		name := result.LocalName()
		if name != "" {
			fmt.Printf("Found device: %s [%s]\n", result.Address.String(), name)
		}
	})
}

func must(action string, err error) {
	if err != nil {
		log.Fatalf("Failed to %s: %v", action, err)
	}
}
