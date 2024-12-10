package main

import (
	"fmt"
	"os"
	"vu/ase/actuator/src/handler"
	"vu/ase/actuator/src/receiver"

	pb_outputs "github.com/VU-ASE/rovercom/packages/go/outputs"
	roverlib "github.com/VU-ASE/roverlib-go/src"

	"github.com/rs/zerolog/log"
)

// The actual program
func run(service roverlib.Service, config *roverlib.ServiceConfiguration) error {
	// Create all necessary queues
	handlerQueue := make(chan *pb_outputs.ControllerOutput, 300) // all incoming messages that need to be processed still

	//
	// Get input streams
	//

	// Get the address of the controller output publisher
	controllerOutput := service.GetReadStream("controller", "decision")

	//
	// Get configuration values
	//
	if config == nil {
		return fmt.Errorf("No configuration provided")
	}

	i2cbus, err := config.GetIntSafe("i2c-bus")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}
	servoTrim, err := config.GetFloatSafe("servo-trim")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}
	servoScaler, err := config.GetFloatSafe("servo-scaler")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}
	enableDiff, err := config.GetIntSafe("electronic-diff")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get electrical differential enable flag")
		return err
	}
	trackWidth, err := config.GetFloatSafe("track-width")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get track-width")
		return err
	}
	fanCap, err := config.GetIntSafe("fan-cap")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get fan-cap")
		return err
	}
	if fanCap < 0 {
		fanCap = 0
	} else if fanCap > 100 {
		fanCap = 100
	}

	// Start all goroutines in a self-restarting way
	go func() {
		for {
			handler.Start(handlerQueue, uint(i2cbus), servoScaler, servoTrim, enableDiff, trackWidth, fanCap)
		}
	}()
	go func() {
		for {
			receiver.Start(*controllerOutput, handlerQueue)
		}
	}()

	// Block forever
	select {}
}

func onTerminate(sig os.Signal) error {
	return nil
}

// Used to start the program with the correct arguments
func main() {
	roverlib.Run(run, onTerminate)
}
