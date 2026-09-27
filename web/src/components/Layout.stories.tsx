import { App } from "@inertiajs/react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect } from "storybook/test";
import { Layout } from "./Layout";

const meta = {
  title: "Layouts/Workspace",
  component: Layout,
  tags: ["autodocs", "ai-generated"],
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "Authenticated workspace navigation and the page's main content. The real Inertia provider supplies synthetic page data; route handlers remain in the application.",
      },
    },
  },
  args: {
    user: { id: "1", email: "person@example.test", admin: false },
    children: <h1 className="text-3xl font-semibold">Workspace</h1>,
  },
  argTypes: { user: { control: false }, children: { control: false } },
  decorators: [
    (Story) => (
      <App
        initialPage={{
          component: "Story",
          url: "/",
          version: null,
          props: { csrfToken: "storybook-only", errors: {} },
          flash: {},
          rescuedProps: [],
          rememberedState: {},
        }}
        initialComponent={() => null}
        resolveComponent={() => () => null}
      >
        {() => <Story />}
      </App>
    ),
  ],
} satisfies Meta<typeof Layout>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Member: Story = {
  play: async ({ canvas }) => {
    await expect(
      canvas.queryByRole("link", { name: "Admin" }),
    ).not.toBeInTheDocument();
  },
};
export const Administrator: Story = {
  args: { user: { id: "2", email: "admin@example.test", admin: true } },
  play: async ({ canvas }) => {
    await expect(canvas.getByRole("link", { name: "Admin" })).toHaveAttribute(
      "href",
      "/admin",
    );
  },
};
export const Narrow: Story = {
  args: Administrator.args,
  decorators: [
    (Story) => (
      <div style={{ width: 320 }}>
        <Story />
      </div>
    ),
  ],
  parameters: {
    docs: {
      description: {
        story:
          "At 320 px the brand and navigation wrap inside the workspace gutter without horizontal overflow.",
      },
    },
  },
  play: async ({ canvas }) => {
    const nav = canvas.getByRole("navigation", { name: "Main" });
    const header = nav.parentElement!;
    await expect(header.scrollWidth).toBeLessThanOrEqual(header.clientWidth);
    await expect(
      canvas.getByRole("button", { name: "Sign out" }),
    ).toBeEnabled();
  },
};
