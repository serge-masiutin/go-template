# Advanced Form Patterns

Advanced form handling for Inertia.js + Go + React.

## Table of Contents

- [useForm Hook](#useform-hook)
- [Data Transforms](#data-transforms)
- [Remember State](#remember-state)
- [Nested Data](#nested-data)
- [Multiple Forms on One Page](#multiple-forms-on-one-page)
- [Conditional Submission](#conditional-submission)
- [Reset and Clear Patterns](#reset-and-clear-patterns)
- [Dynamic Fields](#dynamic-fields)
- [Multi-Step Forms](#multi-step-forms)
- [Client-Side Validation](#client-side-validation)

---

## useForm Hook

Use `useForm` when form data needs to live **outside** the `<Form>` element:

- **Multi-step wizards** — single form state shared across step components
- **Live preview** — form data drives a sibling preview panel (e.g. populate product fields, see the store page update in real time)

```tsx
import { useForm } from '@inertiajs/react'

export default function ProductEditor({ product }: { product: Product }) {
  const form = useForm({
    name: product.name,
    price: product.price,
    description: product.description,
  })

  return (
    <div className="grid grid-cols-2 gap-8">
      <form onSubmit={e => { e.preventDefault(); form.put(`/products/${product.id}`, { headers: { "X-CSRF-Token": csrfToken } }) }}>
        <input value={form.data.name} onChange={e => form.setData('name', e.target.value)} />
        {form.errors.name && <span>{form.errors.name}</span>}
        {/* ... more fields */}
        <button disabled={form.processing}>Save</button>
      </form>

      {/* Preview consumes form.data outside the form */}
      <ProductPreview product={form.data} />
    </div>
  )
}
```

**For everything else, use `<Form>`** — it handles the Inertia visit lifecycle; supply the project CSRF header and implement the matching server input contract.

### useForm API

```tsx
const {
  data,              // Current form data
  setData,           // (key, value) or (values) or (callback)
  post,              // (url, options?) — POST request
  put,               // (url, options?) — PUT request
  patch,             // (url, options?) — PATCH request
  delete: destroy,   // (url, options?) — DELETE request
  processing,        // boolean
  errors,            // { field: "message" }
  hasErrors,         // boolean
  progress,          // Upload progress object
  reset,             // (...fields) — reset to initial values
  clearErrors,       // (...fields) — clear specific errors
  isDirty,           // boolean
  transform,         // (callback) — transform data before send
  wasSuccessful,     // boolean
  recentlySuccessful,// boolean (2s window)
} = useForm({ /* initial data */ })
```

## Data Transforms

Transform data before it's sent to the server. Useful for formatting
dates, removing empty fields, or restructuring data.

### With `<Form>` component:
```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }}
  method="post"
  action="/events"
  transform={(data) => ({
    ...data,
    start_date: formatISO(data.start_date),
    tags: data.tags.filter(Boolean),
  })}
>
```

### With `useForm` hook:
```tsx
const form = useForm({ name: '', price: '' })

const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault()
  form.transform((data) => ({
    ...data,
    price: data.price, // Send a decimal string; the server validates currency/precision.
  }))
  form.post('/products', { headers: { "X-CSRF-Token": csrfToken } })
}
```

## Remember State

Persist form data across navigations so users don't lose their input
when they navigate away and come back.

```tsx
// useForm with remember key
const form = useForm('create-user', {
  name: '',
  email: '',
  role: 'member',
})
// Data persists in history state under the key 'create-user'
```

The first argument is a unique key for the form state. When the user
navigates away and back, the form data is restored.

## Nested Data

Go decodes an explicit nested JSON struct. Match field names and types exactly; there is no automatic nested-attributes persistence.

### With `useForm`:
```tsx
const form = useForm({
  user: {
    name: '',
    email: '',
    address: {
      street: '',
      city: '',
      zip: '',
    },
  },
})

// Update nested field
form.setData('user.address.city', 'New York')
```

### With `<Form>` component:
```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/users">
  {({ errors }) => (
    <>
      <input name="user[name]" />
      <input name="user[email]" />
      <input name="user[address][street]" />
      <input name="user[address][city]" />
    </>
  )}
</Form>
```

Go input contract:
```go
type AddressInput struct {
    Street string `json:"street"`
    City string `json:"city"`
    ZIP string `json:"zip"`
}
type UserInput struct {
    Name string `json:"name"`
    Email string `json:"email"`
    Address AddressInput `json:"address"`
}
// Decode the matching top-level user object strictly, validate once, then map
// to an application-owned command. The operation owns related writes.
```

## Multiple Forms on One Page

Each `<Form>` / `useForm` instance owns the errors from its own submission callback.
Keep the starter's flat `ValidationErrors` and omit `errorBag`, including when two forms
share a field name. Gonertia v3.0.0 does not automatically translate `X-Inertia-Error-Bag`
into a nested `errors[bag]` response. Setting `errorBag` only on the client makes the
client read an absent bag and can silently discard validation failures.

For a feature that truly requires named server error bags, implement that explicit
nested response/flash contract and test both forms, the redirected GET, and partial
reloads before enabling the client option. The following works with the starter:

```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/login">
  {({ errors }) => (
    <>
      <input name="email" />
      {errors.email && <span>{errors.email}</span>}
    </>
  )}
</Form>

<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/register">
  {({ errors }) => (
    <>
      <input name="email" />
      {errors.email && <span>{errors.email}</span>}
    </>
  )}
</Form>
```

## Conditional Submission

Confirm before submitting, or conditionally prevent submission.

```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }}
  method="delete"
  action={`/users/${user.id}`}
  onBefore={() => confirm('Delete this user?')}
>
  {({ processing }) => (
    <button type="submit" disabled={processing}>Delete</button>
  )}
</Form>
```

With `useForm`:
```tsx
const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault()
  if (!confirm('Submit?')) return
  form.post('/users', { headers: { "X-CSRF-Token": csrfToken } })
}
```

## Reset and Clear Patterns

```tsx
const form = useForm({ name: '', email: '', role: 'member' })

// Reset all fields to initial values
form.reset()

// Reset specific fields
form.reset('name', 'email')

// Clear all errors
form.clearErrors()

// Clear specific field errors
form.clearErrors('name', 'email')
```

With `<Form>`:
```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/users">
  {({ reset, clearErrors, wasSuccessful }) => (
    <>
      {/* fields */}
      <button type="button" onClick={() => reset()}>Reset</button>
      <button type="button" onClick={() => clearErrors()}>Clear Errors</button>
    </>
  )}
</Form>
```

## Dynamic Fields

Use `useForm` when dynamic field values need explicit state; a Form with correctly named dynamic inputs is also valid. Use stable row keys when rows can be removed/reordered.

```tsx
const form = useForm({
  items: [{ name: '', quantity: 1 }],
})

const addItem = () => {
  form.setData('items', [...form.data.items, { name: '', quantity: 1 }])
}

const removeItem = (index: number) => {
  form.setData('items', form.data.items.filter((_, i) => i !== index))
}

return (
  <form onSubmit={e => { e.preventDefault(); form.post('/orders', { headers: { "X-CSRF-Token": csrfToken } }) }}>
    {form.data.items.map((item, index) => (
      <div key={index}>
        <input
          value={item.name}
          onChange={e => {
            const items = [...form.data.items]
            items[index] = { ...items[index], name: e.target.value }
            form.setData('items', items)
          }}
        />
        <button type="button" onClick={() => removeItem(index)}>Remove</button>
      </div>
    ))}
    <button type="button" onClick={addItem}>Add Item</button>
    <button type="submit">Submit</button>
  </form>
)
```

## Multi-Step Forms

Use `useForm` (not `<Form>`) with a remember key. Track the current step
with `useState`. Validate per step with `setError`/`clearErrors`, submit
once at the final step.

```tsx
const [step, setStep] = useState(1)

// Remember key preserves progress across navigation
const form = useForm('onboarding', {
  email: '', password: '',  // step 1
  name: '', company: '',    // step 2
  plan: 'starter' as const, // step 3
}).dontRemember('password')

function nextStep() {
  form.clearErrors()
  // Per-step client-side validation
  if (step === 1 && !form.data.email) {
    return form.setError('email', 'Required')
  }
  setStep(s => s + 1)
}

// Single POST at the end — the Go operation validates the complete command
const submit = (e: React.FormEvent) => {
  e.preventDefault()
  form.post('/onboarding', { headers: { "X-CSRF-Token": csrfToken } })
}
```

Each step renders its own fields gated by `{step === N && (...)}`.
Back/forward buttons call `setStep`. Only the final step has
`<button type="submit">`.

## Cache Tag Invalidation

After a form submission, invalidate prefetch caches for pages that show the
affected data. Use `invalidateCacheTags` on `<Form>` or `useForm` submit options:

```tsx
// With <Form>
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/users" invalidateCacheTags={['users', 'dashboard']}>
  {({ errors, processing }) => (
    <>{/* fields */}</>
  )}
</Form>

// With useForm
form.post('/users', {
  headers: { "X-CSRF-Token": csrfToken },
  invalidateCacheTags: ['users', 'dashboard'],
})
```

Tags must match `cacheTags` set on `<Link prefetch cacheTags="...">` — see
`inertia-go-pages` navigation reference for prefetch tag setup.

## Client-Side Validation

Use `setError` for immediate feedback before submitting to the server.
Call `clearErrors` before re-validating to avoid stale messages.

```tsx
function validateAndSubmit(e: React.FormEvent) {
  e.preventDefault()
  form.clearErrors()

  if (form.data.password.length < 12) {
    form.setError('password', 'Must be at least 12 characters')
    return
  }
  if (form.data.password !== form.data.password_confirmation) {
    form.setError('password_confirmation', 'Passwords do not match')
    return
  }

  form.post('/register', { headers: { "X-CSRF-Token": csrfToken } })
}
```

Examples above use `const { csrfToken } = usePage().props` inside the component. Include that binding/import where a snippet is used; never hardcode a token.
