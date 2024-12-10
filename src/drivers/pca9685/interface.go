package pca9685

import (
	"github.com/rs/zerolog/log"
)

const servoLimiter = 0.8

// Set the servo duty. In range -1 (left) to 1(right).
func (pc *PCA9685Controller) SetServo(value float64, servoScaler float64) {
	value = value * -servoLimiter
	value = clamp(value)

	trim := pc.jumpTable[Steer].Trim
	if trim != 0 {
		value = value * abs(1.0-(trim*0.75)) * servoScaler
	} else {
		value = value * servoScaler
	}

	value = ((value + 1) / 2)

	err := pc.SetChannelWithTrim(Steer, value)
	if err != nil {
		log.Error().Err(err).Msg("Error setting Steer value")
	}
}

func (pc *PCA9685Controller) SetServoTrim(value float64) {
	pc.SetTrim(Steer, value)
}

// Set the left motor power. In range -1 to 1.
func (pc *PCA9685Controller) SetLeftMotor(value float64) {
	value = clamp(value)
	value = (value + 1) / 2
	err := pc.SetChannel(LeftThrottle, value)
	if err != nil {
		log.Error().Err(err).Msg("Error setting LeftThrottle value")
	}
}

// Set right motor power. In range -1 to 1.
func (pc *PCA9685Controller) SetRightMotor(value float64) {
	value = clamp(value)
	value = (value + 1) / 2
	err := pc.SetChannel(RightThrottle, value)
	if err != nil {
		log.Error().Err(err).Msg("Error setting RightMotor value")
	}
}

// Set fan power. In range 0 to 1
func (pc *PCA9685Controller) SetFan(value float64) {
	value = clamp(value)
	err := pc.SetChannel(Fan, value)
	if err != nil {
		log.Error().Err(err).Msg("Error setting Fan value")
	}
}

// Private function to clamp a value between -1 and 1
func clamp(value float64) float64 {
	if value > 1 {
		return 1
	}
	if value < -1 {
		return -1
	}
	return value
}
