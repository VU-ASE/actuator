# Processing

This service follows the following steps:

1. It reads a [`ControllerOutput` message](https://github.com/VU-ASE/rovercom/blob/c1d6569558e26d323fecc17d01117dbd089609cc/definitions/outputs/controller.proto#L12), from the `decision` stream as defined in the [*service.yaml* file](https://github.com/VU-ASE/actuator/blob/7492c12e3f4187609f25779549434fa4d05a8115/service.yaml#L14)
2. It takes the [acceleration/throttle](https://github.com/VU-ASE/rovercom/blob/c1d6569558e26d323fecc17d01117dbd089609cc/definitions/outputs/controller.proto#L16) values and the [steering angle](https://github.com/VU-ASE/rovercom/blob/c1d6569558e26d323fecc17d01117dbd089609cc/definitions/outputs/controller.proto#L14C5-L14C29) value and creates an I2C command for the [PCA-9865](https://ase.vu.nl/docs/framework/hardware/Components/carrier-board#wiring-connections) board. This board will then transform the I2C value into a PWM signal that the motors and servo use