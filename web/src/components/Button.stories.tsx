import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, fn } from "storybook/test";
import { Button } from "./Button";
const meta = {
  title: "UI/Button",
  component: Button,
  tags: ["autodocs", "ai-generated"],
  parameters: { layout: "centered" },
  argTypes: {
    variant: { control: "inline-radio", options: ["primary", "secondary"] },
    onClick: { control: false },
  },
  args: { children: "Continue", onClick: fn() },
} satisfies Meta<typeof Button>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Primary: Story = {};
export const Secondary: Story = { args: { variant: "secondary" } };
export const Disabled: Story = { args: { disabled: true } };

export const Working: Story = {
  args: { children: "Signing in…", disabled: true },
};
export const KeyboardAction: Story = {
  play: async ({ canvas, userEvent, args }) => {
    await userEvent.tab();
    await expect(canvas.getByRole("button")).toHaveFocus();
    await userEvent.keyboard("{Enter}");
    await expect(args.onClick).toHaveBeenCalledOnce();
  },
};
export const CssCheck: Story = {
  tags: ["!autodocs", "!dev"],
  play: async ({ canvas }) => {
    await expect(
      getComputedStyle(canvas.getByRole("button")).backgroundColor,
    ).toBe("rgb(37, 98, 72)");
  },
};
