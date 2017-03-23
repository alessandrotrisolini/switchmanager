package logging

import (
	"io"
	"log"
	"errors"
)

const TRACE_LEVEL	uint = 0x01
const DEBUG_LEVEL	uint = 0x02
const INFO_LEVEL	uint = 0x04
const WARN_LEVEL	uint = 0x08
const ERR_LEVEL		uint = 0x10

const TRACE		string	= "TRACE: "
const DEBUG		string	= "DEBUG: "
const INFO		string	= "INFO: "
const WARNING	string	= "WARNING: "
const ERROR		string	= "ERROR: "

const PREFIX	int		= log.Ldate|log.Ltime

type Log struct {
	logtrace	[]InnerLog
	logdebug	[]InnerLog
	loginfo		[]InnerLog
	logwarn		[]InnerLog
	logerr		[]InnerLog

	loglevel	uint
}

type InnerLog struct {
	log		*log.Logger
	h		io.Writer
}

func LogInit(handle io.Writer) (*Log) {
	var l Log

	l.loglevel = INFO_LEVEL | ERR_LEVEL
	
	l.logtrace 	= []InnerLog { {
		log : log.New(handle, TRACE, PREFIX),
		h 	: handle,
	}, }

	l.logdebug 	= []InnerLog { {
		log : log.New(handle, DEBUG, PREFIX),
		h	: handle,
	}, }

	l.loginfo 	= []InnerLog { {
		log : log.New(handle, INFO, PREFIX),
		h	: handle,
	}, }
	
	l.logwarn 	= []InnerLog { {
		log : log.New(handle, WARNING, PREFIX),
		h	: handle,
	}, }

	l.logerr 	= []InnerLog { {
		log : log.New(handle, ERROR, PREFIX),
		h	: handle,
	}, }
	
	return &l
}

/*
 *	Print functions
 */
func (l *Log) Trace(s string) {
	if l.loglevel & TRACE_LEVEL == 0x01 {
		for _, l := range l.logtrace {
			l.log.Println(s)		
		}
	}
}

func (l *Log) Debug(s string) {
	if l.loglevel & DEBUG_LEVEL == 0x01 {
		for _, l := range l.logdebug {
			l.log.Println(s)		
		}
	}
}

func (l *Log) Info(s string) {
	if l.loglevel & INFO_LEVEL == 0x01 {
		for _, l := range l.loginfo {
			l.log.Println(s)		
		}
	}
}

func (l *Log) Warn(s string) {
	if l.loglevel & WARN_LEVEL == 0x01 {
		for _, l := range l.logwarn {
			l.log.Println(s)		
		}
	}
}

func (l *Log) Error(s string) {
	if l.loglevel & ERR_LEVEL == 0x01 {
		for _, l := range l.logerr {
			l.log.Println(s)		
		}
	}
}


/*
 *	Management functions
 */
func (l *Log) AddTraceOutput(h io.Writer) {
	il := InnerLog {
			log : log.New(h, TRACE, PREFIX),
			h	: h, 
		}	
	l.logtrace = append(l.logtrace, il)
}


func (l *Log) AddDebugOutput(h io.Writer) {
	il := InnerLog {
			log : log.New(h, DEBUG, PREFIX),
			h	: h, 
		}	
	l.logdebug = append(l.logdebug, il)
}


func (l *Log) AddInfoOutput(h io.Writer) {
	il := InnerLog {
			log : log.New(h, INFO, PREFIX),
			h	: h, 
		}	
	l.loginfo = append(l.loginfo, il)
}


func (l *Log) AddWarnOutput(h io.Writer) {
	il := InnerLog {
			log : log.New(h, WARNING, PREFIX),
			h	: h, 
		}	
	l.logwarn = append(l.logwarn, il)
}


func (l *Log) AddErrorOutput(h io.Writer) {
	il := InnerLog {
			log : log.New(h, ERROR, PREFIX),
			h	: h, 
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
