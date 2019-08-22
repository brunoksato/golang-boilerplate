package core

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	"github.com/brunoksato/golang-boilerplate/util"
)

type Timestamp struct {
	Time time.Time
}

// Value implements the driver Valuer interface.
func (ts Timestamp) Value() (driver.Value, error) {
	return ts.Time, nil
}

func (ts *Timestamp) Scan(value interface{}) error {
	ts.Time = value.(time.Time)
	return nil
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(ts.Timestamp()))
}

func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	tsInt, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	ts.SetTimestamp(int64(tsInt))
	return nil
}

func (ts Timestamp) Timestamp() int64 {
	return util.NanoToMicro(ts.Time.UnixNano())
}

func (ts *Timestamp) SetTimestamp(val int64) {
	ts.Time = time.Unix(0, util.MicroToNano(val))
}

func NewTimestamp(val int64) Timestamp {
	t := time.Unix(0, util.MicroToNano(val))
	return Timestamp{Time: t}
}

type NullableTimestamp struct {
	Time time.Time
}

// Value implements the driver Valuer interface.
func (ts *NullableTimestamp) Value() (driver.Value, error) {
	if ts == nil {
		return nil, nil
	}

	return ts.Time, nil
}

func (ts *NullableTimestamp) Scan(value interface{}) error {
	if ts == nil {
		return nil
	} else {
		if value != nil {
			t := value.(time.Time)
			ts.Time = t
		}
	}
	return nil
}

func (ts *NullableTimestamp) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}

	return json.Marshal(int64(ts.Timestamp()))
}

func (ts *NullableTimestamp) UnmarshalJSON(data []byte) error {
	if ts == nil {
		return nil
	} else {
		if len(data) > 0 {
			tsInt, err := strconv.Atoi(string(data))
			if err != nil {
				return err
			}
			ts.SetTimestamp(int64(tsInt))
		}
	}
	return nil
}

func (ts NullableTimestamp) Timestamp() int64 {
	return util.NanoToMicro(ts.Time.UnixNano())
}

func (ts *NullableTimestamp) SetTimestamp(val int64) {
	ts.Time = time.Unix(0, util.MicroToNano(val))
}

func NewNullableTimestamp(val int64) *NullableTimestamp {
	t := time.Unix(0, util.MicroToNano(val))
	return &NullableTimestamp{Time: t}
}
