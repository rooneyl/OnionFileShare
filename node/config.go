package node

import (
	"log"
	"time"
)

var Log *log.Logger

const (
	MinNumRoute       int           = 3 // OneWay
	MaxNumFileRequest int           = 3
	HeartBeatRate     time.Duration = 600
	FileChunkSize     int           = 100
	DEBUG             bool          = true
)
