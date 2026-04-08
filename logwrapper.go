// Package logwrapper contains functions to extend the zerolog logger with the function name (method) and filename + line number (at).
// After import it can be called in every function. First initialize FuncLogger and put the given zero logger in TraceFunc.
// Call TraceFunc with defer.
package logwrapper

// Example

// import (logwrapper)
// logger := logwrapper.FuncLogger(app.logger)
// defer logwrapper.TraceFunc(logger)()
// logger.Error().Msg(...)

// The log output will look like this:

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// LogFileNameAndLineNumber is a switch if file name and line number should be logged.
const LogFileNameAndLineNumber = false

func cutFuncName(input string) string {
	const (
		methodStringCount   = 3
		functionStringCount = 2
	)

	stringArray := strings.Split(input, ".")

	if len(stringArray) >= methodStringCount {
		method := stringArray[len(stringArray)-2]

		if strings.HasPrefix(method, "(") {
			return stringArray[len(stringArray)-1]
		}

		return stringArray[len(stringArray)-2] + ":" + stringArray[len(stringArray)-1]
	}

	if len(stringArray) == functionStringCount {
		return stringArray[1]
	}

	return "?" + input + "?"
}

// FuncLogger should be called before TraceFunc. It take a zerolog logger and returns a zerolog logger.
// It adds the function name to the logger and optional the filename and line number when LogFileNameAndLineNumber is true.
func FuncLogger(logger zerolog.Logger) zerolog.Logger {
	pc, file, line, ok := runtime.Caller(1)

	_logger := logger
	if !ok {
		return _logger
	}

	stringArray := strings.Split(file, "/")
	atLocation := stringArray[len(stringArray)-1] + ":" + strconv.Itoa(line)

	details := runtime.FuncForPC(pc)
	name := cutFuncName(details.Name())

	if details != nil {
		_logger = _logger.With().Str("method", name).Logger()
	}

	if LogFileNameAndLineNumber {
		_logger = _logger.With().Str("at", atLocation).Logger()
	}

	return _logger
}

// TraceFunc should be called after the FuncLogger function, it logs the entry and exit of the function at runtime.
// The TraceFunc must be called with defer.
func TraceFunc(logger zerolog.Logger) func() {
	logger.Trace().Msg("ENTRY")
	return func() { logger.Trace().Msg("EXIT") } // this will be deferred
}
