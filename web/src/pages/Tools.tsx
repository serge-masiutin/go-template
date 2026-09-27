import { Form, Head, usePage, usePoll } from "@inertiajs/react";
import { useEffect } from "react";
import { Layout } from "../components/Layout";
import { Button } from "../components/Button";
import { Field } from "../components/Field";
import { Notice } from "../components/Notice";
import { WorkStatus } from "../components/WorkStatus";
import type { AssistantRun, NoteEmail, User } from "../types";

type Props = {
  user: User;
  mailEnabled: boolean;
  aiEnabled: boolean;
  emails: NoteEmail[];
  runs: AssistantRun[];
};

export default function Tools({
  user,
  mailEnabled,
  aiEnabled,
  emails,
  runs,
}: Props) {
  const { csrfToken } = usePage().props;
  const emailPending = emails.some(
    (email) => email.state === "queued" || email.state === "sending",
  );
  const aiPending = runs.some(
    (run) => run.state === "queued" || run.state === "running",
  );
  const { start, stop } = usePoll(
    2000,
    { only: ["emails", "runs"] },
    { autoStart: false },
  );
  useEffect(() => {
    if (emailPending || aiPending) start();
    else stop();
    return stop;
  }, [emailPending, aiPending, start, stop]);

  return (
    <Layout user={user}>
      <Head title="Note tools" />
      <div className="mb-8 space-y-2">
        <h1 className="text-3xl font-semibold">Do more with your notes.</h1>
        <p className="text-sm text-muted">
          Email a copy or ask a question about your 50 latest notes.
        </p>
      </div>
      <div className="grid gap-10 md:grid-cols-2">
        <section aria-labelledby="email-heading" className="space-y-4">
          <h2 id="email-heading" className="text-xl font-semibold">
            Email my notes
          </h2>
          <p className="text-sm text-muted">
            Send a plain-text copy to {user.email}.
          </p>
          {mailEnabled ? (
            <Form
              action="/tools/email"
              method="post"
              headers={{ "X-CSRF-Token": csrfToken }}
            >
              {({ processing, errors }) => (
                <div className="space-y-3">
                  {errors.email && <Notice>{errors.email}</Notice>}
                  <Button type="submit" disabled={processing || emailPending}>
                    Email my notes
                  </Button>
                </div>
              )}
            </Form>
          ) : (
            <p className="text-sm text-muted">
              Email delivery has not been enabled for this workspace.
            </p>
          )}
          <ul
            aria-label="Recent emails"
            aria-live="polite"
            className="space-y-3"
          >
            {emails.map((email) => (
              <li
                key={email.id}
                className="space-y-2 rounded-lg border border-border bg-panel p-4"
              >
                <WorkStatus state={email.state} />
                {email.state === "failed" && (
                  <Notice>
                    Delivery failed or its outcome is uncertain. Check your
                    inbox before requesting another copy.
                  </Notice>
                )}
              </li>
            ))}
          </ul>
        </section>
        <section aria-labelledby="assistant-heading" className="space-y-4">
          <h2 id="assistant-heading" className="text-xl font-semibold">
            Ask about my notes
          </h2>
          {aiEnabled ? (
            <>
              <p className="text-sm text-muted">
                Your question and notes will be sent to the workspace's AI
                provider. The assistant can read your notes; it cannot change
                them or send messages.
              </p>
              <Form
                action="/tools/assistant"
                method="post"
                headers={{ "X-CSRF-Token": csrfToken }}
                resetOnSuccess
                className="space-y-3"
              >
                {({ processing, errors }) => (
                  <>
                    <Field
                      multiline
                      id="question"
                      name="question"
                      label="Your question"
                      rows={3}
                      required
                      maxLength={500}
                      error={errors.question}
                    />
                    <Button type="submit" disabled={processing || aiPending}>
                      Ask assistant
                    </Button>
                  </>
                )}
              </Form>
            </>
          ) : (
            <p className="text-sm text-muted">
              The AI assistant has not been enabled for this workspace.
            </p>
          )}
          <ul
            aria-label="Recent answers"
            aria-live="polite"
            className="space-y-3"
          >
            {runs.map((run) => (
              <li
                key={run.id}
                className="space-y-3 rounded-lg border border-border bg-panel p-4"
              >
                <p className="text-sm font-medium break-words">
                  {run.question}
                </p>
                <WorkStatus state={run.state} />
                {run.state === "completed" && (
                  <p className="text-sm break-words whitespace-pre-wrap">
                    {run.answer}
                  </p>
                )}
                {run.state === "failed" && (
                  <Notice>
                    The assistant could not finish. You can submit a new
                    request.
                  </Notice>
                )}
              </li>
            ))}
          </ul>
          {runs.length > 0 && (
            <p className="text-xs text-muted">
              AI answers can be mistaken. Check important details against your
              notes.
            </p>
          )}
        </section>
      </div>
    </Layout>
  );
}
