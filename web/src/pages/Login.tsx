import { Form, Head, usePage } from "@inertiajs/react";
import { Button } from "../components/Button";
import { Field } from "../components/Field";

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
              <Field
                id="email"
                name="email"
                label="Email"
                type="email"
                autoComplete="username"
                required
                error={errors.email}
              />
              <Field
                id="password"
                name="password"
                label="Password"
                type="password"
                autoComplete="current-password"
                required
              />
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
