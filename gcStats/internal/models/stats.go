package models

import "time"

type Stats struct {
	Malloc    uint64
	NumGC     uint32
	HeapAlloc uint64
	LastGC    uint64
	Time      time.Time
}

var AllStats []Stats
