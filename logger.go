package main

// some ideas came from https://stackoverflow.com/a/49004757/1843338

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogFormat struct {
	TimestampFormat string
}

var globalLogger *logrus.Entry

func init() {
	logger := logrus.New()
	logger.Level = logrus.TraceLevel
	formatter := LogFormat{}
	formatter.TimestampFormat = "2006-01-02T15:04:05.999999Z07:00"
	logger.SetFormatter(&formatter)
	logger.ReportCaller = false
	logger.Out = os.Stderr
	globalLogger = logger.WithFields(logrus.Fields{})
}

func (f *LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Time.Format(f.TimestampFormat))
	b.WriteString(" [")
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("] ")

	if entry.Message != "" {
		b.WriteString(entry.Message)
	}

	if len(entry.Data) > 0 {
		b.WriteString("\n‣ additional info begin\n")
		const payloadKey = `http.response.body`
		for key, value := range entry.Data {
			if key == payloadKey {
				fmt.Fprintf(b, "— %s=↓ below ↓\n%s\n—%s=↑ above ↑", payloadKey, value, payloadKey)
			} else {
				fmt.Fprintf(b, "— %s=", key)
				fmt.Fprint(b, value)

			}
			b.WriteByte('\n')
		}
		b.WriteString(`‣ additional info end`)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}
