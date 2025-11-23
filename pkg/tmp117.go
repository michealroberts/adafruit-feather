/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package tmp117

/**************************************************************************************/

import (
	"machine"
)

/**************************************************************************************/

const DefaultTMP117I2CAddress uint16 = 0x48

/**************************************************************************************/

// Device represents an Adafruit TMP117 temperature sensor connected over I2C.
type Device struct {
	address uint16      // I2C address of the sensor
	bus     machine.I2C // I2C bus to which the sensor is connected
}

/**************************************************************************************/

// NewDevice creates a new Device instance with the specified I2C address and bus.
func NewDevice(address uint16, bus machine.I2C) *Device {
	return &Device{
		address: address,
		bus:     bus,
	}
}

/**************************************************************************************/
