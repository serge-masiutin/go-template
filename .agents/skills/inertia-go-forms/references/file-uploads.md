# File Uploads

File upload handling for Inertia.js + Go with an explicit bounded multipart backend. Upload routes/storage are not installed in this starter; implement and test that boundary before using these examples.

## Table of Contents

- [Basic Upload with Form Component](#basic-upload-with-form-component)
- [Upload with useForm](#upload-with-useform)
- [Progress Tracking](#progress-tracking)
- [Go Multipart Backend](#go-multipart-backend)
- [Multiple Files](#multiple-files)
- [Image Preview](#image-preview)
- [Direct Uploads](#direct-uploads)

---

## Basic Upload with Form Component

Inertia auto-detects file inputs and switches to `multipart/form-data`.

```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/avatars">
  {({ progress, processing }) => (
    <>
      <input type="file" name="avatar" accept="image/*" />
      {progress && (
        <progress value={progress.percentage} max="100">
          {progress.percentage}%
        </progress>
      )}
      <button type="submit" disabled={processing}>Upload</button>
    </>
  )}
</Form>
```

No manual `FormData` construction needed. Include the project X-CSRF-Token header; let the browser set the multipart content type and boundary.

## Upload with useForm

```tsx
const form = useForm<{ avatar: File | null }>({
  avatar: null,
})

const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault()
  form.post('/avatars', {
  headers: { "X-CSRF-Token": csrfToken },
    forceFormData: true, // ensure FormData even if no file selected
  })
}

return (
  <form onSubmit={handleSubmit}>
    <input
      type="file"
      onChange={e => form.setData('avatar', e.target.files?.[0] ?? null)}
    />
    {form.progress && (
      <progress value={form.progress.percentage} max="100" />
    )}
    <button disabled={form.processing}>Upload</button>
  </form>
)
```

## Progress Tracking

Both `<Form>` and `useForm` provide upload progress automatically.

```tsx
// Form component
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/documents" onProgress={(progress) => {
  console.log(`${progress.percentage}% uploaded`)
}}>
  {({ progress }) => (
    progress && <ProgressBar value={progress.percentage} />
  )}
</Form>
```

The `progress` object:
```typescript
{
  percentage: number  // 0-100
  total: number       // total bytes
  loaded: number      // bytes uploaded
}
```

## Go Multipart Backend

The starter caps JSON bodies at 16 KiB. Add a separately bounded upload route/middleware; simply calling ParseMultipartForm behind the current cap does not enable larger uploads.

```go
r.Body = http.MaxBytesReader(w, r.Body, 8<<20)
if err := r.ParseMultipartForm(1 << 20); err != nil {
    http.Error(w, "invalid or oversized upload", http.StatusBadRequest)
    return
}
defer r.MultipartForm.RemoveAll()
file, header, err := r.FormFile("avatar")
if err != nil { http.Error(w, "avatar is required", 400); return }
defer file.Close()
```

Authenticate and check CSRF before parsing. The storage operation must bound size, validate actual content, generate its own object name, ignore client path components and enforce ownership. `header.Filename` and client MIME type are untrusted. Return success only after durable storage and metadata persistence succeed; define cleanup for partial failure.

## Multiple Files

```tsx
<Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/documents">
  {({ progress }) => (
    <>
      <input type="file" name="documents[]" multiple />
      {progress && <progress value={progress.percentage} max="100" />}
      <button type="submit">Upload All</button>
    </>
  )}
</Form>
```

Go multipart field names are literal. Inertia's array conversion may emit indexed keys such as `documents[0]`; implement that exact allowlisted naming contract, with a maximum file count and per-file/total limits. Do not assume Rails bracket parsing. Test the real browser payload against the handler.

## Image Preview

`useState` for preview is local UI state — still use `<Form>`, not `useForm`:

```tsx
function AvatarUpload() {
  const [preview, setPreview] = useState<string | null>(null)
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview) }, [preview])

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setPreview(URL.createObjectURL(file))
    }
  }

  return (
    <Form headers={{ "X-CSRF-Token": csrfToken }} method="post" action="/avatars">
      {({ progress, processing }) => (
        <>
          {preview && <img src={preview} alt="Preview" className="w-24 h-24 rounded-full" />}
          <input type="file" name="avatar" accept="image/*" onChange={handleFileChange} />
          {progress && <progress value={progress.percentage} max="100" />}
          <button type="submit" disabled={processing}>Upload</button>
        </>
      )}
    </Form>
  )
}
```

## Direct Uploads

For large files, a backend can issue a short-lived upload URL for an authorized object key. The browser uploads to storage, then submits the returned upload identifier through an Inertia mutation. The server must verify ownership, expected size/type, completion and single-use semantics before attaching it.

Do not port Active Storage signed IDs or its JavaScript package. Select the actual storage provider first and use its documented signing/verification API. A client-provided public URL is not proof of a completed authorized upload.

Examples above use `const { csrfToken } = usePage().props` inside the component. Include that binding/import where a snippet is used; never hardcode a token.
