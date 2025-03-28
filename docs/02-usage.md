# Usage

## Input

`actuator` expects a `decision` stream (see example below), which is usually sent by the `controller`.

```
actuatorOutput.Write(
    &pb_outputs.SensorOutput{
        SensorId:  2,
        Timestamp: uint64(time.Now().UnixMilli()),
        SensorOutput: &pb_outputs.SensorOutput_ControllerOutput{
            ControllerOutput: &pb_outputs.ControllerOutput{
                SteeringAngle: float32(steerValue),
                LeftThrottle:  float32(speed),
                RightThrottle: float32(speed),
                FrontLights:   false,
            },
        },
    },
)
```

`SensorOutput_ControllerOutput` contains a `ControllerOutput` object, composed of the following fields:

1. `SteeringAngle`, a `float32` value between -1 (left) and 1 (right)
2. `LeftThrottle`, a `float32` value between -1 (full reverse) and 1 (full forward)
3. `RightThrottle`, a `float32` value between -1 (full reverse) and 1 (full forward)
4. `FrontLights`, a `boolean` value previously used to turn on the front lights in the dark, currently not used

## Configuration values

1. **itwoc-bus** - the number of the I2C bus, reserved entirely for the motors, as well as the servo. [More details](https://ase.vu.nl/docs/framework/hardware/Components/carrier-board)
2. **electronic-diff** - takes a boolean value to decide whether or not to slow down the inner wheel while taking a turn. This essentially creates an illusion of a [differential](https://en.wikipedia.org/wiki/Differential_(mechanical_device)).
3. **track-width** - currently unused parameter, previously used for development and setting up the correct values in the **locking-diff**, should not be altered.
4. **servo-scaler** - increases or decreases the total movement range of the servo. Takes a value as a percentage, `1` is `100%` of a standard range. If you wish to increase the range to e.x. `150%`, you would set the value to `1.5`; to decrease to e.x. `90%`, choose value of `0.9`
5. **servo-trim** - used to calibrate the neutral position of the front wheels. Positive values will make the servo steer towards the right in the neutral position, negative - to the left. 
6. **fan-cap** - legacy value, currently not used.
