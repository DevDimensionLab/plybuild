package api

import "github.com/co-pilot-cli/co-pilot/pkg/config"

var CallbackChannel = make(chan bool)

var CurrentProject config.Project
