package request

import "testing"

func TestListOptions_Pagination(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		opts := NewListOptions(PostFilter{})
		if opts.Paginated() {
			t.Fatal("expected no pagination")
		}
		if opts.Offset() != 0 {
			t.Fatalf("expected offset 0, got %d", opts.Offset())
		}
	})

	t.Run("page 2 size 10", func(t *testing.T) {
		opts := NewListOptions(PostFilter{Page: 2, PageSize: 10})
		if !opts.Paginated() {
			t.Fatal("expected pagination")
		}
		if opts.Offset() != 10 {
			t.Fatalf("expected offset 10, got %d", opts.Offset())
		}
	})

	t.Run("partial page params disable pagination", func(t *testing.T) {
		opts := NewListOptions(PostFilter{Page: 1, PageSize: 0})
		if opts.Paginated() {
			t.Fatal("expected pagination disabled when page_size is 0")
		}
	})
}

func TestListOptions_DateRange(t *testing.T) {
	opts := NewListOptions(PostFilter{
		DateRange: &PostFilterDate{
			StartDate: "2026-01-01 00:00:00",
			EndDate:   "2026-01-31 23:59:59",
		},
	})
	if opts.StartDate == nil || opts.EndDate == nil {
		t.Fatal("expected both dates parsed")
	}
	if opts.StartDate.Format("2006-01-02") != "2026-01-01" {
		t.Fatalf("unexpected start date %v", opts.StartDate)
	}
	if opts.EndDate.Format("2006-01-02") != "2026-01-31" {
		t.Fatalf("unexpected end date %v", opts.EndDate)
	}
}

func TestListOptions_EmptyDateStringIgnored(t *testing.T) {
	opts := NewListOptions(PostFilter{
		DateRange: &PostFilterDate{StartDate: "", EndDate: ""},
	})
	if opts.StartDate != nil || opts.EndDate != nil {
		t.Fatal("empty date strings should not create pointers")
	}
}
