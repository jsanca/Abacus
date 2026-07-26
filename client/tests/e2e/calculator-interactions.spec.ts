import { expect, test } from '@playwright/test';
import { CalculatorPage } from './calculator.page';

test.describe('Calculator interaction and validation', () => {
  test('rejects invalid projected values without a calculation request', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    let calculationRequests = 0;
    page.on('request', (request) => { if (request.url().endsWith('/api/v1/calculations')) calculationRequests += 1; });
    await calculator.selectOperation('Division');
    await page.getByLabel('Dividend').fill('10');
    await page.getByLabel('Divisor').fill('0');
    await calculator.calculate();
    await expect(calculator.error()).toContainText('divisor');
    await expect(page.getByLabel('Dividend')).toHaveValue('10');
    await expect(page.getByLabel('Divisor')).toHaveValue('0');
    expect(calculationRequests).toBe(0);

    await calculator.selectOperation('Square Root');
    await page.getByLabel('Number').fill('-1');
    await calculator.calculate();
    await expect(calculator.error()).toContainText('zero or greater');
  });

  test('accepts percentage boundaries and rejects invalid percentage values', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await calculator.selectOperation('Percentage');
    expect(await page.getByRole('textbox', { name: 'Percentage' }).evaluate((input) => input.parentElement?.textContent)).toContain('%');
    await page.getByLabel('Base value').fill('200');
    await page.getByRole('textbox', { name: 'Percentage' }).fill('0');
    await calculator.calculate();
    await expect(calculator.result()).toHaveText('0');
    await calculator.clear();
    await calculator.selectOperation('Percentage');
    await page.getByLabel('Base value').fill('200');
    await page.getByRole('textbox', { name: 'Percentage' }).fill('100');
    await calculator.calculate();
    await expect(calculator.result()).toHaveText('200');
    for (const invalidPercentage of ['-1', '101']) {
      await calculator.clear();
      await calculator.selectOperation('Percentage');
      await page.getByLabel('Base value').fill('200');
      await page.getByRole('textbox', { name: 'Percentage' }).fill(invalidPercentage);
      await calculator.calculate();
      await expect(calculator.error()).toContainText('between 0 and 100');
    }
  });

  test('continues results across binary and unary operations with focus', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await calculator.selectOperation('Division');
    await page.getByLabel('Dividend').fill('144');
    await page.getByLabel('Divisor').fill('12');
    await calculator.calculate();
    await calculator.selectOperation('Multiplication');
    await expect(page.getByLabel('First number')).toHaveValue('12');
    await expect(page.getByLabel('Second number')).toHaveValue('');
    await expect(page.getByLabel('Second number')).toBeFocused();
    await page.getByLabel('Second number').fill('3');
    await calculator.calculate();
    await expect(calculator.result()).toHaveText('36');
    await calculator.selectOperation('Square Root');
    await expect(page.getByLabel('Number')).toHaveValue('36');
    await expect(page.getByLabel('Number')).toBeFocused();
    await expect(page.getByLabel('Second number')).toHaveCount(0);
    await calculator.clear();
    await calculator.selectOperation('Square Root');
    await page.getByLabel('Number').fill('144');
    await calculator.calculate();
    await calculator.selectOperation('Addition');
    await expect(page.getByLabel('First number')).toHaveValue('12');
    await expect(page.getByLabel('Second number')).toBeFocused();
  });

  test('clears through the button and Escape while restoring the default operation', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await calculator.selectOperation('Square Root');
    await page.getByLabel('Number').fill('-1');
    await calculator.calculate();
    await calculator.clear();
    await expect(page.getByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
    await expect(page.getByLabel('First number')).toBeFocused();
    await expect(calculator.result()).toHaveText('—');
    await page.getByLabel('First number').fill('9');
    await page.getByLabel('First number').press('Escape');
    await expect(page.getByLabel('First number')).toHaveValue('');
  });

  test('supports calculator shortcuts without intercepting text entry', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    const shortcuts = [['+', 'Addition'], ['-', 'Subtraction'], ['*', 'Multiplication'], ['/', 'Division'], ['^', 'Exponentiation'], ['r', 'Square Root'], ['%', 'Percentage']] as const;
    for (const [shortcut, operationName] of shortcuts) {
      await page.getByRole('button', { name: 'Clear' }).focus();
      await page.keyboard.press(shortcut);
      await expect(page.getByRole('button', { name: operationName })).toHaveAttribute('aria-pressed', 'true');
    }
    await page.getByLabel('Base value').fill('12');
    await page.getByLabel('Base value').press('+');
    await expect(page.getByRole('button', { name: 'Percentage' })).toHaveAttribute('aria-pressed', 'true');
  });

  test('submits the selected calculation with Enter', async ({ page }) => {
    const calculator = new CalculatorPage(page);
    await calculator.open();
    await page.getByLabel('First number').fill('9');
    await page.getByLabel('Second number').fill('6');
    await page.getByLabel('Second number').press('Enter');
    await expect(page.getByText('9 + 6 = 15')).toBeVisible();
  });
});
