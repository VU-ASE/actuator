# Overview

The `actuator` is the service that converts commands passed from another service into the actual movement of the rover. This service reads in actuator decisions (controller outputs) and turns it into hardware signals to steer motors and servo. It expects a `decision` stream as input from the any service and does not produce an output stream.



