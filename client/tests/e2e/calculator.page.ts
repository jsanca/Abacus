import { expect, type Locator, type Page } from '@playwright/test';

/** CalculatorPage exposes user-facing calculator controls for acceptance tests. */
export class CalculatorPage {
  constructor(private readonly page: Page) {}

  async open(): Promise<void> {
    await this.page.goto('/');
    await this.waitUntilReady();
  }

  async waitUntilReady(): Promise<void> {
    await expect(this.page.getByRole('button', { name: 'Addition' })).toBeVisible();
  }

  async selectOperation(name: string): Promise<void> {
    await this.page.getByRole('button', { name }).click();
  }

  async enterOperand(label: string, value: string): Promise<void> {
    await this.page.getByRole('textbox', { name: label }).fill(value);
  }

  async calculate(): Promise<void> {
    await this.page.getByRole('button', { name: 'Calculate' }).click();
  }

  async clear(): Promise<void> {
    await this.page.getByRole('button', { name: 'Clear' }).click();
  }

  result(): Locator {
    return this.page.getByRole('status', { name: 'Calculation result' });
  }

  status(): Locator {
    return this.page.getByRole('status', { name: 'Calculator status' });
  }

  error(): Locator {
    return this.page.getByRole('alert');
  }
}
