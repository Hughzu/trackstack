import { describe, expect, test } from 'vitest';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const runtimePath = resolve(process.cwd(), 'src/layouts/ClientRuntime.astro');
const runtimeContent = readFileSync(runtimePath, 'utf8');

describe('ClientRuntime confirm modal regression guards', () => {
  test('keeps trigger id values available for path-based delete endpoints', () => {
    expect(runtimeContent).toContain('if (modal.idAttribute) {');
    expect(runtimeContent).toContain('modal.pendingId = trigger.getAttribute(modal.idAttribute);');
    expect(runtimeContent).toContain('resolvedEndpoint = config.endpoint.replaceAll(');
  });

  test('does not force no-payload sentinel when an id attribute exists', () => {
    const sentinelIndex = runtimeContent.indexOf('modal.pendingId = "__no_payload__";');
    const idAttributeIndex = runtimeContent.indexOf('if (modal.idAttribute) {');

    expect(idAttributeIndex).toBeGreaterThan(-1);
    expect(sentinelIndex).toBeGreaterThan(idAttributeIndex);
    expect(runtimeContent).toContain('} else if (modal.sendPayload) {');
  });
});
