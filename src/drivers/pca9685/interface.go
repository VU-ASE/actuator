package pca9685

import (
	"github.com/rs/zerolog/log"
)


// Set the servo duty. In range -1 (left) to 1(right).
func (pc *PCA9685Controller) SetServo(value float64, servoScaler float64) (err error)  {
	value = clamp(value)

	// Calculate the duty cycle range
	maxDuty := pc.jumpTable[Steer].MaxPulseFrac
	minDuty := pc.jumpTable[Steer].MinPulseFrac
	dutyRange := maxDuty - minDuty

	// This gets added or subtracted from the midpoint
	halfRange := dutyRange / 2.0

	// We do not want to "oversteer" the servo, which might damage it
	// so we limit the range
	halfRange *= servoScaler

	// Find the center of the range
	center := minDuty + (dutyRange / 2.0)
	servoTrim := pc.jumpTable[Steer].Trim

	// Calculate the new duty cycle
	// (due to an annoying convention, we do - right and + left)
	dutyCycle := int(center - (halfRange * (value - servoTrim)))
	return pc.pca.SetChannel(int(Steer), 0, dutyCycle)
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
