package crond

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type CrondFunc func()

const (
	// 秒(0~59)
	TimerIndex_Second = 0

	// 分(0~59)
	TimerIndex_Minute = 1

	// 时(0~23)
	TimerIndex_Hour = 2

	// 日(1~31)
	TimerIndex_Day = 3

	// 月份(1~12)
	TimerIndex_Month = 4

	// 星期(1~7)
	TimerIndex_Weekday = 5

	// tmier长度
	TimerLength = 6
)

type CrondJob struct {
	// 定时名称
	Name string

	// 触发时间
	Timer []int

	// 执行方法
	Func CrondFunc
}

func GetBaseExecTime(tn time.Time, loc *time.Location, timer []int) time.Time {
	var month = tn.Month()

	if timer[TimerIndex_Month] != -1 {
		month = time.Month(timer[TimerIndex_Month])
	}

	var (
		day    = tn.Day()
		hour   = tn.Hour()
		minute = tn.Minute()
		second = tn.Second()
	)

	if timer[TimerIndex_Day] != -1 {
		day = timer[TimerIndex_Day]
	}

	if timer[TimerIndex_Hour] != -1 {
		hour = timer[TimerIndex_Hour]
	}

	if timer[TimerIndex_Minute] != -1 {
		minute = timer[TimerIndex_Minute]
	}

	if timer[TimerIndex_Second] != -1 {
		second = timer[TimerIndex_Second]
	}

	t := time.Date(tn.Year(), month, day, hour, minute, second, 0, loc)

	if timer[TimerIndex_Weekday] != -1 {
		weekday := timer[TimerIndex_Weekday]

		var add int
		if int(tn.Weekday()) < weekday {
			add = weekday - int(tn.Weekday()) - 7

		} else {
			add = weekday - int(tn.Weekday())
		}

		if add != 0 {
			t = t.AddDate(0, 0, add)
		}
	}

	return t
}

func GetExecTime(tn time.Time, loc *time.Location, timer []int) time.Time {
	t := GetBaseExecTime(tn, loc, timer)

	for {
		if t.Unix() <= tn.Unix() {
			if timer[TimerIndex_Minute] == -1 {
				// 每分钟定时
				t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+1, t.Second(), t.Nanosecond(), loc)

			} else if timer[TimerIndex_Hour] == -1 {
				// 每小时定时
				t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, t.Minute(), t.Second(), t.Nanosecond(), loc)

			} else if timer[TimerIndex_Weekday] != -1 {
				// 每周定时,向前推进一周
				t = t.AddDate(0, 0, 7)

			} else if timer[TimerIndex_Day] == -1 {
				// 每日定时,向前推进一天
				t = t.AddDate(0, 0, 1)

			} else if timer[TimerIndex_Month] == -1 {
				// 每月定时,向前推进一月
				t = t.AddDate(0, 1, 0)

			} else if timer[TimerIndex_Month] != -1 {
				// 每年定时,向前推进一年
				t = t.AddDate(1, 0, 0)
			}
		} else {
			break
		}
	}

	return t
}

func ParseTimingString(timing string) ([]int, error) {
	timer := make([]int, TimerLength)

	timeStr := strings.Split(timing, " ")

	if len(timeStr) != TimerLength {
		return timer, errors.New("timing param number not match")
	}

	for index, v := range timeStr {
		if v == "*" {
			timer[index] = -1

		} else {
			i, err := strconv.Atoi(v)
			if err != nil {
				return timer, err
			}

			timer[index] = i
		}
	}

	return timer, nil
}

func NewCrondJob(name string, timing string, f CrondFunc) (*CrondJob, error) {
	cj := &CrondJob{
		Name:  name,
		Func:  f,
		Timer: make([]int, TimerLength),
	}

	timer, err := ParseTimingString(timing)
	if err == nil {
		cj.Timer = timer
	}

	return cj, err
}
