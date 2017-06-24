package utility

import (
	"os"
)

import log "gopkg.in/inconshreveable/log15.v2"

// LogManager ...
type LogManager struct {
	LogContext string
	ErrorFile  string
	InfoFile   string
	Logger     log.Logger
}

// InitLog ...
func (logM *LogManager) InitLog() {

	var svrlog = log.New(logM.LogContext, logM.LogContext)
	svrlog.SetHandler(log.MultiHandler(log.StreamHandler(os.Stderr, log.LogfmtFormat()),
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler(logM.ErrorFile, log.JsonFormat())),
		log.LvlFilterHandler(log.LvlInfo, log.Must.FileHandler(logM.InfoFile, log.JsonFormat()))))

	logM.Logger = svrlog
}
