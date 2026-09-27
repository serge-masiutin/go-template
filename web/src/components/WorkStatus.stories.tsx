import type { Meta, StoryObj } from "@storybook/react-vite";
import { WorkStatus } from "./WorkStatus";

const meta = {
  title: "Feedback/WorkStatus",
  component: WorkStatus,
  args: { state: "queued" },
  argTypes: {
    state: {
      control: "select",
      options: ["queued", "sending", "sent", "running", "completed", "failed"],
    },
  },
} satisfies Meta<typeof WorkStatus>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Queued: Story = {};
export const Working: Story = { args: { state: "running" } };
export const Completed: Story = { args: { state: "completed" } };
export const Failed: Story = { args: { state: "failed" } };
