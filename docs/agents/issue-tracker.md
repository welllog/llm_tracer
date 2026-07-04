# Issue Tracker: Local Markdown

Issues for this repository are tracked as markdown files within the `.scratch/` directory.

## Configuration

- **Type**: Local Markdown
- **Directory**: `.scratch/`
- **Pattern**: `.scratch/<feature>/<issue-id>-<slug>.md`

## Workflow

1. **Creation**: Use the `to-issues` skill to create a new markdown file in a feature subdirectory.
2. **Triage**: The `triage` skill scans these files and updates their frontmatter/tags.
3. **Internal PRs**: Not applicable for local markdown tracking.
