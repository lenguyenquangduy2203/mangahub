package enums

import "strings"

type ReadingStatus int

const (
	READING_PLAN_TO_READ ReadingStatus = iota
	READING_READING
	READING_COMPLETED
)

var readingStatusName = map[ReadingStatus]string{
	READING_PLAN_TO_READ: "plan_to_read",
	READING_READING:      "reading",
	READING_COMPLETED:    "completed",
}

var readingStatusValue = map[string]ReadingStatus{
	"plan_to_read": READING_PLAN_TO_READ,
	"reading":      READING_READING,
	"completed":    READING_COMPLETED,
}

func (rs ReadingStatus) String() string {
	return readingStatusName[rs]
}

func (rs ReadingStatus) StringUpper() string {
	return strings.ToUpper(readingStatusName[rs])
}

func ParseReadingStatus(s string) (ReadingStatus, bool) {
	status, ok := readingStatusValue[s]
	return status, ok
}
