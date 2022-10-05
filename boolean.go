package utype

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type Boolean struct {
	null.Bool
}

func BooleanFrom(b bool) Boolean {
	return Boolean{null.BoolFrom(b)}
}

func (b *Boolean) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		b.Valid = false
		return nil
	}

	var str string
	var num int
	var raw null.Bool

	strErr := json.Unmarshal(data, &str)
	numErr := json.Unmarshal(data, &num)
	rawErr := json.Unmarshal(data, &raw)

	if strErr != nil && rawErr != nil && numErr != nil {
		fmt.Println("nothing works")
		return rawErr
	}

	if rawErr == nil {
		b.Bool = raw
		return nil
	}

	if numErr == nil {
		str = strconv.Itoa(num)
	}

	bval, err := convertToBool(str)
	if err != nil {
		return err
	}

	*b = BooleanFrom(bval)

	return nil
}

func (b *Boolean) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*b = BooleanFrom(true)
		return nil
	}

	bval, err := convertToBool(string(text))
	if err != nil {
		return err
	}

	*b = BooleanFrom(bval)

	return nil
}

func (b *Boolean) MarshalText() (text []byte, err error) {
	if !b.Valid {
		return []byte{}, nil
	}

	if b.ValueOrZero() {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}

func (b Boolean) Value() (driver.Value, error) {
	if b.Valid {
		if b.ValueOrZero() {
			return 1, nil
		}
		return 0, nil
	} else {
		return nil, nil
	}
}

func (b *Boolean) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		b.Valid = false
		b.Bool.Bool = false
		return nil
	}

	bs, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}

	var str string
	switch bs := bs.(type) {
	case string:
		str = bs
	case []byte:
		str = string(bs)
	case int:
		str = strconv.Itoa(bs)
	}

	bval, err := convertToBool(str)
	if err != nil {
		return err
	}

	*b = BooleanFrom(bval)

	return nil
}

func convertToBool(str string) (bool, error) {
	var b bool
	if strings.EqualFold(str, "FALSE") || str == "0" {
		b = false
	} else if strings.EqualFold(str, "TRUE") || str == "1" {
		b = true
	} else {
		return false, fmt.Errorf("cannot unmarshal the value %q into boolean type", str)
	}

	return b, nil
}
