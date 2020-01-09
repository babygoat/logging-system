package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/devtools/clouderrorreporting/v1beta1"
	"google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/genproto/googleapis/logging/v2"
)

type (
	Stackdriver struct {
		logging.LogEntry
		clouderrorreporting.ErrorEvent
		Payload map[string]interface{}
	}

	causer interface {
		Cause() error
	}

	serviceContext struct {
		Service string `json:"service"`
		Version string `json:"version"`
	}

	stackTracer interface {
		StackTrace() errors.StackTrace
	}
)

func GetGoroutineState() string {
	stack := make([]byte, 64)
	stack = stack[:runtime.Stack(stack, false)]
	stack = stack[:bytes.Index(stack, []byte("\n"))]

	return string(stack)
}

// Adapt from https://github.com/googleapis/google-cloud-go/issues/1084#issuecomment-474565019
func FormatStack(err error) (buffer []byte) {
	if err == nil {
		return nil
	}

	// find the inner most error with a stack
	inner := err
	for inner != nil {
		if cause, ok := inner.(causer); ok {
			inner = cause.Cause()
			if _, ok := inner.(stackTracer); ok {
				err = inner
			}
		} else {
			break
		}
	}

	if stackTrace, ok := err.(stackTracer); ok {
		buf := bytes.Buffer{}
		buf.WriteString(GetGoroutineState() + "\n")

		// format each frame of the stack to match runtime.Stack's format
		var lines []string
		for _, frame := range stackTrace.StackTrace() {
			pc := uintptr(frame) - 1
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				file, line := fn.FileLine(pc)
				lines = append(lines, fmt.Sprintf("%s()\n\t%s:%d +%#x", fn.Name(), file, line, fn.Entry()))
			}
		}
		buf.WriteString(strings.Join(lines, "\n"))

		buffer = buf.Bytes()
	}
	return
}

func NewStackdriverFormatter(service, version string) *Stackdriver {
	return &Stackdriver{
		ErrorEvent: clouderrorreporting.ErrorEvent{
			ServiceContext: &clouderrorreporting.ServiceContext{
				Service: service,
				Version: version,
			},
		},
	}
}

func (s *Stackdriver) Format(entry *logrus.Entry) ([]byte, error) {
	// Copy customized fields
	s.Payload = make(logrus.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			s.Payload[k] = v.Error()
		default:
			s.Payload[k] = v
		}
	}

	s.Message = entry.Message
	s.Severity = convertLevelToLogSeverity(entry.Level)

	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = new(bytes.Buffer)
	}

	encoder := json.NewEncoder(b)

	if err := encoder.Encode(s); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %+v", err)
	}

	return b.Bytes(), nil
}

func convertLevelToLogSeverity(lvl logrus.Level) ltype.LogSeverity {
	if lvl == logrus.InfoLevel {
		return ltype.LogSeverity_INFO
	}
	return ltype.LogSeverity_ERROR
}
