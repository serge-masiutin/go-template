import { Head } from "@inertiajs/react";
import type { User } from "../types";
import { Layout } from "../components/Layout";

export default function Admin({
  user,
  userCount,
}: {
  user: User;
  userCount: number;
}) {
  return (
    <Layout user={user}>
      <Head title="Admin" />
      <h1 className="mb-6 text-3xl font-semibold">Administration</h1>
      <div className="max-w-sm rounded-lg border border-border bg-panel p-6">
        <p className="text-sm text-muted">Accounts</p>
        <p className="mt-2 text-4xl">{userCount}</p>
      </div>
    </Layout>
  );
}
