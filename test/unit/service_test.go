package unit

import (
    "context"
    "testing"
    "time"
    "go-backend/database"
    "go-backend/internal/models"
    "go-backend/internal/repositories"
    "go-backend/internal/services"
)

func TestService_CreateGetList(t *testing.T) {
    path := t.TempDir() + "/unit1.db"
    db, err := database.Connect(path)
    if err != nil { t.Fatal(err) }
    repo := repositories.NewDeviceRepository(db)
    svc := services.NewDeviceService(repo)
    id1, err := svc.Create(context.Background(), &models.Device{Name: "Phone X", Brand: "Acme", State: models.StateAvailable})
    if err != nil { t.Fatal(err) }
    id2, err := svc.Create(context.Background(), &models.Device{Name: "Phone Y", Brand: "Acme", State: models.StateInactive})
    if err != nil { t.Fatal(err) }
    d, err := svc.Get(context.Background(), id1)
    if err != nil { t.Fatal(err) }
    if d.ID != id1 || d.Name != "Phone X" || d.Brand != "Acme" || d.State != models.StateAvailable { t.Fatalf("unexpected device: %+v", d) }
    list, err := svc.List(context.Background(), "Acme", "")
    if err != nil { t.Fatal(err) }
    if len(list) != 2 { t.Fatalf("expected 2, got %d", len(list)) }
    list, err = svc.List(context.Background(), "", string(models.StateInactive))
    if err != nil { t.Fatal(err) }
    if len(list) != 1 || list[0].ID != id2 { t.Fatalf("filter by state failed: %+v", list) }
}

func TestService_UpdatePatchDeleteRules(t *testing.T) {
    path := t.TempDir() + "/unit2.db"
    db, err := database.Connect(path)
    if err != nil { t.Fatal(err) }
    svc := services.NewDeviceService(repositories.NewDeviceRepository(db))
    id, err := svc.Create(context.Background(), &models.Device{Name: "A", Brand: "B", State: models.StateInUse})
    if err != nil { t.Fatal(err) }
    orig, _ := svc.Get(context.Background(), id)
    err = svc.Update(context.Background(), id, &models.Device{Name: "A2", Brand: "B", State: models.StateInUse, CreatedAt: orig.CreatedAt})
    if err != models.ErrCannotUpdateFields { t.Fatalf("expected ErrCannotUpdateFields, got %v", err) }
    err = svc.Update(context.Background(), id, &models.Device{Name: "A", Brand: "B", State: models.StateInUse, CreatedAt: models.NewFormattedTime(orig.CreatedAt.Time.Add(time.Minute))})
    if err != models.ErrCannotUpdateCreated { t.Fatalf("expected ErrCannotUpdateCreated, got %v", err) }
    err = svc.Update(context.Background(), id, &models.Device{Name: "A", Brand: "B", State: models.StateInactive, CreatedAt: orig.CreatedAt})
    if err != nil { t.Fatalf("unexpected err: %v", err) }
    if err := svc.Patch(context.Background(), id, map[string]any{"created_at": time.Now()}); err != models.ErrCannotUpdateCreated { t.Fatalf("expected ErrCannotUpdateCreated, got %v", err) }
    if err := svc.Delete(context.Background(), id); err != nil { t.Fatalf("unexpected err: %v", err) }
}

func TestService_DeleteInUseBlocked(t *testing.T) {
    path := t.TempDir() + "/unit3.db"
    db, err := database.Connect(path)
    if err != nil { t.Fatal(err) }
    svc := services.NewDeviceService(repositories.NewDeviceRepository(db))
    id, err := svc.Create(context.Background(), &models.Device{Name: "A", Brand: "B", State: models.StateInUse})
    if err != nil { t.Fatal(err) }
    if err := svc.Delete(context.Background(), id); err != models.ErrCannotDeleteInUse { t.Fatalf("expected ErrCannotDeleteInUse, got %v", err) }
}

func TestService_PatchInvalidState(t *testing.T) {
    path := t.TempDir() + "/unit4.db"
    db, err := database.Connect(path)
    if err != nil { t.Fatal(err) }
    svc := services.NewDeviceService(repositories.NewDeviceRepository(db))
    id, err := svc.Create(context.Background(), &models.Device{Name: "C", Brand: "D", State: models.StateAvailable})
    if err != nil { t.Fatal(err) }
    if err := svc.Patch(context.Background(), id, map[string]any{"state": "bad"}); err != models.ErrInvalidState { t.Fatalf("expected ErrInvalidState, got %v", err) }
    if err := svc.Patch(context.Background(), id, map[string]any{"state": 123}); err == nil || err.Error() != "invalid state type" { t.Fatalf("expected invalid state type, got %v", err) }
}

func TestService_UpdateInvalidState(t *testing.T) {
    path := t.TempDir() + "/unit5.db"
    db, err := database.Connect(path)
    if err != nil { t.Fatal(err) }
    svc := services.NewDeviceService(repositories.NewDeviceRepository(db))
    id, err := svc.Create(context.Background(), &models.Device{Name: "E", Brand: "F", State: models.StateAvailable})
    if err != nil { t.Fatal(err) }
    existing, _ := svc.Get(context.Background(), id)
    if err := svc.Update(context.Background(), id, &models.Device{Name: "E", Brand: "F", State: models.State("bad"), CreatedAt: existing.CreatedAt}); err != models.ErrInvalidState { t.Fatalf("expected ErrInvalidState, got %v", err) }
}

func TestModel_ValidateNew(t *testing.T) {
    d := models.Device{Name: "A", Brand: "B", State: models.StateAvailable}
    if err := d.ValidateNew(); err != nil { t.Fatalf("unexpected err: %v", err) }
    if d.CreatedAt.IsZero() { t.Fatalf("expected CreatedAt to be set") }
    d = models.Device{Name: "A", Brand: "B", State: models.State("bad")}
    if err := d.ValidateNew(); err != models.ErrInvalidState { t.Fatalf("expected ErrInvalidState, got %v", err) }
    d = models.Device{Name: "", Brand: "B", State: models.StateAvailable}
    if err := d.ValidateNew(); err == nil { t.Fatalf("expected error for missing name") }
    d = models.Device{Name: "A", Brand: "", State: models.StateAvailable}
    if err := d.ValidateNew(); err == nil { t.Fatalf("expected error for missing brand") }
}

func TestService_ListCombinedFilter(t *testing.T) {
    path := t.TempDir() + "/unit6.db"
    db, err := database.Connect(path)
    if err != nil { t.Fatal(err) }
    svc := services.NewDeviceService(repositories.NewDeviceRepository(db))
    _, _ = svc.Create(context.Background(), &models.Device{Name: "P1", Brand: "Acme", State: models.StateAvailable})
    _, _ = svc.Create(context.Background(), &models.Device{Name: "P2", Brand: "Acme", State: models.StateInactive})
    list, err := svc.List(context.Background(), "Acme", string(models.StateAvailable))
    if err != nil { t.Fatal(err) }
    if len(list) != 1 || list[0].Brand != "Acme" || list[0].State != models.StateAvailable { t.Fatalf("unexpected list: %+v", list) }
}
