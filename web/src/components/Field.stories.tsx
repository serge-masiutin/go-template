import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, fn } from "storybook/test";
import { Field } from "./Field";

const meta = {
  title: "UI/Field",
  component: Field,
  tags: ["autodocs", "ai-generated"],
  parameters: {
    layout: "padded",
    docs: {
      description: {
        component:
          "A named text value with optional instructions and server validation. Native input semantics and form ownership stay with the consuming form.",
      },
    },
  },
  args: {
    id: "example-field",
    name: "email",
    label: "Email",
    onChange: fn(),
  },
  argTypes: {
    type: {
      control: "select",
      options: ["text", "email", "password"],
      table: { category: "Content" },
    },
    multiline: { control: "boolean", table: { category: "Content" } },
    label: { control: "text", table: { category: "Content" } },
    hint: { control: "text", table: { category: "Content" } },
    error: { control: "text", table: { category: "Validation" } },
    disabled: { control: "boolean", table: { category: "State" } },
    readOnly: { control: "boolean", table: { category: "State" } },
    required: { control: "boolean", table: { category: "Validation" } },
    onChange: { control: false, table: { category: "Events" } },
  },
} satisfies Meta<typeof Field>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Default: Story = {
  args: { type: "email", autoComplete: "username" },
};
export const WithHint: Story = {
  args: {
    type: "email",
    hint: "Use the email address associated with your workspace.",
  },
  play: async ({ canvas }) => {
    await expect(canvas.getByLabelText("Email")).toHaveAccessibleDescription(
      "Use the email address associated with your workspace.",
    );
  },
};
export const ValidationError: Story = {
  args: {
    type: "email",
    hint: "Use your work address.",
    error: "Check this email address.",
  },
  play: async ({ canvas }) => {
    const input = canvas.getByLabelText("Email");
    await expect(input).toHaveAttribute("aria-invalid", "true");
    await expect(input).toHaveAccessibleDescription(
      "Use your work address. Check this email address.",
    );
  },
};
export const Filled: Story = {
  args: { type: "email", defaultValue: "person@example.test" },
};
export const Disabled: Story = {
  args: { type: "email", disabled: true, defaultValue: "person@example.test" },
};
export const ReadOnly: Story = {
  args: { type: "email", readOnly: true, defaultValue: "person@example.test" },
};
export const Required: Story = { args: { type: "email", required: true } };
export const Password: Story = {
  args: {
    label: "Password",
    name: "password",
    type: "password",
    autoComplete: "current-password",
    defaultValue: "synthetic-password",
  },
};
export const Multiline: Story = {
  args: {
    multiline: true,
    label: "New note",
    name: "body",
    rows: 4,
    maxLength: 2000,
    required: true,
  },
  play: async ({ canvas, userEvent, args }) => {
    const input = canvas.getByLabelText("New note");
    await userEvent.type(input, "First line{Enter}Second line");
    await expect(input).toHaveValue("First line\nSecond line");
    await expect(args.onChange).toHaveBeenCalled();
  },
};
export const KeyboardEntry: Story = {
  args: { type: "email", autoComplete: "username" },
  play: async ({ canvas, userEvent, args }) => {
    await userEvent.tab();
    const input = canvas.getByLabelText("Email");
    await expect(input).toHaveFocus();
    await userEvent.type(input, "person@example.test");
    await expect(input).toHaveValue("person@example.test");
    await expect(args.onChange).toHaveBeenCalled();
  },
};
