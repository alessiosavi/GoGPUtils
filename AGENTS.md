## GoGPUtils Agent Discovery

When working on Go projects that import `github.com/alessiosavi/GoGPUtils`:

1. **Load skill**: The `gputils-reference` skill provides awareness of available utilities
2. **Query MCP**: For exact function signatures, use the `gputils-graphify` MCP server:
   - `query_graph`: "functions for filtering slices"
   - `get_node`: "sliceutil_Filter" for exact signature
   - `shortest_path`: "sliceutil_Filter" → "sliceutil_GroupBy" for composite flows

**Configuration**: Ensure `.opencode/mcp.json` has the `gputils-graphify` server configured.
