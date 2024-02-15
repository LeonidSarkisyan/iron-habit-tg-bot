package pagination

type HabitPagination struct {
	Page       int
	ExistsMore bool
}

func NewHabitPagination(existsMore bool) HabitPagination {
	return HabitPagination{
		Page:       1,
		ExistsMore: existsMore,
	}
}
