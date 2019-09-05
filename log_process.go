package logprocess

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// LogProcess holds log information for logrus hook
type LogProcess struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	Level     logrus.Level
	PID       bool
	PName     bool
}

// InitLogProcess initialize LogProcess hook
func InitLogProcess(
	keepOutput, pid, processName bool, level logrus.Level,
	formatter logrus.Formatter, writer io.Writer,
) *LogProcess {
	if !keepOutput {
		logrus.SetOutput(ioutil.Discard)
	}

	return &LogProcess{
		Writer:    writer,
		Formatter: formatter,
		Level:     level,
		PID:       pid,
		PName:     processName,
	}
}

// Levels return the supported levels
func (h LogProcess) Levels() []logrus.Level {
	return logrus.AllLevels[:h.Level+1]
}

// Fire execute
func (h LogProcess) Fire(entry *logrus.Entry) error {
	if h.PID {
		entry.Data["pid"] = os.Getpid()
	}

	if h.PName {
		entry.Data["process_name"] = os.Args[0]
	}

	formatted, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.Writer.Write(formatted)
	if err != nil {
		return err
	}

	return nil
}
