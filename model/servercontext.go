package model

import "go.uber.org/zap"

//ServerContext provides a single struct with expected input parameters
//for public functions
type ServerContext struct {
	Linker Linker
	Log    *zap.SugaredLogger
	ID     string
}
