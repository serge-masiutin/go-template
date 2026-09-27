import { afterEach, expect, test } from "vitest";
import { cleanup, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Field } from "./Field";
afterEach(cleanup);

test("associates the label, hint, external description, and current error", () => {
  const { rerender } = render(
    <>
      <p id="privacy">Kept private.</p>
      <Field
        id="email"
        name="email"
        label="Email"
        type="email"
        required
        autoComplete="username"
        hint="Use your work address."
        error="Check this address."
        aria-describedby="privacy"
      />
    </>,
  );
  const field = screen.getByLabelText("Email") as HTMLInputElement;
  expect(field.type).toBe("email");
  expect(field.required).toBe(true);
  expect(field.autocomplete).toBe("username");
  expect(field.getAttribute("aria-describedby")).toBe(
    "privacy email-hint email-error",
  );
  expect(field.getAttribute("aria-invalid")).toBe("true");
  expect(screen.getByRole("alert").id).toBe("email-error");
  rerender(
    <Field
      id="email"
      name="email"
      label="Email"
      hint="Use your work address."
    />,
  );
  const recovered = screen.getByLabelText("Email");
  expect(recovered.getAttribute("aria-describedby")).toBe("email-hint");
  expect(recovered.hasAttribute("aria-invalid")).toBe(false);
  expect(screen.queryByRole("alert")).toBeNull();
});

test("preserves native textarea editing and form submission values", async () => {
  render(
    <form aria-label="Note">
      <Field
        multiline
        id="note"
        name="body"
        label="New note"
        rows={4}
        required
        maxLength={2000}
      />
    </form>,
  );
  const field = screen.getByLabelText("New note") as HTMLTextAreaElement;
  await userEvent.type(field, "First line{Enter}Second line");
  expect(field.rows).toBe(4);
  expect(field.maxLength).toBe(2000);
  expect(
    new FormData(screen.getByRole("form") as HTMLFormElement).get("body"),
  ).toBe("First line\nSecond line");
  expect(field.hasAttribute("aria-describedby")).toBe(false);
});

test("preserves password, read-only and disabled native semantics", async () => {
  render(
    <form aria-label="Credentials">
      <Field
        id="password"
        name="password"
        type="password"
        label="Password"
        autoComplete="current-password"
        defaultValue="existing-value"
        readOnly
      />
      <Field
        id="locked"
        name="locked"
        label="Locked"
        defaultValue="excluded"
        disabled
      />
    </form>,
  );
  const password = screen.getByLabelText("Password") as HTMLInputElement;
  await userEvent.type(password, "ignored");
  expect(password.type).toBe("password");
  expect(password.value).toBe("existing-value");
  expect(password.autocomplete).toBe("current-password");
  const values = new FormData(screen.getByRole("form") as HTMLFormElement);
  expect(values.get("password")).toBe("existing-value");
  expect(values.has("locked")).toBe(false);
});
