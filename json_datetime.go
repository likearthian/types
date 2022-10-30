package utype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gopkg.in/guregu/null.v4"
)

var JSONDateTimeFormat = "2006-01-02 15:04:05"

// JSONDateTime extending zero.Time type with custom JSON Marshalling
// into format 'YYYY-MM-DD HH:MM:SS'
type JSONDateTime struct {
	null.Time
}

func NewJSONDateTime(t time.Time) JSONDateTime {
	return JSONDateTime{null.TimeFrom(t)}
}

func (t *JSONDateTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}

	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		var date time.Time
		if x != "" {
			date, err = time.Parse(JSONDateTimeFormat, x)
			if err != nil {
				return fmt.Errorf("failed to parse time %s", x)
			}
		}

		*t = JSONDateTime{null.TimeFrom(date)}
	case nil:
		t.Valid = false
	case float64:
		u := time.UnixMilli(int64(x))
		*t = NewJSONDateTime(u)
		if u.Year() < 1980 && u.Year() > 2500 {
			t.Valid = false
		}
	default:
		return fmt.Errorf("json: cannot unmarshal %v into Go value of type JSONDateTime", reflect.TypeOf(v).Name())
	}

	return nil
}

func (t JSONDateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		//return (time.Time{}).MarshalJSON()
		return json.Marshal("")
	}
	//return t.Time.MarshalJSON()
	s := t.Time.Time
	return json.Marshal(s.Format(JSONDateTimeFormat))
}

func (t *JSONDateTime) UnmarshalText(text []byte) error {
	date, err := time.Parse(JSONDateTimeFormat, string(text))
	if err != nil {
		return fmt.Errorf("failed to parse time %s", string(text))
	}

	*t = NewJSONDateTime(date)
	return nil
}

func (t *JSONDateTime) MarshalText() (text []byte, err error) {
	str := t.ValueOrZero().Format(JSONDateTimeFormat)
	return []byte(str), nil
}

func (t JSONDateTime) TimeNow() time.Time {
	return time.Now()
}

func (t JSONDateTime) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

func (t JSONDateTime) ToString() string {
	return t.Time.Time.Format(JSONDateTimeFormat)
}

func (t JSONDateTime) Value() (driver.Value, error) {
	if t.Valid {
		return t.ToString(), nil
	} else {
		return nil, nil
	}
}
