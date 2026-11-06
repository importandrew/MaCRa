#!/usr/bin/env node
// PreToolUse hook for Bash: blocks force pushes, pushes to the protected
// branch and PR merges, regardless of how the command is spelled.
// Exit code 2 blocks the tool call and shows stderr to Claude.
const { execFileSync } = require('child_process');

const PROTECTED = new Set(['master', 'main']);
// git global options that consume the next argument.
const GIT_OPTS_WITH_ARG = new Set(['-C', '-c', '--git-dir', '--work-tree', '--namespace', '--exec-path']);

function tokenize(cmd) {
  const tokens = [];
  let cur = '';
  let quote = null;
  let has = false;
  const push = () => {
    if (has) tokens.push(cur);
    cur = '';
    has = false;
  };
  for (let i = 0; i < cmd.length; i++) {
    const ch = cmd[i];
    if (quote) {
      if (ch === quote) quote = null;
      else if (ch === '\\' && quote === '"' && i + 1 < cmd.length) cur += cmd[++i];
      else cur += ch;
      continue;
    }
    if (ch === "'" || ch === '"') {
      quote = ch;
      has = true;
    } else if (ch === '\\' && i + 1 < cmd.length) {
      cur += cmd[++i];
      has = true;
    } else if (/\s/.test(ch)) {
      push();
    } else if (';&|()`'.includes(ch)) {
      push();
      tokens.push(';');
    } else {
      cur += ch;
      has = true;
    }
  }
  push();
  return tokens;
}

function segments(tokens) {
  const out = [[]];
  for (const t of tokens) {
    if (t === ';') out.push([]);
    else out[out.length - 1].push(t);
  }
  return out.filter((s) => s.length);
}

function currentBranch(dir) {
  try {
    return execFileSync('git', ['-C', dir, 'rev-parse', '--abbrev-ref', 'HEAD'], {
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'ignore'],
    }).trim();
  } catch {
    return '';
  }
}

function checkGitPush(args, dir) {
  const positional = [];
  for (let i = 0; i < args.length; i++) {
    const a = args[i];
    if (a === '--') {
      positional.push(...args.slice(i + 1));
      break;
    }
    if (a.startsWith('--force') || a === '--mirror' || a === '--all' || a === '--delete') {
      return `git push ${a} is not allowed`;
    }
    if (/^-[a-zA-Z]*[fd]/.test(a) && !a.startsWith('--')) {
      return `git push ${a} (force/delete) is not allowed`;
    }
    if (['--repo', '--receive-pack', '--exec', '--push-option', '-o'].includes(a)) i++;
    else if (!a.startsWith('-')) positional.push(a);
  }
  const refspecs = positional.slice(1);
  const branch = currentBranch(dir);
  if (refspecs.length === 0 && PROTECTED.has(branch)) {
    return `git push from ${branch} is not allowed; work on a feature/<name> branch`;
  }
  for (const spec of refspecs) {
    if (spec.startsWith('+')) return `force refspec ${spec} is not allowed`;
    let dst = spec.includes(':') ? spec.split(':').pop() : spec;
    if (dst === 'HEAD') dst = branch;
    dst = dst.replace(/^refs\/heads\//, '');
    if (PROTECTED.has(dst)) return `pushing to ${dst} is not allowed; only the owner merges into it`;
  }
  return null;
}

function checkSegment(seg, cwd) {
  // Strip leading env assignments and wrappers.
  let i = 0;
  while (i < seg.length && /^[A-Za-z_][A-Za-z0-9_]*=/.test(seg[i])) i++;
  while (i < seg.length && ['sudo', 'command', 'exec', 'env', 'nohup', 'time'].includes(seg[i])) i++;
  const words = seg.slice(i);
  if (!words.length) return null;
  const bin = words[0].split('/').pop();

  if (bin === 'git') {
    let dir = cwd;
    let j = 1;
    while (j < words.length && words[j].startsWith('-')) {
      if (words[j] === '-C') dir = words[j + 1] || dir;
      j += GIT_OPTS_WITH_ARG.has(words[j]) ? 2 : 1;
    }
    if (words[j] === 'push') return checkGitPush(words.slice(j + 1), dir);
  }
  if (bin === 'gh') {
    // Global flags such as -R <repo> take values, so locate subcommands by name.
    const pr = words.indexOf('pr');
    if (pr > 0 && words.slice(pr + 1).find((w) => !w.startsWith('-')) === 'merge') {
      return 'gh pr merge is not allowed; only the owner merges';
    }
    if (words.includes('api') && words.some((w) => /pulls\/[^/\s]+\/merge|mergePullRequest/.test(w))) {
      return 'merging PRs through gh api is not allowed';
    }
  }
  return null;
}

const SHELLS = new Set(['sh', 'bash', 'zsh', 'dash', 'ksh', 'eval']);

function check(cmd, cwd, depth = 0) {
  if (depth > 3) return null;
  for (const seg of segments(tokenize(cmd))) {
    const reason = checkSegment(seg, cwd);
    if (reason) return reason;
    // Arguments of a nested shell are commands too: sh -c "...", eval "...".
    // Other quoted strings (commit messages, echo) are plain data.
    if (seg.some((w) => SHELLS.has(w.split('/').pop()))) {
      for (const t of seg) {
        const nested = /\s/.test(t) ? check(t, cwd, depth + 1) : null;
        if (nested) return nested;
      }
    }
  }
  return null;
}

let input = '';
process.stdin.on('data', (d) => (input += d));
process.stdin.on('end', () => {
  let data;
  try {
    data = JSON.parse(input);
  } catch {
    process.exit(0);
  }
  const cmd = data?.tool_input?.command;
  if (typeof cmd !== 'string') process.exit(0);
  const reason = check(cmd, data.cwd || process.cwd());
  if (reason) {
    process.stderr.write(`Blocked by .claude/hooks/git-guard.js: ${reason}. See "Working rules" in CLAUDE.md.\n`);
    process.exit(2);
  }
  process.exit(0);
});
