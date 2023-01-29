package stellar

import "log"

type Logger interface {
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

var logger Logger = &SimpleLogger{}

type SimpleLogger struct{}

func (l *SimpleLogger) Debugf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (l *SimpleLogger) Infof(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (l *SimpleLogger) Warnf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (l *SimpleLogger) Errorf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (l *SimpleLogger) Fatalf(fmt string, args ...interface{}) {
	log.Fatalf(fmt, args...)
}

const (
	green  = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	blue   = "\033[97;44m"
	reset  = "\033[0m"
)

type Level int

const (
	Fatal Level = iota
	Error
	Warn
	Info
	Debug
)

var levelState = map[Level]bool{
	Fatal: true,
	Error: true,
	Warn:  true,
	Info:  true,
	Debug: false,
}

// StdLogger for log print on std console.
type StdLogger struct{}

// NewStdLogger ...
func UseStdLogger() {
	logger = &StdLogger{}
}

// SetLogLevel ...
func SetLogLevel(level Level, on bool) {
	levelState[level] = on
}

func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	if levelState[Fatal] {
		log.Fatalf(red+"[Fatal] "+format+reset, v...)
	}
}

func (l *StdLogger) Errorf(format string, v ...interface{}) {
	if levelState[Error] {
		log.Printf(red+"[Error] "+format+reset, v...)
	}

}

func (l *StdLogger) Warnf(format string, v ...interface{}) {
	if levelState[Warn] {
		log.Printf(yellow+"[Warn ] "+format+reset, v...)
	}
}

func (l *StdLogger) Infof(format string, v ...interface{}) {
	if levelState[Info] {
		log.Printf(green+"[Info ] "+format+reset, v...)
	}
}

func (l *StdLogger) Debugf(format string, v ...interface{}) {
	if levelState[Debug] {
		log.Printf(blue+"[Debug] "+format+reset, v...)
	}
}
