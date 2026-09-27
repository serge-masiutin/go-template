package httpapp

import (
	"context"
	"encoding/json"

	"github.com/alexedwards/scs/v2"
	inertia "github.com/romsar/gonertia/v3"
)

type flashProvider struct{ sessions *scs.SessionManager }

func (p flashProvider) FlashErrors(ctx context.Context, values inertia.ValidationErrors) error {
	encoded, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	p.sessions.Put(ctx, "errors", string(encoded))
	return nil
}
func (p flashProvider) GetErrors(ctx context.Context) (inertia.ValidationErrors, error) {
	encoded := p.sessions.PopString(ctx, "errors")
	if encoded == "" {
		return nil, nil
	}
	var values inertia.ValidationErrors
	err := json.Unmarshal([]byte(encoded), &values)
	if err != nil {
		panic(err)
	}
	return values, nil
}
func (p flashProvider) Flash(ctx context.Context, values inertia.Flash) error {
	encoded, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	p.sessions.Put(ctx, "flash", string(encoded))
	return nil
}
func (p flashProvider) GetFlash(ctx context.Context) (inertia.Flash, error) {
	encoded := p.sessions.PopString(ctx, "flash")
	if encoded == "" {
		return nil, nil
	}
	var values inertia.Flash
	err := json.Unmarshal([]byte(encoded), &values)
	if err != nil {
		panic(err)
	}
	return values, nil
}
func (p flashProvider) FlashClearHistory(ctx context.Context) error {
	p.sessions.Put(ctx, "clearHistory", true)
	return nil
}
func (p flashProvider) ShouldClearHistory(ctx context.Context) (bool, error) {
	return p.sessions.PopBool(ctx, "clearHistory"), nil
}
