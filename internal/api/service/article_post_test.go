package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"errors"
	"testing"
)

func TestArticleService_DeleteArticleByID_Cascade(t *testing.T) {
	articles := &mockArticleRepo{}
	tags := &mockTagRepo{}
	categoryLinks := &mockCategoryLinkRepo{}

	svc := NewArticleService(articles, &mockCategoryRepo{}, categoryLinks, tags, nil)
	if err := svc.DeleteArticleByID("a1"); err != nil {
		t.Fatalf("DeleteArticleByID returned error: %v", err)
	}

	if len(articles.deleteCalls) != 1 {
		t.Fatalf("expected article deleted once, got %d calls: %v", len(articles.deleteCalls), articles.deleteCalls)
	}
	if articles.deleteCalls[0] != "a1" {
		t.Fatalf("expected delete id a1, got %s", articles.deleteCalls[0])
	}
	if len(tags.deleteLinkCalls) != 1 || tags.deleteLinkCalls[0].targetID != "a1" || tags.deleteLinkCalls[0].targetType != table.ArticleTag {
		t.Fatalf("expected tag links cleaned for a1, got %+v", tags.deleteLinkCalls)
	}
	if len(categoryLinks.deleteByArticleCalls) != 1 || categoryLinks.deleteByArticleCalls[0] != "a1" {
		t.Fatalf("expected category links cleaned for a1, got %v", categoryLinks.deleteByArticleCalls)
	}
}

func TestArticleService_DeleteArticleByID_PropagatesError(t *testing.T) {
	want := errors.New("db down")
	articles := &mockArticleRepo{deleteErr: want}
	svc := NewArticleService(articles, &mockCategoryRepo{}, &mockCategoryLinkRepo{}, &mockTagRepo{}, nil)
	if err := svc.DeleteArticleByID("a1"); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestSharingService_UpdateSharing_Private(t *testing.T) {
	cases := []struct {
		name     string
		existing int8
		private  bool
		want     int8
	}{
		{name: "public from draft", existing: 0, private: false, want: 1},
		{name: "private from public", existing: 1, private: true, want: 0},
		{name: "edit keeps archive", existing: 2, private: false, want: 2},
		{name: "private overwrites archive", existing: 2, private: true, want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sharings := &mockSharingRepo{
				existing: &table.SharingTable{PostBase: table.PostBase{ID: "s1", Status: tc.existing}},
			}
			svc := NewSharingService(sharings, &mockTagRepo{}, &mockAttachmentRepo{}, &mockAttachmentLinkRepo{})
			err := svc.UpdateSharing(&request.UpdateSharingRequest{
				CreateSharingRequest: request.CreateSharingRequest{
					Content: "hello",
					Author:  "u1",
					Private: tc.private,
				},
				ID: "s1",
			})
			if err != nil {
				t.Fatalf("UpdateSharing returned error: %v", err)
			}
			if sharings.updated == nil {
				t.Fatal("expected sharing update call")
			}
			if sharings.updated.Status != tc.want {
				t.Fatalf("expected status %d, got %d", tc.want, sharings.updated.Status)
			}
		})
	}
}

func TestPostService_UpdatePostStatus_PendingOnlyMessage(t *testing.T) {
	t.Run("message allows pending", func(t *testing.T) {
		posts := &mockPostRepo{}
		svc := NewPostService(&mockArticleRepo{}, &mockSharingRepo{}, &mockMessageRepo{}, posts)
		err := svc.UpdatePostStatus(&request.UpdatePostStatusRequest{
			PostType: "message",
			ID:       "m1",
			Status:   "pending",
		})
		if err != nil {
			t.Fatalf("expected pending allowed for message, got %v", err)
		}
		if posts.lastStatus != 3 {
			t.Fatalf("expected status 3, got %d", posts.lastStatus)
		}
	})

	t.Run("article rejects pending", func(t *testing.T) {
		posts := &mockPostRepo{}
		svc := NewPostService(&mockArticleRepo{}, &mockSharingRepo{}, &mockMessageRepo{}, posts)
		err := svc.UpdatePostStatus(&request.UpdatePostStatusRequest{
			PostType: "article",
			ID:       "a1",
			Status:   "pending",
		})
		if err == nil {
			t.Fatal("expected error for article pending")
		}
	})

	t.Run("publish works", func(t *testing.T) {
		posts := &mockPostRepo{}
		svc := NewPostService(&mockArticleRepo{}, &mockSharingRepo{}, &mockMessageRepo{}, posts)
		err := svc.UpdatePostStatus(&request.UpdatePostStatusRequest{
			PostType: "article",
			ID:       "a1",
			Status:   "publish",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if posts.lastStatus != 1 {
			t.Fatalf("expected status 1, got %d", posts.lastStatus)
		}
	})
}

func TestPostService_SingleListWithPost_PaginationMeta(t *testing.T) {
	articles := &mockArticleRepo{
		articles:     []*table.ArticleTable{{PostBase: table.PostBase{ID: "a1"}}},
		articleTotal: 42,
	}
	svc := NewPostService(articles, &mockSharingRepo{}, &mockMessageRepo{}, &mockPostRepo{})
	res := svc.SingleListWithPost(&request.PostListRequest{
		Filter: request.PostFilter{
			Typ:      "single",
			Page:     2,
			PageSize: 10,
		},
	}, "article")
	if res.Total != 42 {
		t.Fatalf("expected total 42, got %d", res.Total)
	}
	if res.Page != 2 || res.PageSize != 10 {
		t.Fatalf("expected page meta 2/10, got %d/%d", res.Page, res.PageSize)
	}
	if len(res.Articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(res.Articles))
	}
}
