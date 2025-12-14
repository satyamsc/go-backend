package integration

import (
	"bytes"
	"encoding/json"
	"go-backend/database"
	"go-backend/internal/models"
	"go-backend/internal/routers"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlers_CRUD(t *testing.T) {
	path := t.TempDir() + "/http.db"
	db, err := database.Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = db.AutoMigrate(&models.Device{})
	r := routers.New(db)
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(`{"name":"X","brand":"Acme","state":"available"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodGet, "/devices", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/devices/1", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodPatch, "/devices/1", bytes.NewBufferString(`{"state":"inactive"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodDelete, "/devices/1", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlers_ErrorCases(t *testing.T) {
	path := t.TempDir() + "/http2.db"
	db, err := database.Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = db.AutoMigrate(&models.Device{})
	r := routers.New(db)
	req := httptest.NewRequest(http.MethodGet, "/devices/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(`{"name":"A","brand":"B","state":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	var payload struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &payload)
	if payload.Code != "validation_error" {
		t.Fatalf("expected code validation_error, got %s", payload.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(`{"name":"X","brand":"Acme","state":"in-use"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodDelete, "/devices/1", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &payload)
	if payload.Code != "in_use_delete_blocked" {
		t.Fatalf("expected code in_use_delete_blocked, got %s", payload.Code)
	}
	req = httptest.NewRequest(http.MethodPut, "/devices/1", bytes.NewBufferString(`{"name":"X","brand":"Acme","state":"in-use","created_at":"2025-12-14T20:01:13Z"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &payload)
	if payload.Code != "cannot_update_created_at" {
		t.Fatalf("expected code cannot_update_created_at, got %s", payload.Code)
	}
}

func TestDocsAndHealthz(t *testing.T) {
	path := t.TempDir() + "/http3.db"
	db, err := database.Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	r := routers.New(db)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "SwaggerUIBundle") {
		t.Fatalf("docs html missing bundle")
	}
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/openapi.yaml", []byte("openapi: 3.0.3"), 0644)
	_ = os.Chdir(dir)
	req = httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "openapi:") {
		t.Fatalf("openapi yaml missing header")
	}
}
