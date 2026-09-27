import { spawn, execFileSync } from "node:child_process";
import { mkdir, rm, access } from "node:fs/promises";
import { setTimeout as delay } from "node:timers/promises";

await mkdir("tmp", { recursive: true });
await rm("tmp/vite.hot", { force: true });
execFileSync("go", ["build", "-o", "tmp/server", "./cmd/server"], {
  stdio: "inherit",
});
const children = new Set();
let stopping = false;
async function stop(code) {
  if (stopping) return;
  stopping = true;
  const exits = [...children].map(
    (child) => new Promise((resolve) => child.once("exit", resolve)),
  );
  for (const child of children) child.kill("SIGTERM");
  const deadline = setTimeout(() => {
    for (const child of children) child.kill("SIGKILL");
  }, 11_000);
  await Promise.all(exits);
  clearTimeout(deadline);
  await rm("tmp/vite.hot", { force: true });
  process.exitCode = code;
}
function start(command, args) {
  const child = spawn(command, args, { stdio: "inherit" });
  children.add(child);
  child.once("error", (error) => {
    children.delete(child);
    console.error(error.message);
    void stop(1);
  });
  child.once("exit", (code, signal) => {
    children.delete(child);
    if (!stopping) {
      console.error(
        `${command} stopped (${signal ?? code}). Stopping development servers.`,
      );
      void stop(code || 1);
    }
  });
  return child;
}
process.once("SIGINT", () => void stop(130));
process.once("SIGTERM", () => void stop(143));
start(process.execPath, [
  "node_modules/vite/bin/vite.js",
  "--host",
  "localhost",
]);
let ready = false;
for (let attempt = 0; attempt < 100 && !stopping; attempt++) {
  try {
    await access("tmp/vite.hot");
    ready = true;
    break;
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  await delay(100);
}
if (!stopping) {
  if (ready) start("./tmp/server", []);
  else {
    console.error("Vite did not become ready within 10 seconds.");
    await stop(1);
  }
}
