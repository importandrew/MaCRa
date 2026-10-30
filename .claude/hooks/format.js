#!/usr/bin/env node
// PostToolUse hook for Edit|Write: formats the edited file in place.
// Go files go through gofmt; frontend files through the repo's local prettier.
// Never fails the tool call: formatting problems are ignored.
const { execFileSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const PRETTIER_EXT = new Set(['.js', '.jsx', '.ts', '.tsx', '.cjs', '.mjs', '.css', '.json', '.html']);

let input = '';
process.stdin.on('data', (d) => (input += d));
process.stdin.on('end', () => {
  try {
    const data = JSON.parse(input);
    const file = data?.tool_response?.filePath || data?.tool_input?.file_path;
    const root = process.env.CLAUDE_PROJECT_DIR || data.cwd;
    if (!file || !root || !fs.existsSync(file)) return;

    const rel = path.relative(root, file);
    if (rel.startsWith('..') || path.isAbsolute(rel)) return;
    const ext = path.extname(file);
    const run = (bin, args) => execFileSync(bin, args, { stdio: 'ignore', timeout: 20000 });

    if (ext === '.go') {
      run('gofmt', ['-w', file]);
    } else if (rel.startsWith(`frontend${path.sep}`) && PRETTIER_EXT.has(ext)) {
      const prettier = path.join(root, 'frontend', 'node_modules', '.bin', 'prettier');
      if (fs.existsSync(prettier)) run(prettier, ['--write', '--ignore-unknown', file]);
    }
  } catch {
    // Formatting is best-effort.
  }
});
