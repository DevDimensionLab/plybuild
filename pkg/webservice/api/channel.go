package api

import "github.com/devdimensionlab/ply/pkg/config"

var CallbackChannel = make(chan bool)

var CurrentProject config.Project
