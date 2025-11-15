package enums

import "strings"

type MangaStatus int

const (
	MANGA_ONGOING MangaStatus = iota
	MANGA_HIATUS
	MANGA_COMPLETED
)

var mangaStatusName = map[MangaStatus]string{
	MANGA_ONGOING:   "ongoing",
	MANGA_HIATUS:    "hiatus",
	MANGA_COMPLETED: "completed",
}

var mangaStatusValue = map[string]MangaStatus{
	"ongoing":   MANGA_ONGOING,
	"hiatus":    MANGA_HIATUS,
	"completed": MANGA_COMPLETED,
}

func (ms MangaStatus) String() string {
	return mangaStatusName[ms]
}

func (ms MangaStatus) StringUpper() string {
	return strings.ToUpper(mangaStatusName[ms])
}

func ParseMangaStatus(s string) (MangaStatus, bool) {
	status, ok := mangaStatusValue[s]
	return status, ok
}
