// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utility

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	LAYOUT_FORMAT19 = "2006-01-02 15:04:05"
	LAYOUT_FORMAT14 = "20060102150405"
	LAYOUT_FORMAT10 = "2006-01-02"
	LAYOUT_FORMAT8  = "20060102"
	LAYOUT_FORMAT3  = "01-02"
	LAYOUT_FORMAT4  = "15:04:05"
)

var (
	loc *time.Location
)

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}

// Format unix time int64 to string
func Date(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return DateT(t, format)
}

// Format unix time string to string
func DateS(ts string, format string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return Date(i, format)
}

// Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func DateT(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

// GetNowTime get current time
func GetNowTime() time.Time {
	return time.Now().In(loc)
}

// GetUnixTime get unix time int转time
func GetUnixTime(timeInt int64) time.Time {
	return time.Unix(timeInt, 0).In(loc)
}

// GetNowTimestamp get current timestamp
func GetNowTimestamp() int {
	return int(time.Now().Unix())
}

// Format19 format date YYYY-MM-DD HH:mm:ss
func Format19(time time.Time) string {
	return time.In(loc).Format(LAYOUT_FORMAT19)
}

// Format14 format date YYYYMMDDHHmmss
func Format14(time time.Time) string {
	return time.Format(LAYOUT_FORMAT14)
}

// Format10 format date YYYY-MM-DD
func Format10(time time.Time) string {
	return time.Format(LAYOUT_FORMAT10)
}

// Format8 format date YYYYMMDD
func Format8(time time.Time) string {
	return time.Format(LAYOUT_FORMAT8)
}

// Format3 format date MM-DD
func Format3(time time.Time) string {
	time1 := time.In(loc)
	return time1.Format(LAYOUT_FORMAT3)
}

// Format4 format date HH:mm:ss
func Format4(time time.Time) string {
	time1 := time.In(loc)
	return time1.Format(LAYOUT_FORMAT4)
}

// TimeIntervalDay 计算两个时间相隔多少天
func TimeIntervalDay(t1, t2 time.Time) float64 {
	result := t1.Sub(t2).Hours() / 24
	if result < 0 {
		result *= -1
	}
	return result
}

// TimeIntervalDay2 计算两个时间相隔多少天,同一天的时间间隔0天
func TimeIntervalDay2(t1, t2 time.Time) int {
	result := int(t1.Sub(t2).Hours() / 24)
	if result < 0 {
		result *= -1
	}
	return result
}

// TimeIntervalSecond 两个时间相隔秒数
func TimeIntervalSecond(time1 time.Time, time2 time.Time) int {
	d := time2.Sub(time1)
	return int(d / 1e9)
}

// Format19ToTimestamp 格式化的时间(YYYY-MM-DD HH:mm:ss)转时间戳
func Format19ToTimestamp(formatTime string) (int, error) {
	if len(formatTime) != 19 {
		return 0, errors.New("formatTime length error")
	} else {
		sTemp, err := time.ParseInLocation(LAYOUT_FORMAT19, formatTime, loc)
		if err != nil {
			return 0, err
		}
		return int(sTemp.Unix()), nil
	}
}

// Format14ToTimestamp 格式化的时间(YYYYMMDDHHmmss)转时间戳
func Format14ToTimestamp(formatTime string) (int, error) {
	if len(formatTime) != 14 {
		return 0, errors.New("formatTime length error")
	} else {
		sTemp, err := time.ParseInLocation(LAYOUT_FORMAT14, formatTime, loc)
		if err != nil {
			return 0, err
		}
		return int(sTemp.Unix()), nil
	}
}

// Format10ToTimestamp 格式化的时间(YYYY-MM-DD)转时间戳
func Format10ToTimestamp(formatTime string) (int, error) {
	if len(formatTime) != 10 {
		return 0, errors.New("formatTime length error")
	} else {
		sTemp, err := time.ParseInLocation(LAYOUT_FORMAT10, formatTime, loc)
		if err != nil {
			return 0, err
		}
		return int(sTemp.Unix()), nil
	}
}

// Format8ToTimestamp 格式化的时间(YYYYMMDD)转时间戳
func Format8ToTimestamp(formatTime string) (int, error) {
	if len(formatTime) != 8 {
		return 0, errors.New("formatTime length error")
	} else {
		sTemp, err := time.ParseInLocation(LAYOUT_FORMAT8, formatTime, loc)
		if err != nil {
			return 0, err
		}
		return int(sTemp.Unix()), nil
	}
}

// GetTimeIntervalDay 与当前时间的时间间隔 num 间隔的数量 types 间隔的单位 1天 2月 3年 默认间隔单位为天
func GetTimeIntervalDay(num int, types ...int) int {
	var days, months, year int
	if len(types) > 0 {
		if types[0] == 2 {
			months = num
		} else if types[0] == 3 {
			year = num
		} else {
			days = num
		}
	} else {
		days = num
	}
	return int(time.Now().AddDate(year, months, days).Unix())
}

