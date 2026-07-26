import { expect, test } from '@playwright/test';
import { CalculatorPage } from './calculator.page';

const operations = [
  { name: 'Addition', firstLabel: 'First number', secondLabel: 'Second number', first: '10', second: '5', expression: '10 + 5 = 15', result: '15' },
  { name: 'Subtraction', firstLabel: 'First number', secondLabel: 'Second number', first: '5', second: '10', expression: '5 − 10 = -5', result: '-5' },
  { name: 'Multiplication', firstLabel: 'First number', secondLabel: 'Second number', first: '2.5', second: '4', expression: '2.5 × 4 = 10', result: '10' },
  { name: 'Division', firstLabel: 'Dividend', secondLabel: 'Divisor', first: '144', second: '12', expression: '144 ÷ 12 = 12', result: '12' },
  { name: 'Exponentiation', firstLabel: 'Base', secondLabel: 'Exponent', first: '2', second: '8', expression: '2 xʸ 8 = 256', result: '256' },
  { name: 'Square Root', firstLabel: 'Number', first: '144', expression: '√144 = 12', result: '12' },
  { name: 'Percentage', firstLabel: 'Base value', secondLabel: 'Percentage', first: '200', second: '15', expression: '15% of 200 = 30', result: '30' },
];

test.describe('Calculator acceptance', () => {
  for (const operation of operations) {
    test(`calculates ${operation.name} through the real backend`, async ({ page }) => {
      const calculator = new CalculatorPage(page);
      await calculator.open();
      const responsePromise = page.waitForResponse((response) => response.url().endsWith('/api/v1/calculations') && response.request().method() === 'POST');
      await calculator.selectOperation(operation.name);
      await calculator.enterOperand(operation.firstLabel, operation.first);
      if (operation.second && operation.secondLabel) await calculator.enterOperand(operation.secondLabel, operation.second);
      await calculator.calculate();
      expect((await responsePromise).status()).toBe(200);
      await expect(page.getByText(operation.expression)).toBeVisible();
      await expect(calculator.result()).toHaveText(operation.result);
    });
  }

  test('loads backend capabilities with all operations and the default selection', async ({ page }) => {
    const responsePromise = page.waitForResponse((response) => response.url().endsWith('/api/v1/operations'));
    const calculator = new CalculatorPage(page);
    await calculator.open();
    expect((await responsePromise).status()).toBe(200);
    await expect(page.getByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
    for (const operationName of operations.map((operation) => operation.name)) await expect(page.getByRole('button', { name: operationName })).toBeVisible();
  });

  test('keeps backend validation authoritative', async ({ request }) => {
    const response = await request.post('http://localhost:8080/api/v1/calculations', { data: { operationId: 'division', operands: [10, 0] } });
    expect(response.status()).toBe(422);
    await expect(response.json()).resolves.toMatchObject({ code: 'validation_failed' });
  });
});
