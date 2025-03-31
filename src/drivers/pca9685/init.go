package pca9685

import (
	"fmt"

	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-pca9685"
	"github.com/rs/zerolog/log"
)

// ActuationIndex represents the indices for actuation values
type ActuationIndex int

const (
	Steer ActuationIndex = iota
	LeftThrottle
	RightThrottle
	Fan
	/* Add more actuation indices as needed */
)

const travelExtender float64 = 100.0

// PCA9685Controller represents the PCA9685 controller
type PCA9685Controller struct {
	i2cBus    *i2c.Options
	pca       *pca9685.PCA9685
	jumpTable map[ActuationIndex]channelConfig
}

type channelConfig struct {
	InUse        bool    /**< Is the channel in use? */
	MinPulseFrac float64 /**< Min PWM Pulse frac */
	MaxPulseFrac float64 /**< Max PWM Pulse frac */
	Trim         float64 /**< What is the center? 1 == no trim */
}

// NewPCA9685Controller creates a new PCA9685Controller instance
func NewPCA9685Controller(address uint8, bus uint) (*PCA9685Controller, error) {
	// Initialize I2C bus (normally at 3)
	i2cdevice := fmt.Sprintf("/dev/i2c-%d", bus)

	i2cBus, err := i2c.New(address, i2cdevice)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing I2C")
		return nil, err
	}

	pca, err := pca9685.New(i2cBus, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing PCA9685")
		i2cBus.Close()
		return nil, err
	}

	// Set update freq
	err = pca.SetFreq(50)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set PCA9685 freq")
	}

	// Create jump table with actuation indices
	jumpTable := map[ActuationIndex]channelConfig{
		ActuationIndex(Steer): {
			InUse:        true,
			MinPulseFrac: 205, // 1ms up
			MaxPulseFrac: 410, // 2ms up
			Trim:         0,
		},
		ActuationIndex(LeftThrottle): {
			InUse:        true,
			MinPulseFrac: 225,
			MaxPulseFrac: 385,
			Trim:         1,
		},
		ActuationIndex(RightThrottle): {
			InUse:        true,
			MinPulseFrac: 225,
			MaxPulseFrac: 385,
			Trim:         1,
		},
		ActuationIndex(Fan): {
			InUse:        true,
			MinPulseFrac: 225,
			MaxPulseFrac: 385,
			Trim:         1,
		},
		// Add more entries for other channels as needed
	}

	return &PCA9685Controller{
		i2cBus:    i2cBus,
		pca:       pca,
		jumpTable: jumpTable,
	}, nil
}

// Close closes the PCA9685Controller
func (pc *PCA9685Controller) Close() {
	err := pc.pca.Reset()
	if err != nil {
		log.Error().Err(err).Msg("Failed to reset PCA9685 board")
		return
	}

	err = pc.i2cBus.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close I2C bus")
		return
	}

}

// Update the frequency at index @p channel to @p value
func (pc *PCA9685Controller) SetTrim(channel ActuationIndex, trim float64) {

	// Get the value from the map
	config, exists := pc.jumpTable[channel]
	if !exists {
		log.Error().Int("channel", int(channel)).Msg("Error setting trim value")
		return
	}

	// Modify the value
	config.Trim = -1.0 * trim

	// Put the value back in the map
	pc.jumpTable[channel] = config
}

// Update the frequency at index @p channel to @p value for the motors
func (pc *PCA9685Controller) SetChannel(channel ActuationIndex, value float64) error {

	// Calculate duty cycles range
	minDuty := pc.jumpTable[channel].MinPulseFrac
	maxDuty := pc.jumpTable[channel].MaxPulseFrac
	dutyRange := maxDuty - minDuty

	dutyCycle := int(minDuty + (dutyRange * value))

	return pc.pca.SetChannel(int(channel), 0, dutyCycle)
}

func abs(value float64) float64 {
	if value < 0 {
		return value * -1.0
	}
	return value
}

// AllOff sets all channels to their midpoints)
func (pc *PCA9685Controller) AllOff() {
	log.Info().Msg("Turning off all channels")
	for i := ActuationIndex(0); i < 16; i++ {
		var err = pc.SetChannel(i, 0.5)
		if err != nil {
			log.Error().Err(err).Msg("Error turning all channels off")
		}
	}
}
