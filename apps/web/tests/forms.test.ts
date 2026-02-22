import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { expect, test, describe } from 'vitest';

const srcDir = resolve(__dirname, '../src');

function getSourceContent(relativePath: string): string {
  return readFileSync(resolve(srcDir, relativePath), 'utf-8');
}

describe('SigV4 Form Attributes', () => {
  const cases = [
    {
      name: 'login form',
      path: 'pages/login.astro',
      expects: [
        'data-api-form',
        'action="/api/auth/login"',
        'method="POST"',
        'data-redirect="/"',
      ],
    },
    {
      name: 'calories/new form',
      path: 'pages/calories/new.astro',
      expects: [
        'data-api-form',
        'action="/api/calories/log"',
        'method="POST"',
        'data-redirect="/calories"',
      ],
    },
    {
      name: 'calories/settings form',
      path: 'pages/calories/settings.astro',
      expects: [
        'data-api-form',
        'action="/api/calories/target"',
        'method="POST"',
        'data-redirect="/calories"',
      ],
    },
    {
      name: 'calories quick add form',
      path: 'modules/calories/components/QuickAdd.astro',
      expects: [
        'data-api-form',
        'action="/api/calories/log"',
        'method="POST"',
      ],
    },
    {
      name: 'expenses/new form',
      path: 'pages/expenses/new.astro',
      expects: [
        'data-api-form',
        'action="/api/expenses/expense"',
        'method="POST"',
        'data-redirect="/expenses"',
      ],
    },
    {
      name: 'expenses/settings form',
      path: 'pages/expenses/settings.astro',
      expects: [
        'data-api-form',
        'action="/api/expenses/settings"',
        'method="POST"',
        'data-redirect="/expenses"',
      ],
    },
    {
      name: 'expenses checklist form',
      path: 'pages/expenses/settings.astro',
      expects: [
        'id="expense-checklist-form"',
        'data-api-form',
        'action="/api/expenses/checklist"',
        'method="POST"',
      ],
    },
    {
      name: 'expenses recurring form',
      path: 'pages/expenses/settings.astro',
      expects: [
        'id="expense-recurring-form"',
        'data-api-form',
        'action="/api/expenses/recurring"',
        'method="POST"',
      ],
    },
    {
      name: 'expenses checklist complete form',
      path: 'modules/expenses/components/ObligationsList.astro',
      expects: [
        'data-api-form',
        'action="/api/expenses/checklist/complete"',
        'method="POST"',
      ],
    },
    {
      name: 'heat/new form',
      path: 'pages/heat/new.astro',
      expects: [
        'data-api-form',
        'action="/api/heat/refill"',
        'method="POST"',
        'data-redirect="/heat"',
      ],
    },
  ];

  cases.forEach(({ name, path, expects }) => {
    test(`${name} has required data attributes`, () => {
      const content = getSourceContent(path);
      expects.forEach((expectation) => {
        expect(content).toContain(expectation);
      });
    });
  });
});

describe('Form Safety Checks', () => {
  test('calories form uses ApiFormHandler pattern', () => {
    const content = getSourceContent('pages/calories/new.astro');
    expect(content).toContain('FormShell');
    expect(content).toContain('data-api-form');
  });

  test('no inline fetch calls in form templates', () => {
    const files = [
      'pages/login.astro',
      'pages/calories/new.astro',
      'pages/calories/settings.astro',
      'modules/calories/components/QuickAdd.astro',
      'pages/expenses/new.astro',
      'pages/expenses/settings.astro',
      'modules/expenses/components/ObligationsList.astro',
      'pages/heat/new.astro',
    ];

    files.forEach((path) => {
      const content = getSourceContent(path);
      const hasDirectFetch = content.includes('fetch(') && !content.includes('//');
      expect(hasDirectFetch).toBe(false);
    });
  });
});
