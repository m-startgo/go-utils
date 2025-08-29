package mtime

//  依据 生日 计算年龄

func Age(birthDay string) int {
	birthDayTime, err := Parse(birthDay)
	if err != nil {
		return 0
	}
	startUnix := birthDayTime.UnixMilli()

	nowUnix := Now().UnixMilli()

	if nowUnix-startUnix < 0 {
		return 0
	}
	year := (nowUnix - startUnix) / (1000 * 60 * 60 * 24 * 365)

	return int(year)
}
