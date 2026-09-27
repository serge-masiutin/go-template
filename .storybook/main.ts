import type { StorybookConfig } from "@storybook/react-vite";
const config: StorybookConfig = {
  core: { disableTelemetry: true },
  stories: ["../web/src/**/*.stories.tsx"],
  addons: ["@storybook/addon-a11y", "@storybook/addon-docs"],
  framework: "@storybook/react-vite",
  async viteFinal(config) {
    // Storybook must not point the Go server at an unrelated development server.
    config.plugins = config.plugins?.filter(
      (plugin) =>
        !(
          plugin &&
          typeof plugin === "object" &&
          "name" in plugin &&
          plugin.name === "gonertia-hot-file"
        ),
    );
    return config;
  },
  staticDirs: ["../web/public"],
};
export default config;
