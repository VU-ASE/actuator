package receiver

import (
	"vu/ase/actuator/src/handler"

	pb_module_outputs "github.com/VU-ASE/rovercom/packages/go/outputs"

	zmq "github.com/pebbe/zmq4"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// Starts the receiver to listen for incoming messages
func Start(address string, handlerQueue handler.Queue) {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()

	err := subscriber.Connect(address)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to zmq address")
		return
	}

	err = subscriber.SetSubscribe("") // Subscribe to all messages
	if err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to address")
		return
	}

	// Main receiver loop
	for {
		msg, err := subscriber.RecvBytes(0)
		// Don't exit on errors but log them
		if err != nil {
			log.Err(err).Msg("Error while receiving message")
			continue
		}

		// Decode the sensor message
		sensorMsg := pb_module_outputs.SensorOutput{}
		err = proto.Unmarshal(msg, &sensorMsg)
		if err != nil {
			log.Err(err).Msg("Error while decoding message")
			continue
		}

		// Get controller data from the message
		controllerData := sensorMsg.GetControllerOutput()
		if controllerData == nil {
			log.Error().Msg("Received message without controller data")
			continue
		}

		// Add the message to the queue of outstanding messages (and give up ownership)
		handlerQueue <- controllerData
	}
}
