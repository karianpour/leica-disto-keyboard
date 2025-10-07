import noble from "@abandonware/noble";
import robot from "robotjs";

// 🧩 replace these with your known UUIDs (no dashes, lowercase)
const SERVICE_UUID = "3ab10100f8314395b29d570977d5bf94"; // e.g. "ffe0"
const CHARACTERISTIC_UUID = "3ab10101f8314395b29d570977d5bf94"; // e.g. "ffe1"
const scale = 3;

async function connectDisto() {
  console.log("🔍 Scanning for Leica DISTO D110...");

  noble.on("stateChange", async (state) => {
    if (state === "poweredOn") {
      await noble.startScanningAsync([SERVICE_UUID], false);
      console.log("Scanning started...");
    } else {
      await noble.stopScanningAsync();
      console.log("Bluetooth not powered on.");
    }
  });

  noble.on("discover", async (peripheral) => {
    const name = peripheral.advertisement.localName || "";
    if (!name.toLowerCase().includes("disto")) return;

    console.log(`📱 Found device: ${name} (${peripheral.address})`);
    await noble.stopScanningAsync();

    try {
      await peripheral.connectAsync();
      console.log("✅ Connected to DISTO");

      const { characteristics } =
        await peripheral.discoverSomeServicesAndCharacteristicsAsync(
          [SERVICE_UUID],
          [CHARACTERISTIC_UUID]
        );

      const measurementChar = characteristics[0];

      if (!measurementChar) {
        console.error("❌ Measurement characteristic not found");
        return;
      }

      console.log("📡 Subscribing to measurement notifications...");
      await measurementChar.subscribeAsync();

      measurementChar.on("data", (data) => {
        const distance = data.readFloatLE(0).toFixed(scale);
        console.log("📏 Measurement:", distance);
        robot.typeString(distance.toString());
        robot.keyTap("enter");
      });

      console.log("✅ Ready! Trigger a measurement on your DISTO.");
    } catch (err) {
      console.error("⚠️ Connection error:", err);
    }
  });
}

connectDisto();
