package fn

import (
	"scullion/ctx"
	"scullion/log"
	"time"
)

func NewDatetimeRegistrar() Registrar {
	return func(state *ctx.State) {
		dt := datetime{
			logger: state.LoggerWrapper(),
		}
		state.AddFunc("Add", dt.Add)
		state.AddFunc("After", dt.After)
		state.AddFunc("AfterDuration", dt.AfterDuration)
		state.AddFunc("AfterOrEqual", dt.AfterOrEqual)
		state.AddFunc("AfterOrEqualDuration", dt.AfterOrEqualDuration)
		state.AddFunc("Before", dt.Before)
		state.AddFunc("BeforeDuration", dt.BeforeDuration)
		state.AddFunc("BeforeOrEqual", dt.BeforeOrEqual)
		state.AddFunc("BeforeOrEqualDuration", dt.BeforeOrEqualDuration)
		state.AddFunc("Date", dt.Date)
		state.AddFunc("Duration", dt.Duration)
		state.AddFunc("Equal", dt.Equal)
		state.AddFunc("EqualDuration", dt.EqualDuration)
		state.AddFunc("Now", dt.Now)
		state.AddFunc("Since", dt.Since)
		state.AddFunc("Sub", dt.Sub)
	}
}

// See: https://github.com/antonmedv/expr/blob/master/docs/examples/dates_test.go
type datetime struct {
	logger log.Logger
}

func (dt datetime) Date(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		dt.logger.Errorf("date parse: %v", err)
	}
	return t
}
func (dt datetime) Duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		dt.logger.Errorf("duration parse: %v", err)
	}
	return d
}
func (datetime) Now() time.Time                                { return time.Now() }
func (datetime) Equal(a, b time.Time) bool                     { return a.Equal(b) }
func (datetime) Before(a, b time.Time) bool                    { return a.Before(b) }
func (datetime) BeforeOrEqual(a, b time.Time) bool             { return a.Before(b) || a.Equal(b) }
func (datetime) After(a, b time.Time) bool                     { return a.After(b) }
func (datetime) AfterOrEqual(a, b time.Time) bool              { return a.After(b) || a.Equal(b) }
func (datetime) Add(a time.Time, b time.Duration) time.Time    { return a.Add(b) }
func (datetime) Sub(a, b time.Time) time.Duration              { return a.Sub(b) }
func (datetime) EqualDuration(a, b time.Duration) bool         { return a == b }
func (datetime) BeforeDuration(a, b time.Duration) bool        { return a < b }
func (datetime) BeforeOrEqualDuration(a, b time.Duration) bool { return a <= b }
func (datetime) AfterDuration(a, b time.Duration) bool         { return a > b }
func (datetime) AfterOrEqualDuration(a, b time.Duration) bool  { return a >= b }
func (dt datetime) Since(s string) time.Duration               { return time.Since(dt.Date(s)) }
