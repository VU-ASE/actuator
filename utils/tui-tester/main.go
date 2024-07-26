// Package main provides a utility for publishing actuation data and running a TUI to modify the data.
package main

//
// todo: this psackage should be published as a separate package, that uses the roverlib to connecxt to the Actuator service
//

import (
	"strconv"
	// "time"

	pb_module_outputs "github.com/VU-ASE/rovercom/packages/go/outputs"
	"google.golang.org/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

type Queue = chan *pb_module_outputs.ControllerOutput

// Publisher runs the publisher routine to simulate continuous publishing of actuation data.
// The actuation data is encoded to protobuf and sent over a ZeroMQ PUB socket.
func Publisher(msg Queue) {
	// Create a ZeroMQ PUB socket
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close() // Close the socket when the Goroutine exits
	err := publisher.Bind("tcp://*:5555")
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind to port 5555")
		return
	}

	for {
		actuationData := <-msg

		// Encode actuation data to protobuf
		actuationDataBytes, _ := proto.Marshal(actuationData)

		// Send actuation data as a single multipart message
		_, err = publisher.SendBytes(actuationDataBytes, 0)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message!")
			return
		}
	}
}

// RunTUI runs a fancy TUI application that allows the user to modify actuation data.
// It uses tview package for the TUI components, including input fields and buttons.
func RunTUI(msg Queue) {
	// Create a new TUI application
	app := tview.NewApplication()

	controlData := pb_module_outputs.ControllerOutput{}

	// Create a form with input fields and buttons
	form := tview.NewForm().
		AddInputField("LeftThrottle", "", 20, nil, func(text string) {
			l_throttle, err := strconv.ParseFloat(text, 32)
			if err == nil {
				controlData.LeftThrottle = float32(l_throttle)
			}
		}).
		AddInputField("RightThrottle", "", 20, nil, func(text string) {
			r_throttle, err := strconv.ParseFloat(text, 32)
			if err == nil {
				controlData.RightThrottle = float32(r_throttle)
			}
		}).
		AddInputField("Steer", "", 20, nil, func(text string) {
			steer, err := strconv.ParseFloat(text, 32)
			if err == nil {
				controlData.SteeringAngle = float32(steer)
			}
		}).
		AddButton("Send", func() {
			msg <- &controlData
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	// Set the form as the root of the TUI application
	if err := app.SetRoot(form, true).Run(); err != nil {
		panic(err)
	}
}

// main is the entry point of the program.
// It initializes the actuation data, runs the publisher and TUI routines concurrently,
// and prevents the main goroutine from exiting.
func main() {
	// Initialize actuation data
	msgQueue := make(chan *pb_module_outputs.ControllerOutput)

	// Run the publisher and TUI routines concurrently
	go Publisher(msgQueue)
	go RunTUI(msgQueue)

	// Prevent the main goroutine from exiting
	select {}
}