// 与传入时间的时间间隔 num 间隔的数量 types 间隔的单位 1天 2月 3年 默认间隔单位为天
func GetTimeIntervalDayByTimes(times time.Time, num int, types ...int) int {
	var days, months, year int
	if len(types) > 0 {
		if types[0] == 2 {
			months = num
		} else if types[0] == 3 {
			year = num
		} else {
			days = num
		}
	} else {
		days = num
	}
	return int(times.In(loc).AddDate(year, months, days).Unix())
}

// GetTimeIntervalDay 与当前时间的时间间隔 num 间隔的数量 types 间隔的单位 1时 2分 3秒 默认间隔单位为天
func GetTimeIntervalHour(num int, types ...int) int {
	if len(types) > 0 {
		switch types[0] {
		case 1:
			return int(time.Now().Add(-time.Hour * 10).Unix()) // 时
		case 2:
			return int(time.Now().Add(-time.Minute * 10).Unix()) // 分
		case 3:
			return int(time.Now().Add(-time.Second * 10).Unix()) // 秒
		default:
			return int(time.Now().Add(-time.Hour * 10).Unix()) // 时
		}
	}
	return int(time.Now().Add(-time.Hour * 10).Unix()) // 时
}

// GetNightTimestamp 获取前一天的时间戳：day = -1 , 获取今天晚上的时间戳：0 ，获取明天的凌晨时间戳：1
func GetNightTimestamp(day int) int {
	nTime := time.Now().In(loc)
	yesTime := nTime.AddDate(0, 0, day)
	d := yesTime.Format(LAYOUT_FORMAT8)
	res, _ := Format8ToTimestamp(d)
	return res
}

var ChinaZone = time.FixedZone("chinaZone", 8*60*60)

// GetUnixTimeIntervalByDay return unix time interval by day
// offsetDay -1 means yesterday
// offsetDay 0 means current
// offsetDay +1 means tomorrow
// so, .e.g offsetDay +2, offsetDay +3, offsetDay -2, offsetDay -3
func GetUnixTimeIntervalByDay(offsetDay int) []time.Time {
	t := time.Now().In(ChinaZone)

	const day = time.Hour * 24
	nowBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ChinaZone)
	nowEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, ChinaZone)

	return []time.Time{
		nowBegin.Add(day * time.Duration(offsetDay)),
		nowEnd.Add(day * time.Duration(offsetDay)),
	}
}

// GetUnixTimeIntervalByWeek return unix time interval by week
// offsetWeek -1 means last week
// offsetWeek 0 means this week
// offsetWeek +1 means next week
// so, .e.g offsetWeek +2, offsetWeek +3, offsetWeek -2, offsetWeek -3
func GetUnixTimeIntervalByWeek(offsetWeek int) []time.Time {
	t := time.Now().In(ChinaZone)

	const week = time.Hour * 24 * 7
	nowBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ChinaZone)
	nowEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, ChinaZone)

	var begin, end time.Duration
	if t.Weekday() != time.Sunday {
		begin = time.Duration((t.Weekday() - time.Monday)) * time.Hour * -24
		end = time.Duration((time.Saturday - t.Weekday() + 1)) * time.Hour * 24
	} else {
		begin = 6 * time.Hour * -24
		end = 0
	}
	return []time.Time{
		nowBegin.Add(begin + week*time.Duration(offsetWeek)),
		nowEnd.Add(end + week*time.Duration(offsetWeek)),
	}
}

// GetUnixTimeIntervalByHour return unix time interval by hour
// offsetHour -1 means last hour
// offsetHour 0 means this hour
// offsetHour +1 means next hour
// so, .e.g offsetHour +2, offsetHour +3, offsetHour -2, offsetHour -3
func GetUnixTimeIntervalByHour(offsetHour int) []time.Time {
	t := time.Now().In(ChinaZone)

	const hour = time.Hour
	nowBegin := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, ChinaZone)
	nowEnd := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 0, ChinaZone)

	return []time.Time{
		nowBegin.Add(hour * time.Duration(offsetHour)),
		nowEnd.Add(hour * time.Duration(offsetHour)),
	}
}

// GetUnixTimeIntervalByMonth return unix time interval by month
// offsetMonth -1 means last month
// offsetMonth 0 means this month
// offsetMonth +1 means next month
// so, .e.g offsetMonth +2, offsetMonth +3, offsetMonth -2, offsetMonth -3
func GetUnixTimeIntervalByMonth(offsetMonth int) []time.Time {
	t := time.Now().In(ChinaZone)
	nowBegin := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, ChinaZone)
	return []time.Time{
		nowBegin.AddDate(0, offsetMonth, 0),
		nowBegin.AddDate(0, offsetMonth+1, 0).Add(-1 * time.Second),
	}
}

