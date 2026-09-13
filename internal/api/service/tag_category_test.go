package service

import (
	"errors"
	"Mou1ght/internal/domain/model/table"
	"testing"

	"Mou1ght/internal/domain/entity"
)

func findGroupByID(groups []*entity.CategoryGroup, id string) *entity.CategoryGroup {
	for _, g := range groups {
		if g.ID == id {
			return g
		}
		if found := findGroupByID(g.Children, id); found != nil {
			return found
		}
	}
	return nil
}

func TestTagService_TagsList_FillsArticleCount(t *testing.T) {
	repo := &mockTagRepo{
		tags: []table.TagTable{
			{ID: "t1", Label: "go"},
			{ID: "t2", Label: "rust"},
		},
		counts: map[string]int64{"t1": 3},
	}
	svc := NewTagService(repo)
	signs := svc.TagsList()

	if len(signs) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(signs))
	}
	if signs[0].ID != "t1" || signs[0].Count != 3 {
		t.Fatalf("expected t1 count 3, got %+v", signs[0])
	}
	if signs[1].ID != "t2" || signs[1].Count != 0 {
		t.Fatalf("expected t2 count 0, got %+v", signs[1])
	}
}

func TestTagService_TagsList_CountErrorStillReturnsList(t *testing.T) {
	repo := &mockTagRepo{
		tags:      []table.TagTable{{ID: "t1", Label: "go"}},
		countsErr: errors.New("count failed"),
	}
	svc := NewTagService(repo)
	signs := svc.TagsList()
	if len(signs) != 1 || signs[0].Count != 0 {
		t.Fatalf("expected list with zero count on count error, got %+v", signs)
	}
}

func TestDTOService_GetCategoryGroupFromTables_AggregatesTotalCount(t *testing.T) {
	cr := &mockCategoryRepo{counts: map[string]int64{"a": 2, "b": 1}}
	dto := NewDTOService(nil, nil, nil, cr, nil, nil, nil)

	groups := dto.GetCategoryGroupFromTables([]table.CategoryTable{
		{ID: "a", Label: "root-a"},
		{ID: "b", Label: "child-b", ParentID: "a"},
		{ID: "c", Label: "root-c"},
	})

	if len(groups) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(groups))
	}

	a := findGroupByID(groups, "a")
	if a == nil {
		t.Fatal("root a not found")
	}
	if a.Count != 2 || a.TotalCount != 3 {
		t.Fatalf("expected a count=2 total=3, got count=%d total=%d", a.Count, a.TotalCount)
	}
	if len(a.Children) != 1 || a.Children[0].ID != "b" {
		t.Fatalf("expected a to have child b, got %+v", a.Children)
	}

	b := findGroupByID(groups, "b")
	if b == nil || b.Count != 1 || b.TotalCount != 1 {
		t.Fatalf("expected b count=1 total=1, got %+v", b)
	}

	c := findGroupByID(groups, "c")
	if c == nil || c.Count != 0 || c.TotalCount != 0 {
		t.Fatalf("expected c count=0 total=0, got %+v", c)
	}
}
