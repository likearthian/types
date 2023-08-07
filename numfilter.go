package utype

import (
	"fmt"
	"strconv"
	"strings"
)

type IntFilter struct {
	Valid bool
	comp  string
	value int64
}

func (i *IntFilter) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*i = IntFilter{
			Valid: false,
		}
		return nil
	}

	valid, comp, val, err := parseNumFilter(string(text))
	if err != nil {
		return err
	}

	*i = IntFilter{
		Valid: valid,
		comp:  comp,
		value: int64(val),
	}

	return nil
}

func (i *IntFilter) MarshalText() (text []byte, err error) {
	if !i.Valid {
		return []byte{}, nil
	}

	bld := strings.Builder{}
	bld.WriteString(strconv.FormatInt(i.value, 10))
	bld.WriteString("|")
	bld.WriteString(i.comp)

	return []byte(bld.String()), nil
}

type FloatFilter struct {
	Valid bool
	comp  string
	value float64
}

func (f *FloatFilter) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*f = FloatFilter{
			Valid: false,
		}
		return nil
	}

	valid, comp, val, err := parseNumFilter(string(text))
	if err != nil {
		return err
	}

	*f = FloatFilter{
		Valid: valid,
		comp:  comp,
		value: val,
	}

	return nil
}

func (f *FloatFilter) MarshalText() (text []byte, err error) {
	if !f.Valid {
		return []byte{}, nil
	}

	bld := strings.Builder{}
	bld.WriteString(strconv.FormatFloat(f.value, 'f', -1, 64))
	bld.WriteString("|")
	bld.WriteString(f.comp)

	return []byte(bld.String()), nil
}

func parseNumFilter(s string) (valid bool, comp string, val float64, err error) {
	strs := strings.Split(string(s), "|")
	valstr := strings.TrimSpace(strs[0])
	comp = "eq"
	if len(strs) > 1 {
		comp = strings.TrimSpace(strs[1])
	}

	switch comp {
	case "eq", "gt", "gte", "lt", "lte", "n", "nn":
		break
	default:
		return false, "", 0, fmt.Errorf("unknown comparison operator %q", comp)
	}

	val, err = strconv.ParseFloat(valstr, 64)
	if err != nil {
		return false, "", 0, err
	}

	valid = true
	return
}
