package wintime

import "time"

const (
	// 100-nanosecond intervals between January 1, 1601 and January 1, 1970
	intervalsBetween1601And1970 = int64(116444736000000000)
)

// ToTime converts low and high DWORD bits of a Windows FILETIME or similar
// format (that represents the number of 100-nanosecond intervals since
// January 1, 1601 (UTC)) to a time.Time object.
func ToTime(low, high int64) time.Time {
	full := (high << 32) | (low & 0xFFFFFFFF)              // get a full 64-bit FILETIME
	unixNano := (full - intervalsBetween1601And1970) * 100 // get unix nanoseconds
	return time.Unix(0, unixNano)
}
