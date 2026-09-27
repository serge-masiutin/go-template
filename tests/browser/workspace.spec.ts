import { test, expect } from "@playwright/test";

test("sign in, validate, create and delete a private note, sign out", async ({
  page,
}) => {
  const email = process.env.BROWSER_TEST_EMAIL;
  const password = process.env.BROWSER_TEST_PASSWORD;
  if (!email || !password)
    throw new Error(
      "BROWSER_TEST_EMAIL and BROWSER_TEST_PASSWORD are required",
    );
  const failures: string[] = [];
  page.on("pageerror", (error) => failures.push(error.message));
  await page.goto("/");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill("incorrect-password");
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(page.getByRole("alert")).toHaveText(
    "Invalid email or password.",
  );
  await expect(page.getByLabel("Email")).toHaveAttribute(
    "aria-invalid",
    "true",
  );
  await expect(page.getByLabel("Email")).toHaveAccessibleDescription(
    "Invalid email or password.",
  );
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "A place to start." }),
  ).toBeVisible();
  await page.getByLabel("New note").fill("   ");
  await page.getByRole("button", { name: "Add note" }).click();
  await expect(page.getByRole("alert")).toHaveText(
    "Write a note with 1 to 2000 characters.",
  );
  await expect(page.getByLabel("New note")).toHaveAttribute(
    "aria-invalid",
    "true",
  );
  await expect(page.getByLabel("New note")).toHaveAccessibleDescription(
    "Write a note with 1 to 2000 characters.",
  );
  const note = `Browser note ${Date.now()}`;
  await page.getByLabel("New note").fill(note);
  await page.getByRole("button", { name: "Add note" }).click();
  const item = page.getByRole("listitem").filter({ hasText: note });
  await expect(item).toBeVisible();
  await expect(page.getByLabel("New note")).not.toHaveAttribute("aria-invalid");
  await page.setViewportSize({ width: 320, height: 700 });
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= window.innerWidth,
    ),
  ).toBe(true);
  await page.reload();
  await expect(item).toBeVisible();
  await item.getByRole("button", { name: "Delete" }).click();
  await expect(item).toHaveCount(0);
  await page.getByRole("link", { name: "Note tools" }).click();
  await expect(
    page.getByRole("heading", { name: "Do more with your notes." }),
  ).toBeVisible();
  await expect(
    page.getByText("The AI assistant has not been enabled for this workspace."),
  ).toBeVisible();
  await page.getByRole("button", { name: "Sign out" }).click();
  await expect(
    page.getByRole("heading", { name: "Welcome back" }),
  ).toBeVisible();
  await page.goto("/admin");
  await expect(
    page.getByRole("heading", { name: "Welcome back" }),
  ).toBeVisible();
  expect(failures).toEqual([]);
});
