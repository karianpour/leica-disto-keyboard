# Leica Disto Laser Meter Keyboard Emulator

Read the measurement from Leica Disto device, using nodejs / go.
I tested it on D110.

to find the SERVICE_UUID and CHARACTERISTIC_UUID
```js
node discover.js
```

then put the UUIDs in disto-ble.js file and run it
```js
node disto-ble.js
```