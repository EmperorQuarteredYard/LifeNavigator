package refreshTime

import "time"

var timeZone *time.Location = time.FixedZone("CST", 8*3600) //默认采用东八区时间

func SetTimeZone(timeZ int) {
	timeZ %= 24
	timeZone = time.FixedZone("CST", timeZ)
}

// ShouldRefresh 认为一天的开始是UTC+8时间4:00，一周的开始为周一
func ShouldRefresh(strategy string, refreshGap uint8, lastRefreshTime time.Time) bool {
	now := time.Now().In(timeZone)
	last := lastRefreshTime.In(timeZone)
	adjustToDayBoundary := func(t time.Time) time.Time {
		return t.Add(-4 * time.Hour)
	}

	switch strategy {
	case NeverCount:
		return false

	case CountDay:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)
		// 计算相差的整天数
		daysDiff := int((nowAdj.Unix() - lastAdj.Unix()) / 86400)
		return daysDiff >= int(refreshGap)

	case CountWeek:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)

		mondayOfWeek := func(t time.Time) time.Time {
			offset := int((t.Weekday() + 6) % 7)
			date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			return date.AddDate(0, 0, -offset)
		}

		lastMonday := mondayOfWeek(lastAdj)
		nowMonday := mondayOfWeek(nowAdj)
		weeksDiff := int((nowMonday.Unix() - lastMonday.Unix()) / (7 * 86400))
		return weeksDiff >= int(refreshGap)

	case CountMonth:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)

		lastYear, lastMonth, lastDay := lastAdj.Date()
		nowYear, nowMonth, nowDay := nowAdj.Date()

		monthDiff := (nowYear-lastYear)*12 + int(nowMonth-lastMonth)
		if nowDay < lastDay {
			monthDiff--
		}
		return monthDiff >= int(refreshGap)

	case CountYear:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)

		lastYear, lastMonth, lastDay := lastAdj.Date()
		nowYear, nowMonth, nowDay := nowAdj.Date()

		yearDiff := nowYear - lastYear
		if nowMonth < lastMonth || (nowMonth == lastMonth && nowDay < lastDay) {
			yearDiff--
		}
		return yearDiff >= int(refreshGap)

	default:
		// 未知策略默认不刷新
		return false
	}
}
