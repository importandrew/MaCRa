import js from '@eslint/js';
import { fixupPluginRules } from '@eslint/compat';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import prettierRecommended from 'eslint-plugin-prettier/recommended';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import { defineConfig, globalIgnores } from 'eslint/config';
import globals from 'globals';
import tseslint from 'typescript-eslint';

// eslint-plugin-react 7.37 and eslint-plugin-jsx-a11y 6.10 don't support
// ESLint 10 yet: they call context methods that ESLint 10 removed.
// fixupPluginRules() adds those methods back. Drop the wrapper (and the
// matching "overrides" in package.json) once both plugins support ESLint 10.
const withFixedPlugin = (name, plugin, config) => ({
  ...config,
  plugins: { [name]: fixupPluginRules(plugin) },
});

export default defineConfig(
  globalIgnores(['dist']),
  {
    files: ['**/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2021,
      sourceType: 'module',
      globals: { ...globals.browser, ...globals.node, ...globals.es2021 },
      parserOptions: { ecmaFeatures: { jsx: true } },
    },
    settings: { react: { version: 'detect' } },
  },
  js.configs.recommended,
  withFixedPlugin('react', react, react.configs.flat.recommended),
  // The full "recommended" set of eslint-plugin-react-hooks 7: the classic
  // hook rules plus the React Compiler rules (set-state-in-effect, refs,
  // immutability, purity, ...).
  reactHooks.configs.flat.recommended,
  withFixedPlugin('jsx-a11y', jsxA11y, jsxA11y.flatConfigs.recommended),
  prettierRecommended,
  tseslint.configs.recommended,
  {
    plugins: { 'react-refresh': reactRefresh },
    rules: {
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
      'react/react-in-jsx-scope': 'off', // Not needed in React 17+
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'warn',
      'prettier/prettier': 'warn',
      'jsx-a11y/anchor-is-valid': 'off', // Disable if using Next.js Link
      'no-unused-vars': 'off',
      'no-console': 'warn',
    },
  }
);
