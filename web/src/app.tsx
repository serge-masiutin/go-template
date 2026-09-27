import { createInertiaApp } from "@inertiajs/react";
import { createRoot } from "react-dom/client";
import Login from "./pages/Login";
import Home from "./pages/Home";
import Admin from "./pages/Admin";
import Tools from "./pages/Tools";
import "./styles.css";
import "./types";

const pages = { Login, Home, Admin, Tools };

createInertiaApp({
  resolve(name) {
    if (!(name in pages)) throw new Error(`Unknown page: ${name}`);
    return pages[name as keyof typeof pages];
  },
  setup({ el, App, props }) {
    createRoot(el).render(<App {...props} />);
  },
});
