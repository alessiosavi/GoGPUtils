# GoGPUtils Agent Discovery System - Usage & Maintenance Guide

This document describes how to maintain and extend the two-layer agent discovery system for GoGPUtils.

## Overview

The system consists of two complementary layers:

- **Layer 1 (Awareness)**: Static reference skill (`docs/superpowers/skills/gputils/SKILL.md`) that teaches agents GoGPUtils exists and maps common problems to packages.
- **Layer 2 (Precision)**: Dynamic Graphify MCP server backed by `graphify-out/graph.json` that enables semantic queries for exact function signatures.

## Configuration

### MCP Server Setup

The Graphify MCP server is configured in `.opencode/mcp.json`:

```json
{
  "mcpServers": {
    "gputils-graphify": {
      "command": "python3",
      "args": ["-m", "graphify.serve", "graphify-out/graph.json"],
      "description": "GoGPUtils knowledge graph for semantic function discovery"
    }
  }
}
```

**For Claude Desktop**: Add the equivalent entry to `claude_desktop_config.json` with an absolute path to `graphify-out/graph.json`.

### Agent Context

Add the following to `AGENTS.md` (or `.opencode/AGENTS.md`):

```markdown
## GoGPUtils Agent Discovery

When working on Go projects that import `github.com/alessiosavi/GoGPUtils`:

1. **Load skill**: The `gputils-reference` skill provides awareness of available utilities.
2. **Query MCP**: For exact function signatures, use the `gputils-graphify` MCP server:
   - `query_graph`: "functions for filtering slices"
   - `get_node`: "sliceutil_Filter" for exact signature
   - `shortest_path`: "sliceutil_Filter" → "sliceutil_GroupBy" for composite flows

**Configuration**: Ensure `.opencode/mcp.json` has the `gputils-graphify` server configured.
```

## Agent Workflow

When an agent encounters a Go programming task:

1. Skill triggers on Go task context (description matches "writing Go code...").
2. Agent reads skill → learns GoGPUtils exists and sees quick-reference table.
3. Agent uses quick-reference to identify likely package.
4. Agent queries MCP server (`query_graph`) for exact function.
5. Agent calls `get_node` or `shortest_path` for precise signature.
6. Agent generates code with correct import and usage.
7. Integration test (`test_agent_discovery.sh`) can be run to verify the system state.

## Maintenance Checklist

- [ ] **After adding new functions**: Update `SKILL.md` function tables (Key Functions by Package section).
- [ ] **After adding new packages**: Update quick-reference table and add package section with key functions.
- [ ] **After code changes**: Run `graphify . --update` (or equivalent incremental command) to refresh the knowledge graph.
- [ ] **Monthly**: Run `./test_agent_discovery.sh` to verify end-to-end health.
- [ ] **Quarterly**: Review `graphify-out/GRAPH_REPORT.md` for graph health, god nodes, and community structure.

## Troubleshooting

**Skill not loading**:

- Check YAML frontmatter: must start with `---` and contain `name: gputils-reference` and `description:`.
- Verify file is in `.opencode/skills/gputils/SKILL.md` (or global skills directory).

**MCP server not starting**:

- Verify `graphify-out/graph.json` exists and is valid JSON (run `python3 -c "import json; json.load(open('graphify-out/graph.json'))"`).
- Check that the `graphify.serve` module is available (`python3 -m graphify.serve --help`).
- Ensure the path in `mcp.json` is correct (relative or absolute as required by your client).

**Graph missing nodes or stale**:

- Run `graphify . --update` (or full rebuild: `rm -rf graphify-out/ && graphify .`).
- Check `.graphify_incremental.json` mtime for freshness.
- Verify core packages appear in `graph.json` node IDs (lowercase "package_symbol" or "package_filename_go" format).

**False positives / wrong functions returned**:

- Graphify returns confidence scores; prefer EXTRACTED edges (score 1.0) over INFERRED.
- Update the skill's quick-reference table if new canonical names emerge.

**Test failures**:

- Run `./test_agent_discovery.sh` and inspect the failing test number.
- Most failures are due to missing files or stale graph; fix and re-run.

## Extending the System

1. **New utility package**: Add section to `SKILL.md` quick-reference and detailed function table. Rebuild graph.
2. **New function in existing package**: Add to the relevant table in `SKILL.md`. Run incremental graph update.
3. **Changing MCP tool names**: Update both the `mcp.json` tools list and any documentation in `AGENTS.md`.
4. **Versioning**: The skill and graph are versioned with the repo; no separate versioning needed.

## Related Files

- `docs/superpowers/skills/gputils/SKILL.md` — Layer 1 reference skill
- `graphify-out/graph.json` — Knowledge graph (Layer 2)
- `.opencode/mcp.json` — MCP server configuration
- `test_agent_discovery.sh` — Integration test
- `docs/superpowers/plans/2026-07-16-agent-discovery-impl-plan.md` — Original implementation plan
- `docs/superpowers/specs/2026-07-16-gputils-agent-discovery-design.md` — Design specification

## Support

For issues with the discovery system itself (not GoGPUtils code), open an issue with the label `agent-discovery` and include:

- Output of `./test_agent_discovery.sh`
- Relevant section of `graphify-out/GRAPH_REPORT.md`
- MCP client logs (if server fails to start)

---

_End of usage guide._
