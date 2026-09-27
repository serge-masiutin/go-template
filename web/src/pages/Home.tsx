import { Form, Head, router, usePage } from "@inertiajs/react";
import type { Note, User } from "../types";
import { Layout } from "../components/Layout";
import { Button } from "../components/Button";
import { Notice } from "../components/Notice";

export default function Home({ user, notes }: { user: User; notes: Note[] }) {
  const { csrfToken } = usePage().props;
  return (
    <Layout user={user}>
      <Head title="Workspace" />
      <div className="mb-8 space-y-2">
        <p className="text-sm text-muted">YOUR WORKSPACE</p>
        <h1 className="text-3xl font-semibold">A place to start.</h1>
        <p className="text-sm text-muted">Signed in as {user.email}</p>
      </div>
      <section aria-label="Notes" className="max-w-2xl space-y-6">
        <Form
          action="/notes"
          method="post"
          headers={{ "X-CSRF-Token": csrfToken }}
          resetOnSuccess
          className="space-y-3"
        >
          {({ errors, processing }) => (
            <>
              <label htmlFor="body" className="block text-sm font-medium">
                New note
              </label>
              <textarea
                id="body"
                name="body"
                required
                maxLength={2000}
                rows={4}
                className="w-full rounded-md border border-border bg-panel p-3"
              />
              {errors.body && <Notice>{errors.body}</Notice>}
              <Button type="submit" disabled={processing}>
                Add note
              </Button>
            </>
          )}
        </Form>
        {notes.length === 0 ? (
          <p className="text-sm text-muted">
            Your notes will appear here. Only you can see them.
          </p>
        ) : (
          <ul className="space-y-3">
            {notes.map((note) => (
              <li
                key={note.id}
                className="flex items-start justify-between gap-4 rounded-lg border border-border bg-panel p-4"
              >
                <p className="min-w-0 text-sm break-words whitespace-pre-wrap">
                  {note.body}
                </p>
                <Button
                  variant="secondary"
                  onClick={() =>
                    router.delete(`/notes/${note.id}`, {
                      headers: { "X-CSRF-Token": csrfToken },
                    })
                  }
                >
                  Delete
                </Button>
              </li>
            ))}
          </ul>
        )}
        <p className="text-xs text-muted">Showing up to 50 recent notes.</p>
      </section>
    </Layout>
  );
}
