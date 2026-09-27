#!/usr/bin/env python3
"""Refresh reviewed local adaptations; never modify the pinned upstream baselines."""
from pathlib import Path
import difflib
import hashlib
import json
import tarfile

ROOT = Path(__file__).resolve().parents[1]
manifest = ROOT / "config/skills-lock.json"
lock = json.loads(manifest.read_text())


def digest(contents):
    return hashlib.sha256(contents).hexdigest()


for skill in lock["skills"]:
    if skill["mode"] != "adapted":
        continue
    archive = ROOT / skill["archive"]
    if digest(archive.read_bytes()) != skill["archiveSha256"]:
        raise ValueError(f"Corrupted source archive: {skill['name']}")
    with tarfile.open(archive, "r:gz") as stream:
        original = {member.name: stream.extractfile(member).read() for member in stream if member.isfile()}
    if {name: digest(content) for name, content in original.items()} != skill["sourceFiles"]:
        raise ValueError(f"Source file hashes differ: {skill['name']}")
    target = ROOT / ".agents/skills" / skill["name"]
    current = {str(p.relative_to(target)): p.read_bytes() for p in sorted(target.rglob("*")) if p.is_file()}
    patch = []
    mapped = set()
    for filename, contents in original.items():
        replacement = skill["mapping"][filename]
        if replacement is not None and replacement not in current:
            raise ValueError(f"Update the explicit file mapping before removing {replacement}")
        mapped.add(replacement)
        patch.extend(difflib.unified_diff(
            contents.decode().splitlines(True), current.get(replacement, b"").decode().splitlines(True),
            fromfile="original/" + filename, tofile="adapted/" + (replacement or filename),
        ))
    for filename in sorted(current.keys() - mapped):
        patch.extend(difflib.unified_diff([], current[filename].decode().splitlines(True), fromfile="/dev/null", tofile="adapted/" + filename))
    destination = ROOT / skill["patch"]
    destination.write_text("".join(patch))
    skill["patchSha256"] = digest(destination.read_bytes())
    skill["files"] = {name: digest(content) for name, content in current.items()}
manifest.write_text(json.dumps(lock, ensure_ascii=False, indent=2) + "\n")
print("Refreshed adaptation hashes and patches. Review the diff and run bin/ci.")
