import { Form, Link, usePage } from "@inertiajs/react";
import type { ReactNode } from "react";
import type { User } from "../types";
import { Button } from "./Button";

export function Layout({
  user,
  children,
}: {
  user: User;
  children: ReactNode;
}) {
  const { csrfToken } = usePage().props;
  return (
    <div className="mx-auto max-w-5xl px-6">
      <header className="flex flex-wrap items-center justify-between gap-4 border-b border-border py-6">
        <Link href="/" className="font-semibold">
          Go Template
        </Link>
        <nav aria-label="Main" className="flex items-center gap-4 text-sm">
          {user.admin && <Link href="/admin">Admin</Link>}
          <Form
            action="/logout"
            method="post"
            headers={{ "X-CSRF-Token": csrfToken }}
          >
            <Button type="submit" variant="secondary">
              Sign out
            </Button>
          </Form>
        </nav>
      </header>
      <main className="py-10">{children}</main>
    </div>
  );
}
