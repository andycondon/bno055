package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andycondon/bno055"
	"github.com/andycondon/bno055/i2c"
)

func main() {
	i2cBus, err := i2c.NewBus(0x28, 1)
	if err != nil {
		panic(err)
	}

	sensor, err := bno055.NewSensorFromBus(i2cBus)
	if err != nil {
		panic(err)
	}

	err = sensor.UseExternalCrystal(true)
	if err != nil {
		panic(err)
	}

	var (
		isCalibrated       bool
		calibrationOffsets bno055.CalibrationOffsets
		calibrationStatus  *bno055.CalibrationStatus
	)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for !isCalibrated {
		select {
		case <-signals:
			err := i2cBus.Close()
			if err != nil {
				panic(err)
			}
		default:
			calibrationOffsets, calibrationStatus, err = sensor.Calibration()
			if err != nil {
				panic(err)
			}

			isCalibrated = calibrationStatus.IsCalibrated()

			fmt.Printf(
				"\r*** Calibration status (0..3): system=%v, accelerometer=%v, gyroscope=%v, magnetometer=%v",
				calibrationStatus.System,
				calibrationStatus.Accelerometer,
				calibrationStatus.Gyroscope,
				calibrationStatus.Magnetometer,
			)
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("*** Done! Calibration offsets: %v\n", calibrationOffsets)

	// Output
	// *** Calibration status (0..3): system=3, accelerometer=3, gyroscope=3, magnetometer=3
	// *** Done! Calibration offsets: [...]
}
