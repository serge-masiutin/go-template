# Third-party materials

- `data-systems-architecture`: an original standalone synthesis based on Martin Kleppmann and Chris Riccomini's *Designing Data-Intensive Applications*, second edition. [Sources and coverage](.agents/skills/data-systems-architecture/SOURCES.md) record attribution and limitations. The source book and its illustrations are not included; the template license does not license the book.

- Project foundation and local source skills: [serge-masiutin/rails-template](https://github.com/serge-masiutin/rails-template), MIT. Exact revisions are recorded in [skills-lock.json](config/skills-lock.json).
- Evil Martians Agent Skills: [official catalog](https://evilmartians.com/.well-known/agent-skills/index.json) and [public repository](https://github.com/evilmartians/agent-skills). Copyright © 2026 Evil Martians; [MIT notice](third_party/licenses/evilmartians-agent-skills.txt). Original React/Storybook/general skills and adapted Inertia skills retain attribution.
- `layered-go` is adapted from Vladimir Dementyev's (palkan) `layered-rails` through the recorded Rails-template revision. Source: [layered-rails-skills](https://github.com/palkan/layered-rails-skills). Upstream declares MIT in `layered-rails/.claude-plugin/plugin.json` and its gemspec. The method is based on *Layered Design for Ruby on Rails Applications*. The Go adaptation is not an official edition by the author or Evil Martians. Original files, author links, and changes are retained in the [provenance records](third_party/skills).
- `book-to-skill`: [virgiliojr94/book-to-skill](https://github.com/virgiliojr94/book-to-skill), Copyright © 2025 virgiliojr94, [MIT](.agents/skills/book-to-skill/LICENSE.md). The full upstream package is preserved without changes; its revision and hashes are in `config/skills-lock.json`.
- `clear-writing`: a project skill from the Rails template, preserved without changes. Author and book references remain inside the skill. The template's license does not extend to those books.
- Martian Mono: [Evil Martians](https://github.com/evilmartians/mono), SIL Open Font License 1.1; [full text](web/public/fonts/OFL.txt).

Go and npm dependencies carry their own licenses. Manifests and lockfiles define the exact versions in use. Archives in `third_party/skills` preserve source instructions for review and comparison; active skills live only in `.agents/skills`.
