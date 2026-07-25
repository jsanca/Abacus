import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { FormEvent, KeyboardEvent } from 'react';
import { calculate, getOperations } from '../api/calculatorApi';
import type { ApiError, OperationManifest } from '../api/types';
import { validateOperation } from '../model/validation';
import './Calculator.css';

type CalculatorState = 'loading' | 'idle' | 'editing' | 'calculating' | 'resolved' | 'validation-error' | 'service-error';
type FocusTarget = 'first' | 'second';
type FocusIntent = { target: FocusTarget; requestId: number } | undefined;

const isErrorState = (state: CalculatorState): boolean => state === 'validation-error' || state === 'service-error';

/** Calculator renders manifest-driven calculator controls and calculated results. */
export function Calculator() {
  const [manifest, setManifest] = useState<OperationManifest>();
  const [selectedOperationId, setSelectedOperationId] = useState('');
  const [firstOperandText, setFirstOperandText] = useState('');
  const [secondOperandText, setSecondOperandText] = useState('');
  const [calculatorState, setCalculatorState] = useState<CalculatorState>('loading');
  const [message, setMessage] = useState('Loading calculator capabilities…');
  const [expression, setExpression] = useState('');
  const [result, setResult] = useState<number>();
  const [focusIntent, setFocusIntent] = useState<FocusIntent>();
  const firstOperandInput = useRef<HTMLInputElement>(null);
  const secondOperandInput = useRef<HTMLInputElement>(null);
  const focusRequestId = useRef(0);
  const selectedOperation = manifest?.operations.find((operation) => operation.id === selectedOperationId);
  const shortcutMap = useMemo<Record<string, string>>(() => {
    if (!manifest) return {};
    return Object.fromEntries(manifest.operations.map((op) => [op.shortcut, op.id]));
  }, [manifest]);

  const loadManifest = useCallback(async () => {
    setCalculatorState('loading');
    setMessage('Loading calculator capabilities…');
    try {
      const loadedManifest = await getOperations();
      setManifest(loadedManifest);
      setSelectedOperationId(loadedManifest.defaultOperationId);
      setCalculatorState('idle');
      setMessage('Choose an operation and enter values.');
    } catch (error) {
      const apiError = error as ApiError;
      setManifest(undefined);
      setSelectedOperationId('');
      setCalculatorState('service-error');
      setMessage(apiError.message);
    }
  }, []);

  useEffect(() => {
    let disposed = false;
    void Promise.resolve().then(() => {
      if (!disposed) return loadManifest();
    });
    return () => { disposed = true; };
  }, [loadManifest]);

  const requestFocus = useCallback((target: FocusTarget) => {
    focusRequestId.current += 1;
    setFocusIntent({ target, requestId: focusRequestId.current });
  }, []);

  const selectOperation = useCallback((operationId: string) => {
    const nextOperation = manifest?.operations.find((operation) => operation.id === operationId);
    if (!nextOperation) return;
    if (calculatorState === 'resolved' && result !== undefined) {
      setFirstOperandText(String(result));
      setSecondOperandText('');
      setExpression('');
      setResult(undefined);
    }
    setSelectedOperationId(operationId);
    setCalculatorState('editing');
    setMessage(`${nextOperation.name} selected.`);
    requestFocus(nextOperation.arity === 1 ? 'first' : 'second');
  }, [calculatorState, manifest, requestFocus, result]);

  const clearCalculator = useCallback(() => {
    setFirstOperandText('');
    setSecondOperandText('');
    setExpression('');
    setResult(undefined);
    setSelectedOperationId(manifest?.defaultOperationId ?? '');
    setCalculatorState('idle');
    setMessage('Calculator cleared.');
    requestFocus('first');
  }, [manifest, requestFocus]);

  useEffect(() => {
    const handleShortcut = (event: globalThis.KeyboardEvent) => {
      if (event.target instanceof HTMLInputElement || event.metaKey || event.ctrlKey || event.altKey) return;
      if (event.key === 'Escape') { clearCalculator(); return; }
      const operationId = shortcutMap[event.key];
      if (operationId) selectOperation(operationId);
    };
    window.addEventListener('keydown', handleShortcut);
    return () => window.removeEventListener('keydown', handleShortcut);
  }, [clearCalculator, selectOperation, shortcutMap]);

  useEffect(() => {
    if (!focusIntent) return;
    (focusIntent.target === 'first' ? firstOperandInput : secondOperandInput).current?.focus();
  }, [focusIntent]);

  const submitCalculation = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selectedOperation) return;
    const operandTexts = selectedOperation.arity === 1 ? [firstOperandText] : [firstOperandText, secondOperandText];
    const operands = operandTexts.map((operandText) => Number(operandText));
    const validationMessage = validateOperation(selectedOperation, operands);
    if (validationMessage) {
      setCalculatorState('validation-error');
      setMessage(validationMessage);
      return;
    }
    setCalculatorState('calculating');
    setMessage('Calculating…');
    try {
      const calculation = await calculate({ operationId: selectedOperation.id, operands });
      setExpression(`${calculation.expression} = ${calculation.result}`);
      setResult(calculation.result);
      setCalculatorState('resolved');
      setMessage('Calculation complete.');
    } catch (error) {
      const apiError = error as ApiError;
      setCalculatorState(apiError.code === 'validation_failed' ? 'validation-error' : 'service-error');
      setMessage(apiError.message);
    }
  };

  const handleInputKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Escape') { event.preventDefault(); clearCalculator(); }
  };

  if (!manifest || !selectedOperation) return <main className="application-shell"><section className="calculator-card loading-card" aria-live="polite"><p role={isErrorState(calculatorState) ? 'alert' : 'status'}>{message}</p>{calculatorState === 'service-error' && <button className="retry-button" type="button" onClick={() => void loadManifest()}>Retry</button>}</section></main>;

  return <main className="application-shell">
    <section className="calculator-card" aria-label="Abacus calculator">
      <header className="calculator-header"><span className="brand">Abacus</span><span className="shortcut-hint">Enter to calculate · Esc to clear</span></header>
      <section className="result-panel" aria-live="polite" aria-atomic="true">
        <p className="expression">{expression || 'Ready when you are'}</p>
        <output className="result-value">{result ?? '—'}</output>
      </section>
      <form className="calculator-form" onSubmit={submitCalculation}>
        <fieldset className="operator-group"><legend>Operation</legend><div className="operator-grid">
          {manifest.operations.map((operation) => <button className="operator-button" key={operation.id} type="button" aria-pressed={selectedOperation.id === operation.id} aria-label={operation.name} onClick={() => selectOperation(operation.id)}>{operation.symbol}</button>)}
        </div></fieldset>
        <div className="operand-fields">
          <label className="operand-field"><span>{selectedOperation.operands[0].label}</span><input ref={firstOperandInput} value={firstOperandText} onChange={(event) => { setFirstOperandText(event.target.value); setCalculatorState('editing'); }} onKeyDown={handleInputKeyDown} type="text" inputMode="decimal" placeholder={selectedOperation.operands[0].placeholder} aria-describedby="calculator-message" /></label>
          {selectedOperation.arity === 2 && <label className="operand-field"><span>{selectedOperation.operands[1].label}</span><span className="input-with-suffix"><input ref={secondOperandInput} value={secondOperandText} onChange={(event) => { setSecondOperandText(event.target.value); setCalculatorState('editing'); }} onKeyDown={handleInputKeyDown} type="text" inputMode="decimal" placeholder={selectedOperation.operands[1].placeholder} aria-describedby="calculator-message" />{selectedOperation.operands[1].suffix && <span aria-hidden="true">{selectedOperation.operands[1].suffix}</span>}</span></label>}
        </div>
        <p id="calculator-message" className={isErrorState(calculatorState) ? 'status-message error-message' : 'status-message'} role={isErrorState(calculatorState) ? 'alert' : 'status'}>{message}</p>
        <div className="action-row"><button className="clear-button" type="button" onClick={clearCalculator}>Clear</button><button className="calculate-button" type="submit" disabled={calculatorState === 'calculating'}>{calculatorState === 'calculating' ? 'Calculating…' : 'Calculate'}</button></div>
      </form>
    </section>
  </main>;
}
