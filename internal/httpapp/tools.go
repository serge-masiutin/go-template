package httpapp

import (
	"errors"
	"net/http"

	inertia "github.com/romsar/gonertia/v3"
	"github.com/serge-masiutin/go-template/internal/assistant"
	"github.com/serge-masiutin/go-template/internal/notemail"
)

func (a *App) toolsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	emails, err := a.emails.Recent(r.Context(), user.ID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	runs, err := a.assistant.Recent(r.Context(), user.ID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	a.render(w, r, "Tools", inertia.Props{"user": user, "mailEnabled": a.cfg.Mail.Enabled, "aiEnabled": a.cfg.AI.Enabled, "emails": emails, "runs": runs})
}
func (a *App) emailNotes(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	if !a.cfg.Mail.Enabled {
		http.Error(w, "mail is not configured", http.StatusServiceUnavailable)
		return
	}
	var input struct{}
	if err := decodeJSON(r, &input); err != nil {
		http.Error(w, "invalid request body", 400)
		return
	}
	err := a.emails.Request(r.Context(), user.ID)
	if errors.Is(err, notemail.ErrPending) {
		a.invalid(w, r, "/tools", "email", "An email is already pending.")
		return
	}
	if err != nil {
		a.fail(w, r, err)
		return
	}
	a.inertia.Redirect(w, r, "/tools", http.StatusSeeOther)
}
func (a *App) askAssistant(w http.ResponseWriter, r *http.Request) {
	user, ok := a.user(w, r)
	if !ok {
		return
	}
	if !a.cfg.AI.Enabled {
		http.Error(w, "AI is not configured", http.StatusServiceUnavailable)
		return
	}
	var input struct {
		Question string `json:"question"`
	}
	if err := decodeJSON(r, &input); err != nil {
		http.Error(w, "invalid request body", 400)
		return
	}
	err := a.assistant.Request(r.Context(), user.ID, input.Question, a.cfg.AI.Model)
	switch {
	case errors.Is(err, assistant.ErrQuestion):
		a.invalid(w, r, "/tools", "question", "Ask a question with 1 to 500 characters.")
	case errors.Is(err, assistant.ErrPending):
		a.invalid(w, r, "/tools", "question", "An answer is already pending.")
	case err != nil:
		a.fail(w, r, err)
	default:
		a.inertia.Redirect(w, r, "/tools", http.StatusSeeOther)
	}
}
