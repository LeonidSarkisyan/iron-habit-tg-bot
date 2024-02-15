package pagination

import "strconv"

type ReminderPagination struct {
	Page       int
	ExistsMore bool
}

func NewReminderPagination(existsMore bool) ReminderPagination {
	return ReminderPagination{
		Page:       1,
		ExistsMore: existsMore,
	}
}

func (rp *ReminderPagination) NextPage() string {
	return strconv.Itoa(rp.Page + 1)
}

func (rp *ReminderPagination) PrevPage() string {
	return strconv.Itoa(rp.Page - 1)
}
