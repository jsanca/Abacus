import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { Calculator } from '../components/Calculator';
import { configureMockScenario } from '../api/mockCalculatorApi';

const renderReadyCalculator = async () => {
  render(<Calculator />);
  await screen.findByRole('button', { name: 'Addition' });
};

describe('Calculator', () => {
  beforeEach(() => configureMockScenario('healthy'));
  afterEach(cleanup);

  it('renders every operation supplied by the manifest', async () => {
    await renderReadyCalculator();
    for (const operationName of ['Addition', 'Subtraction', 'Multiplication', 'Division', 'Exponentiation', 'Square Root', 'Percentage']) {
      expect(screen.getByRole('button', { name: operationName })).toBeVisible();
    }
  });

  it('shows a loading state before capabilities arrive', () => {
    render(<Calculator />);
    expect(screen.getByText('Loading calculator capabilities…')).toBeVisible();
  });

  it('retries manifest loading after a recoverable service error', async () => {
    const user = userEvent.setup();
    configureMockScenario('manifestFailure');
    render(<Calculator />);
    expect(await screen.findByRole('alert')).toHaveTextContent('capabilities are unavailable');
    configureMockScenario('healthy');
    await user.click(screen.getByRole('button', { name: 'Retry' }));
    expect(screen.getByText('Loading calculator capabilities…')).toBeVisible();
    expect(await screen.findByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
  });

  it('switches to a unary square-root layout with one click', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.click(screen.getByRole('button', { name: 'Square Root' }));
    expect(screen.getByRole('button', { name: 'Square Root' })).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByRole('textbox', { name: 'Number' })).toBeVisible();
    expect(screen.queryByRole('textbox', { name: 'Second number' })).not.toBeInTheDocument();
  });

  it('uses manifest-provided percentage suffix and range validation', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.click(screen.getByRole('button', { name: 'Percentage' }));
    const percentageInput = screen.getByRole('textbox', { name: 'Percentage' });
    expect(percentageInput.parentElement).toHaveTextContent('%');
    await user.type(screen.getByRole('textbox', { name: 'Base value' }), '200');
    await user.type(percentageInput, '101');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Percentage must be between 0 and 100.');
  });

  it('calculates through the mocked API and preserves the expression', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.click(screen.getByRole('button', { name: 'Division' }));
    await user.type(screen.getByRole('textbox', { name: 'Dividend' }), '144');
    await user.type(screen.getByRole('textbox', { name: 'Divisor' }), '12');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(await screen.findByText('144 ÷ 12 = 12')).toBeVisible();
    expect(screen.getByText('12', { selector: 'output' })).toBeVisible();
  });

  it('applies division validation without calculator-specific UI logic', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.click(screen.getByRole('button', { name: 'Division' }));
    await user.type(screen.getByRole('textbox', { name: 'Dividend' }), '10');
    await user.type(screen.getByRole('textbox', { name: 'Divisor' }), '0');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('The divisor must not be zero.');
  });

  it('continues from the previous result after selecting a new operation', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.click(screen.getByRole('button', { name: 'Division' }));
    await user.type(screen.getByRole('textbox', { name: 'Dividend' }), '144');
    await user.type(screen.getByRole('textbox', { name: 'Divisor' }), '12');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    await screen.findByText('144 ÷ 12 = 12');
    await user.click(screen.getByRole('button', { name: 'Addition' }));
    expect(screen.getByRole('textbox', { name: 'First number' })).toHaveValue('12');
    const secondNumber = screen.getByRole('textbox', { name: 'Second number' });
    expect(secondNumber).toHaveValue('');
    await waitFor(() => expect(secondNumber).toHaveFocus());
  });

  it('clears the current calculation and returns focus to the first operand', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    const firstNumber = screen.getByRole('textbox', { name: 'First number' });
    await user.type(firstNumber, '9');
    await user.click(screen.getByRole('button', { name: 'Clear' }));
    await waitFor(() => expect(firstNumber).toHaveFocus());
    expect(firstNumber).toHaveValue('');
  });

  it('restores the default binary operation and focus when clearing Square Root', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.click(screen.getByRole('button', { name: 'Square Root' }));
    await user.type(screen.getByRole('textbox', { name: 'Number' }), '25');
    await user.click(screen.getByRole('button', { name: 'Clear' }));
    const firstNumber = screen.getByRole('textbox', { name: 'First number' });
    expect(screen.getByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByRole('textbox', { name: 'Second number' })).toBeVisible();
    await waitFor(() => expect(firstNumber).toHaveFocus());
  });

  it('keeps keyboard shortcuts effective across state transitions without duplicate changes', async () => {
    await renderReadyCalculator();
    fireEvent.keyDown(window, { key: 'r' });
    expect(await screen.findByRole('button', { name: 'Square Root' })).toHaveAttribute('aria-pressed', 'true');
    fireEvent.keyDown(window, { key: '+' });
    expect(await screen.findByRole('button', { name: 'Addition' })).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getAllByRole('textbox')).toHaveLength(2);
  });

  it('shows a service error when the mocked calculation service fails', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    configureMockScenario('serviceFailure');
    await user.type(screen.getByRole('textbox', { name: 'First number' }), '1');
    await user.type(screen.getByRole('textbox', { name: 'Second number' }), '2');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('temporarily unavailable');
  });

  it('calculates when Enter is pressed in an operand field', async () => {
    const user = userEvent.setup();
    await renderReadyCalculator();
    await user.type(screen.getByRole('textbox', { name: 'First number' }), '2');
    await user.type(screen.getByRole('textbox', { name: 'Second number' }), '3{Enter}');
    expect(await screen.findByText('2 + 3 = 5')).toBeVisible();
  });
});
