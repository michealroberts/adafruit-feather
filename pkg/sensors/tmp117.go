/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package sensors

/**************************************************************************************/

import (
	"machine"

	"github.com/michealroberts/adafruit-feather/pkg/device"
)

/**************************************************************************************/

const DefaultTMP117I2CAddress uint16 = 0x48

/**************************************************************************************/

const TMP117LSBResolution float64 = 0.0078125

/**************************************************************************************/

type TMP117Device struct {
	device.Device
}

/**************************************************************************************/

// NewTMP117Device creates a new Device instance for TMP117 sensor using the default I2C
// address and the specified bus.
func NewTMP117Device(bus machine.I2C) *TMP117Device {
	return &TMP117Device{
		Device: device.Device{
			Address: DefaultTMP117I2CAddress,
			Bus:     bus,
		},
	}
}

/**************************************************************************************/

// Read the temperature from the TMP117 sensor (in S.I. Kelvin units).
func (d *TMP117Device) ReadTemperature() (float64, error) {
	// Make 2-bytes (16 bits) buffer for reading out the temperature:
	buffer := make([]byte, 2)

	if err := d.Bus.Tx(d.Address, []byte{0x00}, buffer); err != nil {
		return 0.0, err
	}

	// Convert the 2-byte buffer to a temperature value:
	value := int16(buffer[0])<<8 | int16(buffer[1])

	// Convert the raw temperature value to Kelvin, using the Celsius resolution of each
	// LSB in the TMP117 reading (0.0078125 °C/LSB):
	temperature := (float64(value) * TMP117LSBResolution) + 273.15

	return temperature, nil
}

/**************************************************************************************/
