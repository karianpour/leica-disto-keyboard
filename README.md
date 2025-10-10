# Leica Disto Laser Meter Keyboard Emulator

Read the measurement from Leica Disto device, using nodejs / go.
I tested it on D110.

## NodeJS
to find the SERVICE_UUID and CHARACTERISTIC_UUID
```bash
node discover.js
```

then put the UUIDs in disto-ble.js file and run it
```bash
node disto-ble.js
```

## Go
To discover
```bash
cd go/discover
go run .
```
```bash
cd go
go run .
```