import type { Preview } from "@storybook/react-vite";
import "../web/src/styles.css";
const preview: Preview = {
  parameters: { layout: "centered", a11y: { test: "error" } },
};
export default preview;
