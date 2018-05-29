package model

import "go.uber.org/zap"

//ServerContext provides a single struct with expected input parameters
//for public functions
type ServerContext struct {
	linker Linker
	log    *zap.SugaredLogger
}
