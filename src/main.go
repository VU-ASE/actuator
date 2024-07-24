package main

import (
	"os"
	"vu/ase/actuator/src/handler"
	"vu/ase/actuator/src/receiver"

	pb_module_outputs "github.com/VU-ASE/pkg-CommunicationDefinitions/v2/packages/go/outputs"
	pb_systemmanager_messages "github.com/VU-ASE/pkg-CommunicationDefinitions/v2/packages/go/systemmanager"

	servicerunner "github.com/VU-ASE/pkg-ServiceRunner/v2/src"
	"github.com/rs/zerolog/log"
)

// The actual program
func run(service servicerunner.ResolvedService, sysmanInfo servicerunner.SystemManagerInfo, tuningState *pb_systemmanager_messages.TuningState) error {
	// Create all necessary queues
	handlerQueue := make(chan *pb_module_outputs.ControllerOutput, 300) // all incoming messages that need to be processed still

	// Get the address of the controller output publisher
	controllerOutputAddr, err := service.GetDependencyAddress("controller", "decision")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get address of controller output publisher")
		return err
	}

	// Get the I2C bus number
	i2cbus, err := servicerunner.GetTuningInt("i2c-bus", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}

	servoTrim, err := servicerunner.GetTuningFloat("servo-trim", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}

	servoScaler, err := servicerunner.GetTuningFloat("servo-scaler", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return err
	}

	enableDiff, err := servicerunner.GetTuningInt("electronic-diff", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get electrical differential enable flag")
		return err
	}

	trackWidth, err := servicerunner.GetTuningFloat("track-width", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get track-width")
		return err
	}

	fanCap, err := servicerunner.GetTuningInt("fan-cap", tuningState)
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

func onTuningState(tuningState *pb_systemmanager_messages.TuningState) {
	log.Warn().Msg("Received new tuning state")

	servoTrim, err := servicerunner.GetTuningFloat("servo-trim", tuningState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get I2C bus number")
		return
	}
	handler.SetServoTrim(servoTrim)

	servoScaler, err := servicerunner.GetTuningFloat("servo-scaler", tuningState)
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
	servicerunner.Run(run, onTuningState, onTerminate, false)
}
