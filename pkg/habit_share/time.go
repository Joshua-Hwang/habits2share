package habit_share

import (
	"time"
)

// Habit Sharing is done on the granularity of days. We can ignore timezones.
const DateFormat = "2006-01-02"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalText(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	t.Time, err = time.Parse(DateFormat, string(b))
	return
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	t.Time, err = time.Parse(`"`+DateFormat+`"`, string(b))
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}
