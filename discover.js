const noble = require("@abandonware/noble");

async function startScanning() {
  noble.on("stateChange", async (state) => {
    if (state === "poweredOn") {
      console.log("Scanning for DISTO D110...");
      await noble.startScanningAsync([], false); // [] = all devices, or add UUID filter
    } else {
      noble.stopScanningAsync();
    }
  });

  noble.on("discover", async (peripheral) => {
    const name = peripheral.advertisement.localName;
    console.log("Discovered:", name);

    if (!name) return;

    if (name.includes("DISTO") || name.includes("D110")) {
      console.log("Found DISTO device:", peripheral.address, name);
      await noble.stopScanningAsync();
      connectToDisto(peripheral);
    }
  });
}

async function connectToDisto(device) {
  try {
    await device.connectAsync();
    console.log("Connected to DISTO");

    const { services, characteristics } =
      await device.discoverSomeServicesAndCharacteristicsAsync([], []);
    // If you have known UUIDs for service/characteristic you can filter here

    // Print discovered services & characteristics to inspect them
    console.log(
      "Services:",
      services.map((s) => s.uuid)
    );
    console.log(
      "Characteristics:",
      characteristics.map((c) => c.uuid)
    );

    // Suppose we find the correct characteristic for measurement, call it measureChar
    // (you must inspect/experiment to get correct UUID)
    const measureChar = characteristics.find(
      (c) => /* some condition, e.g. c.uuid === 'xxxx' */ false
    );

    if (measureChar) {
      await measureChar.subscribeAsync();
      measureChar.on("data", (data, isNotification) => {
        // data is a Buffer
        const value = parseMeasurement(data);
        console.log("Measurement:", value);
        // do something with it
      });
    } else {
      console.log("Measurement characteristic not found");
    }
  } catch (err) {
    console.error("Error connecting to DISTO:", err);
  }
}

function parseMeasurement(buffer) {
  // You’ll need to inspect what comes back: maybe ASCII? maybe float? maybe little-endian?
  // Example:
  try {
    // If it's ASCII like "12.34", buffer might contain string
    const s = buffer.toString("utf8");
    const num = parseFloat(s);
    if (!isNaN(num)) return num;
  } catch (e) {}
  // fallback: interpret bytes
  return buffer.readFloatLE ? buffer.readFloatLE(0) : null;
}

startScanning().catch(console.error);
