package main

import (
	"os"
	"vu/ase/actuator/src/handler"
	"vu/ase/actuator/src/receiver"

	pb_core_messages "github.com/VU-ASE/rovercom/packages/go/core"
	pb_module_outputs "github.com/VU-ASE/rovercom/packages/go/outputs"

	roverlib "github.com/VU-ASE/roverlib/src"
	"github.com/rs/zerolog/log"
)

// The actual program
func run(service roverlib.ResolvedService, sysmanInfo roverlib.CoreInfo, tuningState *pb_core_messages.TuningState) error {
	// Create all necessary queues
	handlerQueue := make(chan *pb_module_outputs.ControllerOutput, 300) // all incoming messages that need to be processed still

	// Get the address of the controller output publisher
	controllerOutputAddr, err := service.GetDependencyAddress("controller", "decision")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get address of controller output publisher")
		return err
	}

	// Get the I2C bus number
	i2cbus, err := roverlib.GetTuningInt("i2c-bus", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}

	servoTrim, err := roverlib.GetTuningFloat("servo-trim", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}

	servoScaler, err := roverlib.GetTuningFloat("servo-scaler", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}

	enableDiff, err := roverlib.GetTuningInt("electronic-diff", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get electrical differential enable flag")
		return err
	}

	trackWidth, err := roverlib.GetTuningFloat("track-width", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get track-width")
		return err
	}

	fanCap, err := roverlib.GetTuningInt("fan-cap", tuningState)
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
			receiver.Start(controllerOutputAddr, handlerQueue)
		}
	}()

	// Block forever
	select {}
}

func onTuningState(tuningState *pb_core_messages.TuningState) {
	log.Warn().Msg("Received new tuning state")

	servoTrim, err := roverlib.GetTuningFloat("servo-trim", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return
	}
	handler.SetServoTrim(servoTrim)

	servoScaler, err := roverlib.GetTuningFloat("servo-scaler", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return
	}
	handler.SetServoScaler(servoScaler)
}

func onTerminate(sig os.Signal) {
	log.Info().Msg("Received signal, terminating")
	handler.OnTerminate()
}

// Used to start the program with the correct arguments
func main() {
	roverlib.Run(run, onTuningState, onTerminate, false)
}
