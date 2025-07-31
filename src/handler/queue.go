package handler

import (
	pb_module_outputs "github.com/VU-ASE/rovercom/v2/packages/go/outputs"
)

type Queue = chan *pb_module_outputs.ControllerOutput
