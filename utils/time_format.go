package utils

import (
    "fmt"
    "strings"
    "time"

    log "github.com/sirupsen/logrus"
)

var (
    timeFormat = "2006-01-02 15:04:05"
    nowTime    = time.Now() // 当日时间
)

// GetFirstDateOfMonth 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
}

// GetLastDateOfMonth 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
    return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetStartTime 获取某一天的0点时间
func GetStartTime(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func GetEndTime(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// GetDayStartOrEndTime 获取当天时间的开始和结束
func GetDayStartOrEndTime(d time.Time) (string, string) {
    curTime := time.Unix(d.Unix(), 0).Format(timeFormat)
    todayCurTime := strings.Split(curTime, " ")
    startTime := todayCurTime[0] + " 00:00:00" // 当日开始时间
    endTime := todayCurTime[0] + " 59:59:59"   // 当日结束时间
    return startTime, endTime
}

// GetMonthStartOrEndTime 获取当月时间的开始和结束
func GetMonthStartOrEndTime(d time.Time) (string, string) {
    timeStart := GetFirstDateOfMonth(d)
    timeEnd := GetLastDateOfMonth(d)
    monthTimeStart := time.Unix(timeStart.Unix(), 0).Format("2006-01-02 15:04:05") // 当月开始时间
    monthTimeEnd := time.Unix(timeEnd.Unix(), 0).Format("2006-01-02 15:04:05")
    splitTimeEnd := strings.Split(monthTimeEnd, " ")
    monthEndTime := splitTimeEnd[0] + " 59:59:59" // 当月结束时间
    return monthTimeStart, monthEndTime
}

func TimeAdd() {
    t := nowTime.Unix()
    t += 86400 //增加一天
    curTime := time.Unix(t, 0).Format(timeFormat)
    fmt.Println("当天时间加一天：", curTime)
    timeStart, timeEnd := GetDayStartOrEndTime(nowTime)
    fmt.Println("当天开始时间：", timeStart, "当天结束时间：", timeEnd)
    monthTimeStart, monthTimeEnd := GetMonthStartOrEndTime(nowTime)
    fmt.Println("当月开始时间：", monthTimeStart, "当月结束时间：", monthTimeEnd)
    fmt.Println("当天开始时间：", GetStartTime(nowTime), "当天结束时间：", GetEndTime(nowTime))
}

func GetTimeUnixMilli(timeNow string) int64 {
    t, err := time.ParseInLocation("Jan 2 15:04:05.000", timeNow, time.Local)
    if err != nil {
        log.Errorf("GetTimeUnixMilli 无法解析时间timeNow：%v，Error：%v", timeNow, err)
        return time.Now().UnixMilli()
    }
    time2 := strings.Replace(fmt.Sprintf("%v", t), "0000", fmt.Sprintf("%v", time.Now().Year()), 1)
    if len(time2) > 22 {
        t1, err := time.ParseInLocation("2006-01-02 15:04:05.000", time2[:23], time.Local)
        if err != nil {
            log.Errorf("GetTimeUnixMilli 无法解析时间timeNow：%v，Error：%v", time2[:23], err)
            return time.Now().UnixMilli()
        }
        log.Infof("GetTimeUnixMilli 优化后时间格式输出：%v，转义后时间戳：%v", time2[:23], t1.UnixMilli())
        return t1.UnixMilli()
    }
    return time.Now().UnixMilli()
}

func GetTimeUnix(timeNow string) int64 {
    t, err := time.Parse("Jan 2 15:04:05.000", timeNow)
    if err != nil {
        log.Errorf("GetTimeUnixMilli 无法解析时间：%v", err)
        return 0
    }
    _, month, day := t.Date()
    hour, min, sec := t.Clock()
    timeStamp := time.Date(time.Now().Year(), month, day, hour, min, sec, 0, time.Local).Unix()
    return timeStamp
}

func GetTimeFormat(timeNow string) int64 {
    t, err := time.Parse("2006-01-02 15:04:05", timeNow)
    if err != nil {
        log.Errorf("GetTimeUnixMilli 无法解析时间：%v", err)
        return 0
    }
    return t.Unix()
}
