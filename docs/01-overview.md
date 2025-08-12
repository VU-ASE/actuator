import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Overview

## Purpose 

The `actuator` service takes a "steering and acceleration decision" and turns it into hardware signals to spin the rear motors, steer the front wheels and to spin up a fan?

## Installation

To install this service, the latest release of [`roverctl`](https://ase.vu.nl/docs/framework/Software/rover/roverctl/installation) should be installed for your system and your Rover should be powered on.

<Tabs groupId="installation-method">
<TabItem value="roverctl" label="Using roverctl" default>

1. Install the service from your terminal
```bash
# Replace ROVER_NUMBER with your the number label on your Rover (e.g. 7)
roverctl service install -r <ROVER_NUMBER> https://github.com/VU-ASE/actuator/releases/latest/download/actuator.zip 
```

</TabItem>
<TabItem value="roverctl-web" label="Using roverctl-web">

1. Open `roverctl-web` for your Rover
```bash
# Replace ROVER_NUMBER with your the number label on your Rover (e.g. 7)
roverctl -r <ROVER_NUMBER>
```
2. Click on "install a service" button on the bottom left, and click "install from URL"
3. Enter the URL of the latest release:
```
https://github.com/VU-ASE/actuator/releases/latest/download/actuator.zip 
```

</TabItem>
</Tabs>

Follow [this tutorial](https://ase.vu.nl/docs/tutorials/write-a-service/upload) to understand how to use an ASE service. You can find more useful `roverctl` commands [here](/docs/framework/Software/rover/roverctl/usage)

## Requirements

- A [PCA-9865](https://ase.vu.nl/docs/framework/hardware/Components/carrier-board#wiring-connections) board should be connected to the Debix over I2C
    - If you want to configure the used I2C bus, you should modify this service's [*service.yaml*](https://github.com/VU-ASE/actuator/blob/main/service.yaml) and upload your service again.

## Inputs

As defined in the [*service.yaml*](https://github.com/VU-ASE/actuator/blob/main/service.yaml), this service depends on the following read streams:

- `decision`, exposed by a `controller` service:
    - From this stream, it will read a [`ControllerOutput`](https://github.com/VU-ASE/rovercom/blob/c1d6569558e26d323fecc17d01117dbd089609cc/definitions/outputs/controller.proto#L12) messages wrapped in a [`SensorOutput` wrapper message](https://github.com/VU-ASE/rovercom/blob/main/definitions/outputs/wrapper.proto)

## Outputs

As defined in the [*service.yaml*](https://github.com/VU-ASE/actuator/blob/main/service.yaml), this service does not expose any write streams.