// 传入日期  YYYY-MM-DD  获取当天23:59:59的时间戳
func GetTimeSomeOneDayNext(offset string) (int, error) {
	t, err := time.ParseInLocation(LAYOUT_FORMAT19, offset+" 23:59:59", time.Local)
	return int(t.Unix()), err
}

// 前两小时的时间戳
func TwoHourBefore() int {
	timeNow := GetNowTime()
	h, _ := time.ParseDuration("-1h")
	return int(timeNow.Add(2 * h).Unix())
}

// 今日零点时间戳
func ZeroHourByToday() int {
	s := GetNowTime().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", s, loc)
	return int(t.Unix())
}

func TimeZoneConvert(unixtime int64, offset string) (string, error) {
	t := time.Unix(unixtime, 0)
	d, err := time.ParseDuration(offset)
	if err != nil {
		return "", err
	}

	loc, err := time.LoadLocation("Asia/Shanghai") //北京时间-12小时
	if err != nil {
		return "", err
	}

	return t.In(loc).Add(d).Format("2006-01-02 15:04:05"), nil
}

func GetDateAndWeekByUnix(unixTime int64) (string, string) {
	chinaTime := time.Unix(unixTime, 0).In(ChinaZone)
	var date, week string
	y, m, d := chinaTime.Date()
	date = fmt.Sprintf("%d-%02d-%02d", y, m, d)

	switch chinaTime.Weekday() {
	case 0:
		week = "周日"
	case 1:
		week = "周一"
	case 2:
		week = "周二"
	case 3:
		week = "周三"
	case 4:
		week = "周四"
	case 5:
		week = "周五"
	case 6:
		week = "周六"
	}
	return date, week
}

// 获取上个月第一天0时间戳
func GetLastMonthFirstDay() int64 {
	times := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc)
	return times.AddDate(0, -1, 0).Unix()
}

// 获取当月第一天0时间戳
func GetMonthFirstDay() int64 {
	times := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc)
	return times.Unix()
}

// 获取本周第一天0点时间戳
func FirstDayByCurrentWeek() time.Time {
	s := strconv.Itoa(-(int(time.Now().Weekday()-1) * 24)) + "h"
	d, _ := time.ParseDuration(s)
	m := time.Now().Add(d)
	return time.Date(m.Year(), m.Month(), m.Day(), 0, 0, 0, 0, loc)
}

// 获取上周第一天0点时间戳
func FirstDayByPreWeek() time.Time {
	d, _ := time.ParseDuration(strconv.Itoa(-7*24) + "h")
	return FirstDayByCurrentWeek().Add(d)
}

func MonthDay(d time.Time) string {
	return d.In(loc).Format("01-02")
}

func YearDay(d time.Time) string {
	return d.In(loc).Format("2006-01-02")
}

// 昨天零点时间
func ZeroHourByLastday() int {
	s := GetNowTime().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", s, loc)
	d, err := time.ParseDuration("-24h")
	if err != nil {
		return 0
	}
	return int(t.Add(d).Unix())
}

// 将时间转换为时间戳
func ConvertTime(s string) int {
	t, err := time.ParseInLocation(LAYOUT_FORMAT19, s, loc)
	if err != nil {
		return 0
	}
	return int(t.Unix())
}

// 时间差额对比 开始时间和结束时间差距的天数对比
// sTime 开始时间
// eTime 结束时间
// num 距离数量
// timeType 单位 1.天 2.月 3.年 默认是天
func TimeDifferenceComparison(sTime, eTime, num int, timeType ...int) bool {
	t := time.Unix(int64(sTime), 0)
	if len(timeType) > 0 {
		switch timeType[0] {
		case 1: // 天
			return t.AddDate(0, 0, num).Unix() > int64(eTime)
		case 2: // 月
			return t.AddDate(0, num, 0).Unix() > int64(eTime)
		case 3: // 年
			return t.AddDate(num, 0, 0).Unix() > int64(eTime)
		}
	}
	return t.AddDate(0, 0, num).Unix() > int64(eTime)
}

//获取统计时间段的每一天
func GetEveryDay(startTime, EndTime int) []string {
	var date []string
	if startTime == EndTime {
		tm := time.Unix(int64(startTime), 0)
		str := YearDay(tm)
		date = append(date, str)
	} else {
		num := (EndTime - startTime) / (24 * 60 * 60)
		for j := 0; j <= num; j++ {
			tm := time.Unix(int64(startTime+24*60*60*j), 0)
			str := YearDay(tm)
			date = append(date, str)
		}
	}
	return date
}