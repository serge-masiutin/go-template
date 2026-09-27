# Flash Toast with Sonner

Preserve the upstream React hook and persistent-layout integration. Sonner/shadcn are optional additions; the starter does not install a toast library or emit notice/alert messages by default.

## Go Setup

Implement the session-backed Gonertia FlashProvider and declare the exact `flashDataType` in TypeScript before using these keys.

```go
ctx := inertia.SetFlash(r.Context(), inertia.Flash{"notice": "Profile saved."})
engine.Redirect(w, r.WithContext(ctx), "/profile", http.StatusSeeOther)
```

Use SetValidationErrors for field validation, not a generic toast. Keep notice/alert a deliberate payload contract; no Rails initializer or automatic flash-key exposure exists.

## useFlash Hook

**Important:** `toast` is imported from `'sonner'` (the package), NOT from `@/components/ui/sonner` (that only exports `Toaster`).

```tsx
// web/src/hooks/use-flash.ts
import { router, usePage } from '@inertiajs/react'
import { useEffect, useRef } from 'react'
import { toast } from 'sonner' // NOT from '@/components/ui/sonner'

function showFlash(flash: FlashData) {
  if (flash.alert) toast.error(flash.alert)
  if (flash.notice) toast(flash.notice)
}

export function useFlash() {
  const { flash } = usePage()
  const toastShowed = useRef(false)

  // Show flash from initial page load
  useEffect(() => {
    if (!toastShowed.current) {
      toastShowed.current = true
      showFlash(flash)
    }
  }, [flash])

  // Listen for flash events (client-side flash, redirects)
  useEffect(() => {
    return router.on('flash', (event) => {
      showFlash(event.detail.flash)
    })
  }, [])
}
```

## Layout Integration

Use in persistent layout (runs once, covers all pages):

```tsx
// web/src/layouts/persistent-layout.tsx
import { Toaster } from 'sonner'
import { useFlash } from '@/hooks/use-flash'

export function PersistentLayout({ children }) {
  useFlash()
  return <>{children}<Toaster /></>
}
```
