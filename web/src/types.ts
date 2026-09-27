export type User = { id: string; email: string; admin: boolean };
export type Note = { id: string; body: string };

export type SharedProps = { csrfToken: string };

declare module "@inertiajs/core" {
  interface InertiaConfig {
    sharedPageProps: SharedProps;
  }
}

export type NoteEmail = {
  id: string;
  state: "queued" | "sending" | "sent" | "failed";
  createdAt: string;
};
export type AssistantRun = {
  id: string;
  question: string;
  state: "queued" | "running" | "completed" | "failed";
  answer: string;
  createdAt: string;
};
