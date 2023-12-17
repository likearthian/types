package utype

import (
	"fmt"
	"strconv"
	"strings"
)

type NumFilter[T int64 | float64] struct {
	Valid bool
	Comp  string
	Value T
}

func (nf *NumFilter[T]) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*nf = NumFilter[T]{
			Valid: false,
		}
		return nil
	}

	valid, comp, val, err := parseGenericNumFilter[T](string(text))
	if err != nil {
		return err
	}

	*nf = NumFilter[T]{
		Valid: valid,
		Comp:  comp,
		Value: val,
	}

	return nil
}

func (nf *NumFilter[T]) MarshalText() (text []byte, err error) {
	if !nf.Valid {
		return []byte{}, nil
	}

	bld := strings.Builder{}
	bld.WriteString(fmt.Sprintf("%v", nf.Value))
	bld.WriteString("|")
	bld.WriteString(nf.Comp)

	return []byte(bld.String()), nil
}

type IntFilter struct {
	Valid bool
	Comp  string
	Value int64
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
		Comp:  comp,
		Value: int64(val),
	}

	return nil
}

func (i *IntFilter) MarshalText() (text []byte, err error) {
	if !i.Valid {
		return []byte{}, nil
	}

	bld := strings.Builder{}
	bld.WriteString(strconv.FormatInt(i.Value, 10))
	bld.WriteString("|")
	bld.WriteString(i.Comp)

	return []byte(bld.String()), nil
}

type FloatFilter struct {
	Valid bool
	Comp  string
	Value float64
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
		Comp:  comp,
		Value: val,
	}

	return nil
}

func (f *FloatFilter) MarshalText() (text []byte, err error) {
	if !f.Valid {
		return []byte{}, nil
	}

	bld := strings.Builder{}
	bld.WriteString(strconv.FormatFloat(f.Value, 'f', -1, 64))
	bld.WriteString("|")
	bld.WriteString(f.Comp)

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

func parseGenericNumFilter[T int64 | float64](s string) (valid bool, comp string, val T, err error) {
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

	fval, err := strconv.ParseFloat(valstr, 64)
	if err != nil {
		return false, "", 0, err
	}

	val = T(fval)

	valid = true
	return
}
