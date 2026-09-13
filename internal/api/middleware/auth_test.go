package middleware

import (
	"io"
	"Mou1ght/internal/pkg/util"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "mou1ght-middleware-test")
	if err != nil {
		panic(err)
	}
	cfg := []byte(`security:
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

func newAuthApp() *fiber.App {
	app := fiber.New()
	app.Get("/protected", Auth, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"uid": c.Locals("uid")})
	})
	return app
}

func doAuthRequest(t *testing.T, app *fiber.App, authHeader string) (*http.Response, error) {
	t.Helper()
	req := httptest.NewRequest("GET", "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	resp, err := app.Test(req, -1)
	return resp, err
}

func TestAuth(t *testing.T) {
	app := newAuthApp()
	validToken, err := util.ReleaseToken("admin-1")
	if err != nil {
		t.Fatalf("release token: %v", err)
	}

	t.Run("missing authorization header returns 403", func(t *testing.T) {
		resp, err := doAuthRequest(t, app, "")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 403 {
			t.Fatalf("expected 403, got %d", resp.StatusCode)
		}
	})

	t.Run("non-bearer scheme returns 403", func(t *testing.T) {
		resp, err := doAuthRequest(t, app, "Token abc")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 403 {
			t.Fatalf("expected 403, got %d", resp.StatusCode)
		}
	})

	t.Run("garbage token returns 401", func(t *testing.T) {
		resp, err := doAuthRequest(t, app, "Bearer not-a-jwt")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 401 {
			t.Fatalf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("token signed with visitor key returns 401", func(t *testing.T) {
		visitorToken, err := util.ReleaseVisitorToken("127.0.0.1", "ua")
		if err != nil {
			t.Fatal(err)
		}
		resp, err := doAuthRequest(t, app, "Bearer "+visitorToken)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 401 {
			t.Fatalf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("valid token passes and sets uid local", func(t *testing.T) {
		resp, err := doAuthRequest(t, app, "Bearer "+validToken)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "admin-1") {
			t.Fatalf("expected uid in response body, got %s", string(body))
		}
	})
}
