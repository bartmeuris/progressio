package progressio

import "fmt"
import "time"

// SecondFormatter represents a duration in seconds
type SecondFormatter int64

// Weeks returns the amount of whole weeks represented by the
// SecondFormatter instance
// Minimum value: 0
func (s SecondFormatter) Weeks() int64 {
	return int64(s) / (86400 * 7)
}

// Days returns the amount of days of a partial week represented
// by the SecondFormatter instance
// Minimum value: 0, Maximum Value: 6
func (s SecondFormatter) Days() int64 {
	return (int64(s) / 86400) % 7
}

// Hours returns the amount of hours of the last partial day
// represented by the SecondFormatter instance
// Minumum value: 0, Maximum value: 23
func (s SecondFormatter) Hours() int64 {
	return (int64(s) / 3600) % 24
}

// Minutes returns the amount of minutes of the last partial hour
// represented by the SecondFormatter instance
// Minumum value: 0, Maximum value: 59
func (s SecondFormatter) Minutes() int64 {
	return (int64(s) / 60) % 60
}

// Seconds returns the amount of seconds of the last partial minute
// represented by the SecondFormatter instance
// Minumum value: 0, Maximum value: 59
func (s SecondFormatter) Seconds() int64 {
	return int64(s) % 60
}

func addCountString(ostr string, val int64, str string) string {
	if val == 0 {
		return ostr
	}
	if len(ostr) != 0 {
		ostr += ", "
	}
	if val < 0 {
		val = -val
	}
	if val == 1 {
		return ostr + "1 " + str
	}
	return fmt.Sprintf("%s%d %ss", ostr, val, str)
}

// String returns the string representation of the SecondFormatter
// instance, specifying (if applicable): the amount of weeks, days,
// hours, minutes and seconds it represents
func (s SecondFormatter) String() string {
	sret := ""
	sret = addCountString(sret, s.Weeks(), "week")
	sret = addCountString(sret, s.Days(), "day")
	sret = addCountString(sret, s.Hours(), "hour")
	sret = addCountString(sret, s.Minutes(), "minute")
	sret = addCountString(sret, s.Seconds(), "second")
	if len(sret) == 0 {
		sret = "0 seconds"
	} else if s < 0 {
		sret += " ago"
	}
	return sret
}

// FormatDuration returns the string representation of the specified
// time.Duration, in the (if applicable) amount of weeks, days, hours,
// minutes and seconds it represents
func FormatDuration(dur time.Duration) string {
	return SecondFormatter(dur.Seconds()).String()
}

// FormatSeconds returns a string representing the (if applicable)
// amount of weeks, days, hours, minutes and seconds the amount of
// seconds it is passed as a parameter.
func FormatSeconds(seconds int64) string {
	return SecondFormatter(seconds).String()
}
