# Contributing to Canvus Go API

## Branching and Commit Practices

- Use feature branches for new work (e.g., `feat/feature-name`, `fix/bug-description`).
- Merge via pull requests only.
- Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for all commit messages:
  - Format: `<type>(<scope>): <description>`
  - Example: `feat(api): add new endpoint for user data`
  - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, etc.
  - Use `BREAKING CHANGE:` in the footer for breaking changes.

## Windows/PowerShell Compatibility

- All development and commands must be compatible with Windows and PowerShell.
- Do not use Linux commands, shell syntax, or make Linux-specific assumptions.
- Validate all file paths and scripts for Windows compatibility.

## General Guidelines

- Update `tasks.md` with every new plan or after each major prompt cycle.
- Before every git commit, summarize what was done in `tasks.md`.
- Keep documentation and code up to date.
