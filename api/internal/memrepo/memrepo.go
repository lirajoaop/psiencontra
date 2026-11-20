// Package memrepo provides in-memory implementations of the repository
// interfaces used by the service layer. It exists only to support tests:
// it lets service, handler and router tests run without PostgreSQL.
//
// Misses return gorm.ErrRecordNotFound so repository.IsNotFound keeps
// working exactly as it does against the real database.
package memrepo

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/schemas"
	"gorm.io/gorm"
)

// Store holds every table in memory. Use its accessor methods to get the
// per-table repositories, which share the same underlying data.
type Store struct {
	mu        sync.Mutex
	sessions  map[uuid.UUID]*schemas.Session
	responses map[uuid.UUID][]schemas.Response
	results   map[uuid.UUID]*schemas.Result
	users     map[uuid.UUID]*schemas.User
	nextID    uint

	// Error injection: when non-nil, the matching operation fails with it.
	CreateSessionErr error
	CreateResultErr  error
	HistoryErr       error
}

func New() *Store {
	return &Store{
		sessions:  map[uuid.UUID]*schemas.Session{},
		responses: map[uuid.UUID][]schemas.Response{},
		results:   map[uuid.UUID]*schemas.Result{},
		users:     map[uuid.UUID]*schemas.User{},
	}
}

func (s *Store) Sessions() *Sessions   { return &Sessions{s} }
func (s *Store) Responses() *Responses { return &Responses{s} }
func (s *Store) Results() *Results     { return &Results{s} }
func (s *Store) Users() *Users         { return &Users{s} }

// --- sessions ---

type Sessions struct{ s *Store }

func (r *Sessions) Create(sess *schemas.Session) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if r.s.CreateSessionErr != nil {
		return r.s.CreateSessionErr
	}
	if sess.ID == uuid.Nil {
		sess.ID = uuid.New()
	}
	if sess.QuestionnaireType == "" {
		sess.QuestionnaireType = "simple"
	}
	if sess.CreatedAt.IsZero() {
		sess.CreatedAt = time.Now()
	}
	cp := *sess
	r.s.sessions[sess.ID] = &cp
	return nil
}

func (r *Sessions) FindByID(id uuid.UUID) (*schemas.Session, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	sess, ok := r.s.sessions[id]
	if !ok {
		return &schemas.Session{}, gorm.ErrRecordNotFound
	}
	cp := *sess
	return &cp, nil
}

func (r *Sessions) ClaimAnonymous(sessionID, userID uuid.UUID) (int64, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	sess, ok := r.s.sessions[sessionID]
	if !ok || sess.UserID != nil {
		return 0, nil
	}
	uid := userID
	sess.UserID = &uid
	return 1, nil
}

func (r *Sessions) FindCompletedByUserID(userID uuid.UUID) ([]schemas.Session, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if r.s.HistoryErr != nil {
		return nil, r.s.HistoryErr
	}
	var out []schemas.Session
	for _, sess := range r.s.sessions {
		if sess.UserID == nil || *sess.UserID != userID {
			continue
		}
		res, ok := r.s.results[sess.ID]
		if !ok {
			continue
		}
		cp := *sess
		resCopy := *res
		cp.Result = &resCopy
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

// --- responses ---

type Responses struct{ s *Store }

func (r *Responses) CreateBatch(responses []schemas.Response) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, resp := range responses {
		r.s.nextID++
		resp.ID = r.s.nextID
		r.s.responses[resp.SessionID] = append(r.s.responses[resp.SessionID], resp)
	}
	return nil
}

func (r *Responses) FindBySessionID(sessionID uuid.UUID) ([]schemas.Response, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	out := append([]schemas.Response(nil), r.s.responses[sessionID]...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].QuestionID < out[j].QuestionID })
	return out, nil
}

// --- results ---

type Results struct{ s *Store }

func (r *Results) Create(res *schemas.Result) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if r.s.CreateResultErr != nil {
		return r.s.CreateResultErr
	}
	r.s.nextID++
	res.ID = r.s.nextID
	cp := *res
	r.s.results[res.SessionID] = &cp
	return nil
}

func (r *Results) FindBySessionID(sessionID uuid.UUID) (*schemas.Result, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	res, ok := r.s.results[sessionID]
	if !ok {
		return &schemas.Result{}, gorm.ErrRecordNotFound
	}
	cp := *res
	return &cp, nil
}

// --- users ---

type Users struct{ s *Store }

func (r *Users) Create(u *schemas.User) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, existing := range r.s.users {
		if existing.Email == u.Email {
			return gorm.ErrDuplicatedKey
		}
	}
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	cp := *u
	r.s.users[u.ID] = &cp
	return nil
}

func (r *Users) FindByID(id uuid.UUID) (*schemas.User, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if u, ok := r.s.users[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *Users) FindByEmail(email string) (*schemas.User, error) {
	return r.findBy(func(u *schemas.User) bool { return u.Email == email })
}

func (r *Users) FindByGoogleID(googleID string) (*schemas.User, error) {
	return r.findBy(func(u *schemas.User) bool { return googleID != "" && u.GoogleID == googleID })
}

func (r *Users) Update(u *schemas.User) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if _, ok := r.s.users[u.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	cp := *u
	r.s.users[u.ID] = &cp
	return nil
}

func (r *Users) findBy(match func(*schemas.User) bool) (*schemas.User, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, u := range r.s.users {
		if match(u) {
			cp := *u
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
