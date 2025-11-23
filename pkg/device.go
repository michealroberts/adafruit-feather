/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package adafruit

/**************************************************************************************/

import (
	"machine"
)

/**************************************************************************************/

// Device represents an Adafruit device connected over I2C.
type Device struct {
	address uint16      // I2C address of the sensor
	bus     machine.I2C // I2C bus to which the sensor is connected
}

/**************************************************************************************/
