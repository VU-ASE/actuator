package handler

import (
	"time"
	diff "vu/ase/actuator/src/differential" // Import the package that contains the GetDiff function
	"vu/ase/actuator/src/drivers/pca9685"

	"github.com/rs/zerolog/log"
)

const MSG_TIMEOUT = 500 * time.Millisecond

var pca *pca9685.PCA9685Controller
var handlerQueue Queue

// OTA tunable parameters
var servoScaler float64
var servoTrim float64

func Start(queue Queue, i2cbus uint, serScaler float64, serTrim float64, enableDiff bool, trackWidth float64, fanCap int) {
	servoScaler = serScaler
	servoTrim = serTrim

	// Main handler loop
	newPca, err := pca9685.NewPCA9685Controller(0x40, i2cbus)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Initialize pca driver!")
		return
	}
	handlerQueue = queue
	pca = newPca
	defer pca.Close()

	// Apply the servo trim
	pca.SetServoTrim(servoTrim)

	// Keep track of when the last message was received
	// if this was more than .5 seconds ago, we should stop the motors
	lastMessageTime := time.Now()
	go func() {
		for {
			if time.Since(lastMessageTime) > MSG_TIMEOUT {
				pca.AllOff()
				pca.SetFan(0)
				// log.Warn().Msg("No message received for 500ms, stopping motors")
			}
			time.Sleep(MSG_TIMEOUT)
		}
	}()

	for {
		// Receive the pointer to the next message
		msg := <-handlerQueue
		if msg == nil {
			continue
		}
		lastMessageTime = time.Now()

		// Apply the differential
		if enableDiff {
			l, r := diff.GetDiff(
				float64(msg.SteeringAngle),
				float64(msg.LeftThrottle),
				float64(msg.RightThrottle),
				trackWidth,
			) // Call the GetDiff function from the imported package
			msg.LeftThrottle = float32(l)
			msg.RightThrottle = float32(r)
		}
		// log.Debug().float64("steering angle", msg.SteeringAngle).float64("left motor", msg.LeftThrottle).float64("right motor", msg.RightThrottle).Msg("New message available")

		// Process the message (let the drivers handle this)
		pca.SetServo(float64(msg.SteeringAngle), servoScaler)
		pca.SetLeftMotor(float64(msg.LeftThrottle))
		pca.SetRightMotor(float64(msg.RightThrottle))

		// Apply fan cap to fan speed
		fanSpeed := float64(msg.FanSpeed) * (float64(fanCap) / 100.0)
		pca.SetFan(fanSpeed)
	}
}

func SetServoTrim(trim float64) {
	servoTrim = trim
}

func SetServoScaler(scaler float64) {
	servoScaler = scaler
}

func OnTerminate() {
	log.Info().Msg("Requested handler to terminate")
	if pca == nil {
		log.Warn().Msg("Handler not running, nothing to terminate")
		return
	}

	handlerQueue = nil // Stop receiving messages, and wait for the last message to be processed
	// time.Sleep(50 * time.Millisecond)
	pca.AllOff()
	pca.SetFan(0)
	pca.Close()
	// time.Sleep(400 * time.Millisecond) // Wait for I2C to finish (I guess buffer flush?)
	log.Info().Msg("Handler terminated")
}
