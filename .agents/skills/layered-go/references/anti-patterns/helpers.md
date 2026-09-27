# Helper Anti-Patterns

Helpers doing work that belongs in templates or components.

## HTML Construction in Helpers

**Problem:** Helpers building HTML programmatically instead of providing logic for templates.

```go
// BAD: presentation structure and unescaped user input assembled as strings.
func UserCard(name string) string {
    return "<section><h2>" + name + "</h2></section>"
}
```

**Issues:**
- HTML structure hidden in string-building code, harder to read and modify
- No template preview, harder to collaborate with designers
- Logic and markup tightly coupled
- Testing requires rendering, not unit testable
- Misses React/Storybook benefits (typed props, stories, composition)

**Signal:** Nested HTML concatenation, presentation branches in string builders, or unsafe HTML escaping.

**Fix:** Extract to a React component with JSX.

```tsx
type UserCardProps = { name: string; role: string };
export function UserCard({ name, role }: UserCardProps) {
  return <section><h2>{name}</h2><p>{role}</p></section>;
}
```

```tsx
<UserCard name={user.name} role={user.roleLabel} />
```

**Rule of thumb:** Repeated nested markup with stable semantics belongs to a component. A line-count threshold alone is not a reason to extract.
