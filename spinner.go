package main

import (
	"github.com/theckman/yacspin"
	"time"
)

var cfg = yacspin.Config{
	Frequency:       100 * time.Millisecond,
	CharSet:         yacspin.CharSets[59],
	SuffixAutoColon: true,
	StopCharacter:   "✓",
	StopColors:      []string{"fgGreen"},
}
