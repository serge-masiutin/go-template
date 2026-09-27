import { readFileSync, writeFileSync, readdirSync } from "node:fs";
import { execFileSync } from "node:child_process";
import path from "node:path";

const [modulePath, ...extra] = process.argv.slice(2);
if (
  extra.length ||
  !modulePath ||
  !/^[a-z0-9.-]+\.[a-z]+\/[A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+$/.test(modulePath)
) {
  throw new Error(
    "Usage: bin/configure github.com/owner/repository (three path segments)",
  );
}
const oldModule = JSON.parse(
  execFileSync("go", ["mod", "edit", "-json"], { encoding: "utf8" }),
).Module.Path;
const name = modulePath.split("/").at(-1).toLowerCase();
function goFiles(directory) {
  return readdirSync(directory, { withFileTypes: true }).flatMap((entry) => {
    const filename = path.join(directory, entry.name);
    return entry.isDirectory()
      ? goFiles(filename)
      : filename.endsWith(".go")
        ? [filename]
        : [];
  });
}
const changes = ["cmd", "internal", "tests"]
  .flatMap(goFiles)
  .map((filename) => [
    filename,
    readFileSync(filename, "utf8").replaceAll(
      `"${oldModule}/`,
      `"${modulePath}/`,
    ),
  ]);
for (const filename of ["package.json", "package-lock.json"]) {
  const manifest = JSON.parse(readFileSync(filename, "utf8"));
  manifest.name = name;
  if (filename === "package-lock.json") manifest.packages[""].name = name;
  changes.push([filename, JSON.stringify(manifest, null, 2) + "\n"]);
}
execFileSync("go", ["mod", "edit", "-module", modulePath]);
for (const [filename, contents] of changes) writeFileSync(filename, contents);
console.log(
  `Configured ${modulePath}. Run bin/ci and review the diff before committing.`,
);
