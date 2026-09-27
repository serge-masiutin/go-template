package httpapp

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	inertia "github.com/romsar/gonertia/v3"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/notes"
)

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	items, err := a.notes.List(r.Context(), user.ID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	a.render(w, r, "Home", inertia.Props{"user": user, "notes": items})
}

func (a *App) loginPage(w http.ResponseWriter, r *http.Request) { a.render(w, r, "Login", nil) }

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	if !a.loginLimit.Allow(r.RemoteAddr, time.Now()) {
		w.Header().Set("Retry-After", "60")
		http.Error(w, "too many sign-in attempts", http.StatusTooManyRequests)
		return
	}
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &input); err != nil {
		http.Error(w, "invalid request body", 400)
		return
	}
	email, password := input.Email, input.Password
	if len(email) > accounts.MaxEmailBytes || len(password) > 72 {
		a.invalid(w, r, "/login", "email", "Invalid email or password.")
		return
	}
	user, err := a.accounts.Authenticate(r.Context(), email, password)
	if errors.Is(err, accounts.ErrCredentials) {
		a.invalid(w, r, "/login", "email", "Invalid email or password.")
		return
	}
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if err := a.sessions.RenewToken(r.Context()); err != nil {
		a.fail(w, r, err)
		return
	}
	a.sessions.Put(r.Context(), "userID", user.ID)
	a.sessions.Remove(r.Context(), "csrf")
	a.inertia.Redirect(w, r.WithContext(inertia.ClearHistory(r.Context())), "/", http.StatusSeeOther)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	if err := a.sessions.Destroy(r.Context()); err != nil {
		a.fail(w, r, err)
		return
	}
	a.inertia.Redirect(w, r.WithContext(inertia.ClearHistory(r.Context())), "/login", http.StatusSeeOther)
}

func (a *App) createNote(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	var input struct {
		Body string `json:"body"`
	}
	if err := decodeJSON(r, &input); err != nil {
		http.Error(w, "invalid request body", 400)
		return
	}
	err := a.notes.Create(r.Context(), user.ID, input.Body)
	if errors.Is(err, notes.ErrBody) {
		a.invalid(w, r, "/", "body", "Write a note with 1 to 2000 characters.")
		return
	}
	if err != nil {
		a.fail(w, r, err)
		return
	}
	a.inertia.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) deleteNote(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid note ID", 400)
		return
	}
	err = a.notes.Delete(r.Context(), user.ID, id)
	if errors.Is(err, notes.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		a.fail(w, r, err)
		return
	}
	a.inertia.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) admin(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	if !user.Admin {
		http.Error(w, "forbidden", 403)
		return
	}
	count, err := a.accounts.Count(r.Context())
	if err != nil {
		a.fail(w, r, err)
		return
	}
	a.render(w, r, "Admin", inertia.Props{"user": user, "userCount": count})
}

func (a *App) invalid(w http.ResponseWriter, r *http.Request, destination, field, message string) {
	ctx := inertia.SetValidationErrors(r.Context(), inertia.ValidationErrors{field: message})
	a.inertia.Redirect(w, r.WithContext(ctx), destination, http.StatusSeeOther)
}
