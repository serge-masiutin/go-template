import { Form, Head, usePage } from "@inertiajs/react";
import { Button } from "../components/Button";
import { Notice } from "../components/Notice";

export default function Login() {
  const { csrfToken } = usePage().props;
  return (
    <main className="mx-auto flex min-h-screen max-w-md items-center px-6">
      <Head title="Sign in" />
      <section className="w-full space-y-6 rounded-xl border border-border bg-panel p-8">
        <div className="space-y-2">
          <p className="text-sm text-muted">GO TEMPLATE</p>
          <h1 className="text-2xl font-semibold">Welcome back</h1>
          <p className="text-sm text-muted">Sign in to your workspace.</p>
        </div>
        <Form
          action="/login"
          method="post"
          headers={{ "X-CSRF-Token": csrfToken }}
          className="space-y-4"
        >
          {({ errors, processing }) => (
            <>
              <div className="space-y-2">
                <label htmlFor="email" className="block text-sm">
                  Email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="username"
                  required
                  className="w-full rounded-md border border-border px-3 py-2"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="password" className="block text-sm">
                  Password
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  className="w-full rounded-md border border-border px-3 py-2"
                />
              </div>
              {errors.email && <Notice>{errors.email}</Notice>}
              <Button type="submit" disabled={processing}>
                {processing ? "Signing in…" : "Sign in"}
              </Button>
            </>
          )}
        </Form>
      </section>
    </main>
  );
}
