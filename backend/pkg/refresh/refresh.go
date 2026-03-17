package refresh

import "time"

const (
	TypeRefreshNever   = 0
	TypeRefreshYearly  = 1
	TypeRefreshMonthly = 2
	TypeRefreshWeekly  = 3
	TypeRefreshDaily   = 4
	TypeRefreshHourly  = 5
)

const (
	NameRefreshNever   = "Never"
	NameRefreshYearly  = "Yearly"
	NameRefreshMonthly = "Monthly"
	NameRefreshWeekly  = "Weekly"
	NameRefreshDaily   = "Daily"
	NameRefreshHourly  = "Hourly"
)

var typeNameIndex = map[int]string{
	TypeRefreshNever:   NameRefreshNever,
	TypeRefreshYearly:  NameRefreshYearly,
	TypeRefreshMonthly: NameRefreshMonthly,
	TypeRefreshWeekly:  NameRefreshWeekly,
	TypeRefreshDaily:   NameRefreshDaily,
	TypeRefreshHourly:  NameRefreshHourly,
}

var nameTypeIndex = map[string]int{
	NameRefreshNever:   TypeRefreshNever,
	NameRefreshYearly:  TypeRefreshYearly,
	NameRefreshMonthly: TypeRefreshMonthly,
	NameRefreshWeekly:  TypeRefreshWeekly,
	NameRefreshDaily:   TypeRefreshDaily,
	NameRefreshHourly:  TypeRefreshHourly,
}

var timeZone *time.Location = time.FixedZone("CST", 8*3600) //默认采用东八区时间

func SetTimeZone(timeZ int) {
	timeZ %= 24
	if timeZ < 0 {
		timeZ += 24
	}
	timeZone = time.FixedZone("CST", timeZ)
}
func ComposeRefreshInterval(refreshGap, strategy int) uint32 {
	return uint32(refreshGap<<8 + strategy&0xFF)
}

func CreateRefreshInterval(refreshGap uint32, strategy string) uint32 {
	Type, ok := nameTypeIndex[strategy]
	if !ok {
		return 0
	}
	return ComposeRefreshInterval(int(refreshGap), int(Type))
}

func ReduceRefreshInterval(refreshInterval uint32) (refreshGap uint32, strategy int) {
	return refreshInterval >> 8, int(refreshInterval & 0xFF)
}

func AnalysisRefreshInterval(refreshInterval uint32) (uint32, string) {
	Gap, Type := ReduceRefreshInterval(refreshInterval)
	Strategy, ok := typeNameIndex[Type]
	if !ok {
		return 0, NameRefreshNever
	}
	return Gap, Strategy
}

// ShouldRefresh 认为一天的开始是UTC+8时间4:00，一周的开始为周一
func ShouldRefresh(RefreshInterval uint32, lastRefreshTime time.Time) bool {
	strategy := RefreshInterval & 0xFF
	refreshGap := RefreshInterval >> 8
	now := time.Now().In(timeZone)
	last := lastRefreshTime.In(timeZone)
	adjustToDayBoundary := func(t time.Time) time.Time {
		return t.Add(-4 * time.Hour)
	}

	switch strategy {
	case TypeRefreshNever:
		return false

	case TypeRefreshDaily:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)
		// 计算相差的整天数
		daysDiff := int((nowAdj.Unix() - lastAdj.Unix()) / 86400)
		return daysDiff >= int(refreshGap)

	case TypeRefreshWeekly:
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

	case TypeRefreshMonthly:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)

		lastYear, lastMonth, lastDay := lastAdj.Date()
		nowYear, nowMonth, nowDay := nowAdj.Date()

		monthDiff := (nowYear-lastYear)*12 + int(nowMonth-lastMonth)
		if nowDay < lastDay {
			monthDiff--
		}
		return monthDiff >= int(refreshGap)

	case TypeRefreshYearly:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)

		lastYear, lastMonth, lastDay := lastAdj.Date()
		nowYear, nowMonth, nowDay := nowAdj.Date()

		yearDiff := nowYear - lastYear
		if nowMonth < lastMonth || (nowMonth == lastMonth && nowDay < lastDay) {
			yearDiff--
		}
		return yearDiff >= int(refreshGap)
	case TypeRefreshHourly:
		lastAdj := adjustToDayBoundary(last)
		nowAdj := adjustToDayBoundary(now)
		hourDiff := (nowAdj.Unix() - lastAdj.Unix()) / 3600
		return hourDiff >= int64(refreshGap)

	default:
		// 未知策略默认不刷新
		return false
	}
}
