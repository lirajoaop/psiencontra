package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/schemas"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// fakeUserRepo is an in-memory UserRepository used by AuthService tests.
// It mimics the behaviour of the real GORM-backed repo for the few cases
// AuthService cares about, including returning gorm.ErrRecordNotFound on
// misses so repository.IsNotFound keeps working.
type fakeUserRepo struct {
	byID    map[uuid.UUID]*schemas.User
	byEmail map[string]*schemas.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byID:    make(map[uuid.UUID]*schemas.User),
		byEmail: make(map[string]*schemas.User),
	}
}

func (r *fakeUserRepo) Create(u *schemas.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if _, exists := r.byEmail[u.Email]; exists {
		return errors.New("duplicate email")
	}
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *fakeUserRepo) FindByID(id uuid.UUID) (*schemas.User, error) {
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeUserRepo) FindByEmail(email string) (*schemas.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeUserRepo) FindByGoogleID(googleID string) (*schemas.User, error) {
	for _, u := range r.byID {
		if u.GoogleID == googleID && googleID != "" {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeUserRepo) Update(u *schemas.User) error {
	if _, exists := r.byID[u.ID]; !exists {
		return gorm.ErrRecordNotFound
	}
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func newTestService() (*AuthService, *fakeUserRepo) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo, "test-secret-do-not-use-in-prod")
	return svc, repo
}

// --- JWT round-trip ---

func TestIssueAndParseToken_RoundTrip(t *testing.T) {
	svc, _ := newTestService()
	id := uuid.New()

	token, err := svc.IssueToken(id)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != id {
		t.Errorf("UserID mismatch: got %s want %s", claims.UserID, id)
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	svc, _ := newTestService()
	token, _ := svc.IssueToken(uuid.New())

	other := NewAuthService(newFakeUserRepo(), "different-secret")
	if _, err := other.ParseToken(token); err == nil {
		t.Fatal("expected error for token signed with different secret")
	}
}

func TestParseToken_Garbage(t *testing.T) {
	svc, _ := newTestService()
	if _, err := svc.ParseToken("not-a-token"); err == nil {
		t.Fatal("expected error for garbage token")
	}
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	svc, repo := newTestService()
	user, token, err := svc.Register("Test@Example.com  ", "supersafe123", "  Joao  ")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email to be normalized, got %q", user.Email)
	}
	if user.Name != "Joao" {
		t.Errorf("expected name to be trimmed, got %q", user.Name)
	}
	if user.PasswordHash == "" {
		t.Fatal("expected password hash to be set")
	}
	if user.PasswordHash == "supersafe123" {
		t.Fatal("password hash must not equal plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("supersafe123")); err != nil {
		t.Errorf("password hash does not validate: %v", err)
	}

	stored, _ := repo.FindByEmail("test@example.com")
	if stored == nil || stored.ID != user.ID {
		t.Error("user not persisted in repo")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _ := newTestService()
	if _, _, err := svc.Register("dup@example.com", "supersafe123", ""); err != nil {
		t.Fatalf("first Register: %v", err)
	}
	_, _, err := svc.Register("dup@example.com", "anothersafe1", "")
	if !errors.Is(err, ErrEmailAlreadyUsed) {
		t.Fatalf("expected ErrEmailAlreadyUsed, got %v", err)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	svc, _ := newTestService()
	_, _, err := svc.Register("short@example.com", "1234567", "")
	if !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestRegister_EmptyEmail(t *testing.T) {
	svc, _ := newTestService()
	_, _, err := svc.Register("   ", "supersafe123", "")
	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	svc, _ := newTestService()
	if _, _, err := svc.Register("login@example.com", "supersafe123", "Joao"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	user, token, err := svc.Login("LOGIN@example.com", "supersafe123")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if user.Email != "login@example.com" {
		t.Errorf("got email %q", user.Email)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _ := newTestService()
	_, _, _ = svc.Register("wrong@example.com", "supersafe123", "")
	_, _, err := svc.Login("wrong@example.com", "incorrect")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	svc, _ := newTestService()
	_, _, err := svc.Login("ghost@example.com", "whatever")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials (no leak), got %v", err)
	}
}

func TestLogin_GoogleOnlyUser(t *testing.T) {
	svc, repo := newTestService()
	// Simulate a user that signed up via Google: no PasswordHash.
	googleUser := &schemas.User{
		Email:    "google@example.com",
		GoogleID: "google-sub-123",
		Name:     "Google User",
	}
	if err := repo.Create(googleUser); err != nil {
		t.Fatalf("seed: %v", err)
	}
	_, _, err := svc.Login("google@example.com", "anything")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for google-only user, got %v", err)
	}
}

// --- UpsertGoogleUser ---

func TestUpsertGoogleUser_NewUser(t *testing.T) {
	svc, repo := newTestService()
	user, token, err := svc.UpsertGoogleUser("sub-1", "new@example.com", "New User", "https://avatar")
	if err != nil {
		t.Fatalf("UpsertGoogleUser: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if user.GoogleID != "sub-1" {
		t.Errorf("google id not set: %q", user.GoogleID)
	}
	if user.AvatarURL != "https://avatar" {
		t.Errorf("avatar not set: %q", user.AvatarURL)
	}
	if _, err := repo.FindByEmail("new@example.com"); err != nil {
		t.Errorf("user not persisted: %v", err)
	}
}

func TestUpsertGoogleUser_ExistingByGoogleID(t *testing.T) {
	svc, _ := newTestService()
	// First call creates.
	first, _, err := svc.UpsertGoogleUser("sub-2", "again@example.com", "First", "")
	if err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	// Second call with same google id returns the same user (no duplicate).
	second, _, err := svc.UpsertGoogleUser("sub-2", "different@example.com", "Other", "")
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	if first.ID != second.ID {
		t.Errorf("expected same user ID, got %s vs %s", first.ID, second.ID)
	}
}

func TestUpsertGoogleUser_LinksExistingEmail(t *testing.T) {
	svc, repo := newTestService()
	// User exists from email/password registration.
	if _, _, err := svc.Register("link@example.com", "supersafe123", "Original"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	// Then they sign in with Google using the same email.
	user, _, err := svc.UpsertGoogleUser("sub-3", "link@example.com", "Google Name", "https://pic")
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if user.GoogleID != "sub-3" {
		t.Errorf("expected google id to be linked, got %q", user.GoogleID)
	}
	// Existing name should be preserved (it was non-empty).
	if user.Name != "Original" {
		t.Errorf("expected existing name to be preserved, got %q", user.Name)
	}
	// Avatar was empty before, should now be filled.
	if user.AvatarURL != "https://pic" {
		t.Errorf("expected avatar to be filled, got %q", user.AvatarURL)
	}
	// And the user in the repo should be the same record (not a duplicate).
	stored, _ := repo.FindByEmail("link@example.com")
	if stored.ID != user.ID {
		t.Error("upsert created a duplicate instead of linking")
	}
}
