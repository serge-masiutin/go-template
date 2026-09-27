import type { Meta, StoryObj } from "@storybook/react-vite";
import { Notice } from "./Notice";
const meta = {
  title: "UI/Notice",
  component: Notice,
  tags: ["autodocs"],
  parameters: { layout: "padded" },
  args: { children: "Write a note with 1 to 2000 characters." },
} satisfies Meta<typeof Notice>;
export default meta;
export const ValidationError: StoryObj<typeof meta> = {};
