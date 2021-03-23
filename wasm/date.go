package wasm

import (
	"syscall/js"
	"time"
)

// Date is an instance of a JS Date.
// The zero value of this struct is not a valid Date.
type Date struct {
	Object
}

// ToTime converts the date to a Go Time type.
// It uses the JavaScript `getTime` method to get Unix time elapsed in milliseconds.
// ToTime will not have accuracy to nanoseconds, only milliseconds.
// If it is unable to call `getTime` the function will return an error.
func (d Date) ToTime() (time.Time, error) {
	var millis int64

	getTime, err := d.Get("getTime")
	if err != nil {
		return time.Time{}, err
	}

	// `This` context matters for calling JS.
	// Calling `Call` directly on Date means potential panics may occur.
	// https://stackoverflow.com/questions/17899598
	value := getTime.Call("call", d)
	if err := FromJSValue(value, &millis); err != nil {
		return time.Time{}, err
	}
	seconds := millis / 1e3
	offsetMillis := millis - (seconds * 1e3)
	nano := offsetMillis * 1e6

	return time.Unix(seconds, nano), nil
}

// FromJSValue turns a JS value to a Date.
func (d *Date) FromJSValue(value js.Value) error {
	var err error
	d.Object, err = NewObject(value)
	return err
}

// NewDate returns a JS Date from the provied time.Time.
// Converts from Go value to JS through Unix time elapsed.
func NewDate(t time.Time) Date {
	// Using `t.UnixNano()` will overflow for values before 1678 or after 2262.
	// Because JS only has millisecond precision, we can manually use the Nanosecond offset to retain dates farther away
	// from Unix time.
	millis := t.Unix() * 1e3
	millis += int64(t.Nanosecond() % 1e9 / 1e6)

	date, err := Global().Expect(js.TypeFunction, "Date")
	if err != nil {
		panic("Date constructor not found")
	}

	return mustJSValueToDate(date.New(millis))
}

func mustJSValueToDate(v js.Value) Date {
	var d Date
	err := d.FromJSValue(v)
	if err != nil {
		panic("Expected a Date from JS Standard library")
	}

	return d
}
