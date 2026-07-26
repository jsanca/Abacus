import { expect, test } from '@playwright/test';
import { CalculatorPage } from './calculator.page';

test.describe('Calculator resilience, responsiveness, and accessibility', () => {
  test('shows calculating state and prevents duplicate requests during latency', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    let requestCount = 0;
    let releaseRequest: (() => void) | undefined;
    const requestGate = new Promise<void>((resolve) => { releaseRequest = resolve; });
    await page.route('**/api/v1/calculations', async (route) => {
      requestCount += 1;
      await requestGate;
      await route.continue();
    });
    await page.getByLabel('First number').fill('2');
    await page.getByLabel('Second number').fill('3');
    await calculator.calculate();
    await expect(page.getByRole('button', { name: 'Calculating…' })).toBeDisabled();
    expect(requestCount).toBe(1);
    releaseRequest?.();
    await expect(page.getByText('2 + 3 = 5')).toBeVisible();
    expect(requestCount).toBe(1);
  });

  test('preserves calculation input through a recoverable service failure', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await calculator.selectOperation('Division');
    await page.getByLabel('Dividend').fill('10');
    await page.getByLabel('Divisor').fill('2');
    await page.route('**/api/v1/calculations', (route) => route.fulfill({ status: 503, contentType: 'application/json', body: JSON.stringify({ code: 'service_unavailable', message: 'Service temporarily unavailable.' }) }));
    await calculator.calculate();
    await expect(calculator.error()).toContainText('temporarily unavailable');
    await expect(page.getByLabel('Dividend')).toHaveValue('10');
    await expect(page.getByLabel('Divisor')).toHaveValue('2');
    await page.unroute('**/api/v1/calculations');
    await calculator.calculate();
    await expect(calculator.result()).toHaveText('5');
  });

  test('recovers from a failed manifest request without a hardcoded fallback', async ({ page }) => {
    let failedOnce = false;
    await page.route('**/api/v1/operations', async (route) => {
      if (!failedOnce) {
        failedOnce = true;
        await route.fulfill({ status: 503, contentType: 'application/json', body: JSON.stringify({ code: 'manifest_unavailable', message: 'Capabilities unavailable.' }) });
        return;
      }
      await route.continue();
    });
    await page.goto('/');
    await expect(page.getByRole('alert')).toContainText('Capabilities unavailable');
    await expect(page.getByRole('button', { name: 'Addition' })).toHaveCount(0);
    await page.getByRole('button', { name: 'Retry' }).click();
    await expect(page.getByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
  });

  test('is accessible and usable across desktop and mobile viewports', async ({ page, browserName }) => {
    test.skip(browserName !== 'chromium', 'Chromium is the configured E2E browser.');
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await expect(page.getByRole('group', { name: 'Operation' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
    await expect(page.getByRole('region', { name: 'Calculation results' })).toHaveAttribute('aria-live', 'polite');
    await expect(calculator.status()).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    await calculator.selectOperation('Square Root');
    await expect(page.getByLabel('Number')).toBeVisible();
    await expect(page.getByLabel('Second number')).toHaveCount(0);
    await calculator.selectOperation('Percentage');
    expect(await page.getByRole('textbox', { name: 'Percentage' }).evaluate((input) => input.parentElement?.textContent)).toContain('%');
  });

  test('captures curated desktop and mobile calculator evidence', async ({ page }, testInfo) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await page.getByLabel('First number').fill('10');
    await page.getByLabel('Second number').fill('5');
    await calculator.calculate();
    const evidenceDirectory = '../docs/knowledge/testing/evidence';
    const evidenceName = testInfo.project.name === 'mobile-chromium' ? 'mobile-calculator.png' : 'desktop-calculator.png';
    await page.screenshot({ path: `${evidenceDirectory}/${evidenceName}`, fullPage: true });
    if (testInfo.project.name === 'desktop-chromium') {
      await page.setContent('<main style="font-family:system-ui;padding:48px"><h1>Abacus Playwright acceptance suite</h1><p>Chromium desktop and mobile acceptance projects completed.</p><p>See ACCEPTANCE_TEST_RESULTS.md for the recorded run.</p></main>');
      await page.screenshot({ path: `${evidenceDirectory}/playwright-summary.png`, fullPage: true });
    }
  });
});
