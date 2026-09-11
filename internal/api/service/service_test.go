package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "mou1ght-service-test")
	if err != nil {
		panic(err)
	}
	cfg := []byte(`blog:
  board:
    need_reviewed: false
database:
  dsn: test
  type: sqlite
security:
  jwt_key: test-jwt
  visitor_jwt_key: test-visitor
`)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), cfg, 0o644); err != nil {
		panic(err)
	}
	old, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
	code := m.Run()
	_ = os.Chdir(old)
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func TestMessageService_UpdateMessage_Ownership(t *testing.T) {
	t.Run("forbidden when author mismatch", func(t *testing.T) {
		repo := &mockMessageRepo{
			getMessage: &table.MessageTable{
				PostBase: table.PostBase{ID: "m1", Status: 1},
				AuthorIP: "owner-jti",
			},
		}
		svc := NewMessageService(repo)
		err := svc.UpdateMessage(&request.UpdateMessageRequest{
			CreateMessageRequest: request.CreateMessageRequest{
				Content:      "x",
				VisitorToken: "intruder-jti",
			},
			ID: "m1",
		})
		var fe *fiber.Error
		if !errors.As(err, &fe) || fe.Code != 403 {
			t.Fatalf("expected 403, got %v", err)
		}
		if repo.updated != nil {
			t.Fatal("update should not be called when forbidden")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockMessageRepo{getMessageErr: gorm.ErrRecordNotFound}
		svc := NewMessageService(repo)
		err := svc.UpdateMessage(&request.UpdateMessageRequest{
				CreateMessageRequest: request.CreateMessageRequest{VisitorToken: "jti"},
			ID:                   "missing",
		})
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("expected record not found, got %v", err)
		}
	})

	t.Run("owner can update", func(t *testing.T) {
		repo := &mockMessageRepo{
			getMessage: &table.MessageTable{
				PostBase: table.PostBase{ID: "m1", Status: 1},
				AuthorIP: "owner-jti",
			},
		}
		svc := NewMessageService(repo)
		err := svc.UpdateMessage(&request.UpdateMessageRequest{
			CreateMessageRequest: request.CreateMessageRequest{
				Content:      "updated",
				VisitorToken: "owner-jti",
				Position:     request.MessagePosition{X: 1, Y: 2, Z: 3},
			},
			ID: "m1",
		})
		if err != nil {
			t.Fatalf("owner update failed: %v", err)
		}
		if repo.updated == nil || repo.updated.Content != "updated" {
			t.Fatalf("expected content updated, got %+v", repo.updated)
		}
		if repo.updated.Status != 1 {
			t.Fatalf("expected status preserved as 1, got %d", repo.updated.Status)
		}
	})
}

