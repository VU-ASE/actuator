package receiver

import (
	"vu/ase/actuator/src/handler"

	roverlib "github.com/VU-ASE/roverlib-go/v2/src"

	"github.com/rs/zerolog/log"
)

// Starts the receiver to listen for incoming messages
func Start(controllerOutput roverlib.ReadStream, handlerQueue handler.Queue) {

	// Main receiver loop
	for {
		msg, err := controllerOutput.Read()
		if err != nil {
			log.Error().Err(err).Msg("Failed to read message")
			continue
		}

		// Get controller data from the message
		controllerData := msg.GetControllerOutput()
		if controllerData == nil {
			log.Error().Msg("Received message without controller data")
			continue
		}

		// Add the message to the queue of outstanding messages (and give up ownership)
		handlerQueue <- controllerData
	}
}
