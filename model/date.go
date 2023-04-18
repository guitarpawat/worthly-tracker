package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format(time.DateOnly)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func NewDate(in string) (Date, error) {
	datetime, err := time.Parse(time.DateOnly, in)
	if err != nil {
		return Date{}, err
	}

	date := Date(datetime)
	return date, nil
}

func MustNewDate(in string) Date {
	datetime, err := time.Parse(time.DateOnly, in)
	if err != nil {
		panic(err)
	}

	date := Date(datetime)
	return date
}

func FromNullDecimal(d decimal.NullDecimal) *decimal.Decimal {
	if !d.Valid {
		return nil
	}

	return &d.Decimal
}
