package log

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

const defaultTimestampFormat = time.RFC3339

type LogstashFormatter struct {
	Type string // if not empty use for logstash type field.

	// TimestampFormat sets the format used for timestamps.
	TimestampFormat string
}

func (f *LogstashFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := make(logrus.Fields)
	for k, v := range entry.Data {
		fields[k] = v
	}

	fields["@version"] = 1

	timeStampFormat := f.TimestampFormat

	if timeStampFormat == "" {
		timeStampFormat = defaultTimestampFormat
	}

	fields["@timestamp"] = entry.Time.Format(timeStampFormat)

	// set message field
	v, ok := entry.Data["message"]
	if ok {
		fields["fields.message"] = v
	}
	fields["message"] = entry.Message

	// set level field
	v, ok = entry.Data["level"]
	if ok {
		fields["fields.level"] = v
	}
	fields["level"] = entry.Level.String()

	// set type field
	if f.Type != "" {
		v, ok = entry.Data["type"]
		if ok {
			fields["fields.type"] = v
		}
		fields["type"] = f.Type
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
