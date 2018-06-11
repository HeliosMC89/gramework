// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

var defaultMaxHackAttempts int32 = 5

// New App
func New() *App {
	logger := &log.Logger{
		Level:   log.DebugLevel,
		Handler: cli.New(os.Stdout),
	}
	flags := &Flags{
		values: make(map[string]Flag),
	}
	defFWLimit := int64(-1)
	defBlockTimeout := int64(-1)
	maxHackAttempts := defaultMaxHackAttempts
	app := &App{
		Flags:          flags,
		flagsQueue:     flagsToRegister,
		Logger:         logger,
		domainListLock: new(sync.RWMutex),
		firewall: &firewall{
			blockList:      make(map[string]int64),
			MaxReqPerMin:   &defFWLimit,
			BlockTimeout:   &defBlockTimeout,
			requestCounter: make(map[string]int64),
		},
		firewallInit:              new(sync.Once),
		domains:                   make(map[string]*Router),
		middlewaresMu:             new(sync.RWMutex),
		middlewaresAfterRequestMu: new(sync.RWMutex),
		preMiddlewaresMu:          new(sync.RWMutex),
		middlewares:               make([]func(*Context), 0),
		middlewaresAfterRequest:   make([]func(*Context), 0),
		preMiddlewares:            make([]func(*Context), 0),
		seed:                      uintptr(time.Now().Nanosecond()),
		maxHackAttempts:           &maxHackAttempts,
	}

	app.defaultRouter = &Router{
		router: newRouter(),
		app:    app,
	}

	return app
}
