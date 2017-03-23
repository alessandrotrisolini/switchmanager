package logging

import (
	"errors"
	"io"
	"log"
)

const TRACE_LEVEL uint = 0x01
const DEBUG_LEVEL uint = 0x02
const INFO_LEVEL uint = 0x04
const WARN_LEVEL uint = 0x08
const ERR_LEVEL uint = 0x10

const TRACE string = "TRACE: "
const DEBUG string = "DEBUG: "
const INFO string = "INFO: "
const WARNING string = "WARNING: "
const ERROR string = "ERROR: "

const PREFIX int = log.Ldate | log.Ltime

type Log struct {
	logtrace []innerLog
	logdebug []innerLog
	loginfo  []innerLog
	logwarn  []innerLog
	logerr   []innerLog

	loglevel uint
	initialized bool
}

type innerLog struct {
	log *log.Logger
	h   io.Writer
}

var l Log = Log{initialized: false}

func LogInit(handle io.Writer) error {
	if !l.initialized {
		l.initialized = true
		l.loglevel = INFO_LEVEL | ERR_LEVEL

		l.logtrace = []innerLog{{
			log: log.New(handle, TRACE, PREFIX),
			h:   handle,
		}}

		l.logdebug = []innerLog{{
			log: log.New(handle, DEBUG, PREFIX),
			h:   handle,
		}}

		l.loginfo = []innerLog{{
			log: log.New(handle, INFO, PREFIX),
			h:   handle,
		}}

		l.logwarn = []innerLog{{
			log: log.New(handle, WARNING, PREFIX),
			h:   handle,
		}}

		l.logerr = []innerLog{{
			log: log.New(handle, ERROR, PREFIX),
			h:   handle,
		}}
	} else {
		return errors.New("A log already exists: get it by calling logging.GetLogger()")
	}

	return nil
}

func GetLogger() *Log { return &l }

/*
 *	Print functions
 */
func (l *Log) Trace(a ...interface{}) {
	if l.loglevel&TRACE_LEVEL == 0x01 {
		for _, l := range l.logtrace {
			l.log.Println(a)
		}
	}
}

func (l *Log) Debug(a ...interface{}) {
	if l.loglevel&DEBUG_LEVEL == 0x02 {
		for _, l := range l.logdebug {
			l.log.Println(a)
		}
	}
}

func (l *Log) Info(a ...interface{}) {
	if l.loglevel&INFO_LEVEL == 0x04 {
		for _, l := range l.loginfo {
			l.log.Println(a)
		}
	}
}

func (l *Log) Warn(a ...interface{}) {
	if l.loglevel&WARN_LEVEL == 0x08 {
		for _, l := range l.logwarn {
			l.log.Println(a)
		}
	}
}

func (l *Log) Error(a ...interface{}) {
	if l.loglevel&ERR_LEVEL == 0x10 {
		for _, l := range l.logerr {
			l.log.Println(a)
		}
	}
}

/*
 *	Management functions
 */
func (l *Log) AddTraceOutput(h io.Writer) {
	il := innerLog{
		log: log.New(h, TRACE, PREFIX),
		h:   h,
	}
	l.logtrace = append(l.logtrace, il)
}

func (l *Log) AddDebugOutput(h io.Writer) {
	il := innerLog{
		log: log.New(h, DEBUG, PREFIX),
		h:   h,
	}
	l.logdebug = append(l.logdebug, il)
}

func (l *Log) AddInfoOutput(h io.Writer) {
	il := innerLog{
		log: log.New(h, INFO, PREFIX),
		h:   h,
	}
	l.loginfo = append(l.loginfo, il)
}

func (l *Log) AddWarnOutput(h io.Writer) {
	il := innerLog{
		log: log.New(h, WARNING, PREFIX),
		h:   h,
	}
	l.logwarn = append(l.logwarn, il)
}

func (l *Log) AddErrorOutput(h io.Writer) {
	il := innerLog{
		log: log.New(h, ERROR, PREFIX),
		h:   h,
	}
	l.logerr = append(l.logerr, il)
}

func (l *Log) SetDebugLevel(level uint) error {
	if level > 0x1F {
		return errors.New("Invalid log level")
	}
	l.loglevel = level
	return nil
}
