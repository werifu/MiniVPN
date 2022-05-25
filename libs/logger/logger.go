package logger

import (
	"github.com/withmandala/go-log"
	"os"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout).WithColor()
}
