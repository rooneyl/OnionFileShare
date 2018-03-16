package node

import (
	"log"
)

var Log *log.Logger

const (
	MinNumRoute   int = 3
	HeartBeatRate int = 600
	FileChunkSize int = 1000
)
