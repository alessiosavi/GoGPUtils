#!/bin/bash
set -e

echo "=== GoGPUtils Agent Discovery Integration Test ==="
echo

# Test 1: Skill file exists and has content
echo "[TEST 1] Verify skill file"
if [ ! -f "docs/superpowers/skills/gputils/SKILL.md" ]; then
	echo "FAIL: Skill file not found"
	exit 1
fi
if ! grep -q "sliceutil.Filter" "docs/superpowers/skills/gputils/SKILL.md"; then
	echo "FAIL: Skill missing key function reference"
	exit 1
fi
echo "PASS: Skill file exists and references key functions"
echo

# Test 2: Graph exists and is valid JSON
echo "[TEST 2] Verify graph JSON"
if [ ! -f "graphify-out/graph.json" ]; then
	echo "FAIL: Graph file not found"
	exit 1
fi
python3 -c "
import json
with open('graphify-out/graph.json') as f:
    data = json.load(f)
assert 'nodes' in data, 'Missing nodes'
assert 'links' in data, 'Missing links'
assert len(data['nodes']) > 0, 'Empty nodes'
print(f'PASS: Graph has {len(data[\"nodes\"])} nodes, {len(data[\"links\"])} edges')
"
echo

# Test 3: Key nodes exist in graph
echo "[TEST 3] Verify key graph nodes"
python3 -c "
import json
with open('graphify-out/graph.json') as f:
    data = json.load(f)
nodes = {n['id'] for n in data['nodes']}
required = ['sliceutil', 'stringutil', 'mathutil', 'fileutil', 'cryptoutil', 'randutil', 'collection', 'textnorm']
for pkg in required:
    assert any(pkg in n for n in nodes), f'Missing package: {pkg}'
print('PASS: All 8 core packages found in graph')
"
echo

# Test 4: MCP config exists (optional, may not be committed yet)
echo "[TEST 4] Verify MCP configuration"
if [ -f ".opencode/mcp.json" ]; then
	if grep -q "gputils-graphify" ".opencode/mcp.json"; then
		echo "PASS: MCP configuration found"
	else
		echo "WARN: MCP config exists but missing gputils-graphify"
	fi
else
	echo "WARN: MCP config not found (manual configuration may be needed)"
fi
echo

# Test 5: Verify skill can be parsed (YAML frontmatter)
echo "[TEST 5] Verify skill YAML frontmatter"
python3 -c "
import re
with open('docs/superpowers/skills/gputils/SKILL.md') as f:
    content = f.read()
assert content.startswith('---'), 'Missing YAML frontmatter start'
assert 'name: gputils-reference' in content, 'Missing name field'
assert 'description:' in content, 'Missing description field'
print('PASS: YAML frontmatter valid')
"
echo

echo "=== All integration tests passed ==="
