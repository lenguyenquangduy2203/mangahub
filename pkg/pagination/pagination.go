package pagination

const (
	DEFAULT_LIMIT = 50
	MAX_LIMIT     = 100
)

type Paginated interface {
	GetTotal() int
	GetLimit() int
	GetOffset() int
	GetResults() any
}

func CalculateOffset(page int, limit int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * limit
}
