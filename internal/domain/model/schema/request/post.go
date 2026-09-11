package request

import "time"

// Sign is a struct for category or tag
type Sign struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type PostListRequest struct {
	Filter PostFilter     `json:"filter"`
	Data   PostFilterData `json:"data"`
}

type UpdatePostStatusRequest struct {
	PostType string `json:"post_type"`
	ID       string `json:"post_id"`
	Status   string `json:"status"`
}

type PostFilter struct {
	Typ           string          `json:"type"`
	DateRange     *PostFilterDate `json:"date_range,omitempty"`
	Sort          string          `json:"sort"`
	Page          int             `json:"page,omitempty"`
	PageSize      int             `json:"page_size,omitempty"`
	OnlyPublished bool            `json:"-"`
}

type PostFilterData struct {
	Keyword []string `json:"keyword"`
}

type PostFilterDate struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ListOptions is the shared query input for paginated post lists.
// Page/PageSize <= 0 means no pagination (return all matching rows).
type ListOptions struct {
	StartDate      *time.Time
	EndDate        *time.Time
	Page           int
	PageSize       int
	OnlyPublished  bool
	AscendingOrder bool
}

func (o ListOptions) Paginated() bool {
	return o.Page > 0 && o.PageSize > 0
}

func (o ListOptions) Offset() int {
	if !o.Paginated() {
		return 0
	}
	return (o.Page - 1) * o.PageSize
}

func NewListOptions(filter PostFilter) ListOptions {
	opts := ListOptions{
		Page:           filter.Page,
		PageSize:       filter.PageSize,
		OnlyPublished:  filter.OnlyPublished,
		AscendingOrder: filter.Sort == "asc",
	}
	if filter.DateRange == nil {
		return opts
	}
	if filter.DateRange.StartDate != "" {
		if sd, err := time.Parse("2006-01-02 15:04:05", filter.DateRange.StartDate); err == nil {
			opts.StartDate = &sd
		}
	}
	if filter.DateRange.EndDate != "" {
		if ed, err := time.Parse("2006-01-02 15:04:05", filter.DateRange.EndDate); err == nil {
			opts.EndDate = &ed
		}
	}
	return opts
}