func TestAttachmentService_Delete(t *testing.T) {
	t.Run("requires id", func(t *testing.T) {
		svc := NewAttachmentService(&mockAttachmentRepo{}, &mockAttachmentLinkRepo{})
		if err := svc.Delete(""); !errors.Is(err, ErrAttachmentIDRequired) {
			t.Fatalf("expected ErrAttachmentIDRequired, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := NewAttachmentService(
			&mockAttachmentRepo{getAttachmentErr: gorm.ErrRecordNotFound},
			&mockAttachmentLinkRepo{},
		)
		if err := svc.Delete("x"); !errors.Is(err, ErrAttachmentNotFound) {
			t.Fatalf("expected ErrAttachmentNotFound, got %v", err)
		}
	})

	t.Run("in use", func(t *testing.T) {
		svc := NewAttachmentService(
			&mockAttachmentRepo{getAttachment: &table.AttachmentTable{ID: "x", StoragePath: "a/b.png"}},
			&mockAttachmentLinkRepo{refCount: 2},
		)
		if err := svc.Delete("x"); !errors.Is(err, ErrAttachmentInUse) {
			t.Fatalf("expected ErrAttachmentInUse, got %v", err)
		}
	})

	t.Run("deletes unused", func(t *testing.T) {
		atts := &mockAttachmentRepo{
			getAttachment: &table.AttachmentTable{ID: "x", StoragePath: "image/2026/01/x.png"},
		}
		svc := NewAttachmentService(atts, &mockAttachmentLinkRepo{refCount: 0})
		if err := svc.Delete("x"); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		if len(atts.deleted) != 1 || atts.deleted[0] != "x" {
			t.Fatalf("expected attachment deleted, got %v", atts.deleted)
		}
	})
}

func TestUserService_ChangePassword(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	t.Run("wrong old password", func(t *testing.T) {
		users := &mockUserRepo{
			getUser: &table.UserTable{ID: "u1", Password: string(hashed)},
		}
		svc := NewUserService(users, &mockArticleRepo{}, &mockSharingRepo{})
		err := svc.ChangePassword("u1", &request.ChangePasswordRequest{
			OldPassword: "bad",
			NewPassword: "newpass",
		})
		if err == nil || err.Error() != "old password incorrect" {
			t.Fatalf("expected old password incorrect, got %v", err)
		}
	})

	t.Run("same password rejected", func(t *testing.T) {
		users := &mockUserRepo{
			getUser: &table.UserTable{ID: "u1", Password: string(hashed)},
		}
		svc := NewUserService(users, &mockArticleRepo{}, &mockSharingRepo{})
		err := svc.ChangePassword("u1", &request.ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "oldpass",
		})
		if err == nil {
			t.Fatal("expected same password rejected")
		}
	})

	t.Run("success", func(t *testing.T) {
		users := &mockUserRepo{
			getUser: &table.UserTable{ID: "u1", Password: string(hashed)},
		}
		svc := NewUserService(users, &mockArticleRepo{}, &mockSharingRepo{})
		if err := svc.ChangePassword("u1", &request.ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "newpass",
		}); err != nil {
			t.Fatalf("change password failed: %v", err)
		}
		if users.updatedPwd == "" {
			t.Fatal("expected password update call")
		}
		if bcrypt.CompareHashAndPassword([]byte(users.updatedPwd), []byte("newpass")) != nil {
			t.Fatal("updated hash does not match new password")
		}
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Run("username conflict", func(t *testing.T) {
		users := &mockUserRepo{
			getUser: &table.UserTable{ID: "u1", UserName: "alice", Email: "a@x.com"},
			byName: map[string]*table.UserTable{
				"bob": {ID: "u2", UserName: "bob"},
			},
		}
		svc := NewUserService(users, &mockArticleRepo{}, &mockSharingRepo{})
		_, err := svc.UpdateProfile("u1", &request.UpdateUserProfileRequest{UserName: "bob"})
		if err == nil || err.Error() != "user name already exists" {
			t.Fatalf("expected username conflict, got %v", err)
		}
	})

	t.Run("updates changed fields only", func(t *testing.T) {
		users := &mockUserRepo{
			getUser: &table.UserTable{ID: "u1", UserName: "alice", Email: "old@x.com"},
		}
		svc := NewUserService(users, &mockArticleRepo{}, &mockSharingRepo{})
		info, err := svc.UpdateProfile("u1", &request.UpdateUserProfileRequest{
			Email: "new@x.com",
		})
		if err != nil {
			t.Fatalf("update profile failed: %v", err)
		}
		if info == nil {
			t.Fatal("expected user entity")
		}
		if users.profile == nil {
			t.Fatal("expected profile update")
		}
		if _, ok := users.profile["email"]; !ok {
			t.Fatalf("expected email field, got %v", users.profile)
		}
		if _, ok := users.profile["user_name"]; ok {
			t.Fatalf("username should not be updated, got %v", users.profile)
		}
	})
}

func TestTagService_UpdateTag(t *testing.T) {
	t.Run("empty label", func(t *testing.T) {
		svc := NewTagService(&mockTagRepo{})
		if err := svc.UpdateTag("t1", &request.UpdateTagRequest{Label: ""}); err == nil {
			t.Fatal("expected empty label error")
		}
	})

	t.Run("success", func(t *testing.T) {
		tags := &mockTagRepo{getTag: &table.TagTable{ID: "t1", Label: "old"}}
		svc := NewTagService(tags)
		if err := svc.UpdateTag("t1", &request.UpdateTagRequest{Label: "new"}); err != nil {
			t.Fatalf("update tag failed: %v", err)
		}
		if tags.updated == nil || tags.updated.Label != "new" {
			t.Fatalf("expected label new, got %+v", tags.updated)
		}
	})
}
