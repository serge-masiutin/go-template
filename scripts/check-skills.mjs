import { readFileSync, readdirSync, existsSync } from "node:fs";
import { createHash } from "node:crypto";
import path from "node:path";

const root = process.cwd();
const lock = JSON.parse(readFileSync("config/skills-lock.json", "utf8"));
const sha = (contents) => createHash("sha256").update(contents).digest("hex");
function files(directory) {
  return readdirSync(directory, { withFileTypes: true }).flatMap((entry) => {
    const filename = path.join(directory, entry.name);
    return entry.isDirectory() ? files(filename) : [filename];
  });
}
const failures = [];
function check(condition, message) {
  if (!condition) failures.push(message);
}
function markdown(filename) {
  const text = readFileSync(filename, "utf8");
  let fence = null;
  const prose = [];
  for (const line of text.split("\n")) {
    const marker = line.match(/^\s*(`{3,}|~{3,})/);
    if (marker) {
      if (!fence) fence = marker[1];
      else if (marker[1][0] === fence[0] && marker[1].length >= fence.length)
        fence = null;
    } else if (!fence) prose.push(line);
  }
  check(fence === null, `${filename}: unclosed code fence`);
  const withoutInlineCode = prose.join("\n").replace(/(`+)[\s\S]*?\1/g, "");
  for (const [, raw] of withoutInlineCode.matchAll(/\[[^\]]*\]\(([^)]+)\)/g)) {
    const url = raw.replace(/^<|>$/g, "").split(/\s+"/)[0].split("#")[0];
    if (!url || /^[a-z][\w+.-]*:/.test(url)) continue;
    check(
      existsSync(path.resolve(path.dirname(filename), decodeURI(url))),
      `${filename}: missing local link ${url}`,
    );
  }
}
const names = lock.skills.map((skill) => skill.name).sort();
const actual = readdirSync(".agents/skills", { withFileTypes: true })
  .filter((e) => e.isDirectory())
  .map((e) => e.name)
  .sort();
check(
  JSON.stringify(names) === JSON.stringify(actual),
  "Skill catalog differs from lock",
);
check(new Set(names).size === names.length, "Duplicate skill name");
for (const skill of lock.skills) {
  const directory = path.join(".agents/skills", skill.name);
  const actualFiles = files(directory)
    .map((f) => path.relative(directory, f))
    .sort();
  check(
    JSON.stringify(actualFiles) ===
      JSON.stringify(Object.keys(skill.files).sort()),
    `${skill.name}: file inventory differs`,
  );
  for (const [filename, digest] of Object.entries(skill.files)) {
    const target = path.join(directory, filename);
    check(
      existsSync(target) && sha(readFileSync(target)) === digest,
      `${skill.name}/${filename}: content differs from reviewed lock`,
    );
  }
  const entry = readFileSync(path.join(directory, "SKILL.md"), "utf8");
  check(entry.startsWith("---\n"), `${skill.name}: missing frontmatter`);
  check(
    new RegExp(`^name: ${skill.name}$`, "m").test(entry),
    `${skill.name}: invalid name`,
  );
  check(/^description: .+/m.test(entry), `${skill.name}: missing description`);
  if (skill.mode !== "upstream") {
    for (const key of ["archive", "patch"]) {
      check(
        sha(readFileSync(skill[key])) === skill[key + "Sha256"],
        `${skill.name}: ${key} integrity failed`,
      );
    }
    if (skill.mode === "preserved") {
      check(
        JSON.stringify(skill.sourceFiles) === JSON.stringify(skill.files),
        `${skill.name}: preserved source changed`,
      );
    } else {
      for (const filename of actualFiles.filter((f) => f.endsWith(".md"))) {
        const target = path.join(directory, filename);
        markdown(target);
        check(
          !/^```(?:ruby|erb)\s*$/m.test(readFileSync(target, "utf8")),
          `${target}: executable Rails example remains`,
        );
      }
      for (const target of Object.values(skill.mapping).filter(Boolean)) {
        check(
          existsSync(path.join(directory, target)),
          `${skill.name}: mapped target missing ${target}`,
        );
      }
    }
  }
}
for (const filename of [
  "README.md",
  "AGENTS.md",
  "THIRD_PARTY.md",
  ...files("docs"),
].filter((f) => f.endsWith(".md")))
  markdown(filename);
if (failures.length) {
  console.error(failures.join("\n"));
  process.exit(1);
}
console.log(
  `Verified ${names.length} skills: exact upstream bytes, adaptation provenance, local links and code fences.`,
);
