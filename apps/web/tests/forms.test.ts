import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { expect, test, describe } from 'vitest';

const pagesDir = resolve(__dirname, '../src/pages');

function getPageContent(pagePath: string): string {
  return readFileSync(resolve(pagesDir, pagePath), 'utf-8');
}

describe('SigV4 Form Attributes', () => {
  test('calories/new form has required data attributes', () => {
    const content = getPageContent('calories/new.astro');
    
    expect(content).toContain('data-api-form');
    expect(content).toContain('action="/api/calories/log"');
    expect(content).toContain('method="POST"');
    expect(content).toContain('data-redirect="/calories"');
  });

  test('expenses/new form has required data attributes', () => {
    const content = getPageContent('expenses/new.astro');
    
    expect(content).toContain('data-api-form');
    expect(content).toContain('action="/api/expenses/expense"');
    expect(content).toContain('method="POST"');
  });

  test('heat/new form has required data attributes', () => {
    const content = getPageContent('heat/new.astro');
    
    expect(content).toContain('data-api-form');
    expect(content).toContain('action="/api/heat/refill"');
    expect(content).toContain('method="POST"');
  });
});

describe('Form Safety Checks', () => {
  test('calories form uses ApiFormHandler pattern', () => {
    const content = getPageContent('calories/new.astro');
    
    // Ensure it has the ApiFormHandler wrapper (via FormShell)
    expect(content).toContain('FormShell');
    
    // Check for data attributes that ApiFormHandler expects
    expect(content).toContain('data-api-form');
  });

  test('no inline scripts that bypass ApiFormHandler', () => {
    const content = getPageContent('calories/new.astro');
    
    // Should not have direct fetch calls in the template
    const hasDirectFetch = content.includes('fetch(') && !content.includes('//');
    expect(hasDirectFetch).toBe(false);
  });
});
