package progressio

import "fmt"
import "time"

type SecondFormatter int64

func (s SecondFormatter) Weeks() int64 {
	return int64(s) / (86400 * 7)
}

func (s SecondFormatter) Days() int64 {
	return (int64(s) / 86400) % 7
}

func (s SecondFormatter) Hours() int64 {
	return (int64(s) / 3600) % 24
}

func (s SecondFormatter) Minutes() int64 {
	return (int64(s) / 60) % 60
}
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

func FormatDuration(dur time.Duration) string {

	return SecondFormatter(dur.Seconds()).String()
}

func FormatSeconds(seconds int64) string {
	return SecondFormatter(seconds).String()
}

