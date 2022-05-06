package api

import "github.com/devdimensionlab/co-pilot/pkg/config"

var CallbackChannel = make(chan bool)

var CurrentProject config.Project
