# Form Object for Complex Input

Replace a controller orchestrating multiple models with a typed input and an explicit registration operation.

## Contents

- Before (fat controller)
- After (form object)

## Before

```go
// BAD: HTTP parsing, validation, team assignment, and persistence are mixed.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    team := r.FormValue("department")
    role := "member"
    if strings.HasSuffix(email, "@company.com") { role = "employee" }
    if err := h.users.CreateWithProfile(r.Context(), email, team, role); err != nil {
        http.Error(w, "registration failed", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

## After

```go
// Transport DTO: parse and validate external input at the HTTP boundary.
type RegistrationInput struct {
    Email string `json:"email"`
    Password string `json:"password"`
    Name string `json:"name"`
    Department string `json:"department"`
}
func (command RegistrationCommand) Validate() error {
    address, err := mail.ParseAddress(input.Email)
    if err != nil || address.Address != input.Email { return ErrEmail }
    if len(input.Password) < 12 || len(input.Password) > 72 { return ErrPassword }
    if strings.TrimSpace(input.Name) == "" { return ErrName }
    return nil
}
// Registration.Register owns user/profile/team changes in one transaction.
// Validate the DTO once; keep team/role business rules in Registration, not the DTO.
command := accounts.RegistrationCommand{
    Email: input.Email, Password: input.Password, Name: input.Name, Department: input.Department,
}
user, err := registration.Register(ctx, command)
if err != nil { return fmt.Errorf("register user: %w", err) }
return user, nil
```

`RegistrationCommand` is owned by the application package (`accounts`). The HTTP package maps its request DTO to that command; the application never imports the HTTP package. Transport field validation and domain invariants have distinct owners. Test the mapping and rollback of every related write.
