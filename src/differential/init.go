package differential

import (
	"math"
)

func calculateMotorAngle(steerVal float32, throttle float32, trackWidth float32) (float32, float32) {
	dr := trackWidth // Distance between the centre of the rear tires
	if steerVal == 0 {
		return throttle, throttle
	}

	delta := math.Abs(float64(steerVal)) * 50

	innerRadius := 20 / float32(math.Tan(delta*math.Pi/180))
	outerRadius := innerRadius + dr

	adjustedThrottle := (innerRadius / outerRadius) * throttle

	if adjustedThrottle < 0 {
		adjustedThrottle = 0
	}

	if steerVal > 0 {
		return throttle, adjustedThrottle
	}

	return adjustedThrottle, throttle
}

func GetDiff(steeringAngle float32, leftThrottle float32, rightThrottle float32, trackWidth float32) (float32, float32) {

	if steeringAngle < 0.0 {
		// going left so make sure right wheel stays same
		return calculateMotorAngle(steeringAngle, rightThrottle, trackWidth)
	} else {
		// going right so make sure left wheel stays same
		return calculateMotorAngle(steeringAngle, leftThrottle, trackWidth)
	}
}
