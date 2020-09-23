package variable

import (
	"binadesa2020-backend/lib/clog"
	"time"
)

// DateTimeNowPtr to pointer
func DateTimeNowPtr() *time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		clog.Fatal(err, "Get location clock")
	}
	now := time.Now().In(loc)
	return &now
}
