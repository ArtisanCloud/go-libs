package carbon

import (
	"github.com/golang-module/carbon"
)

type CarbonDatetime struct {
	C *carbon.Carbon
}

var DefaultTimeZone = carbon.UTC

const (
	DATETIME_FORMAT = "Y-m-d H:i:s"
	TIME_FORMAT     = "H:i:s"
)

func CreateCarbonDatetime(c carbon.Carbon) (dt *CarbonDatetime) {
	dt = &CarbonDatetime{
		&c,
	}
	return dt
}

func (dt *CarbonDatetime) SetDatetime(c carbon.Carbon) {
	dt.C = &c
}

func (dt *CarbonDatetime) SetTimezone(timezone string) *CarbonDatetime {
	dt.C.SetTimezone(timezone)
	dt.C.AddHours(8)

	return dt
}

func GetCarbonNow() carbon.Carbon {
	return carbon.Now(DefaultTimeZone)
}
