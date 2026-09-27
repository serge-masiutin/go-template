export type User = { id: string; email: string; admin: boolean };
export type Note = { id: string; body: string };

export type SharedProps = { csrfToken: string };

declare module "@inertiajs/core" {
  interface InertiaConfig {
    sharedPageProps: SharedProps;
  }
}
