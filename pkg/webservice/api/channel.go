package api

import "github.com/devdimensionlab/plybuild/pkg/config"

var CallbackChannel = make(chan bool)

var CurrentProject config.Project
