import type { ApiError, CalculationRequest, CalculationResponse, OperationDefinition, OperationManifest } from './types';
import { validateOperation } from '../model/validation';

type MockScenario = 'healthy' | 'manifestFailure' | 'serviceFailure';

const operationManifest: OperationManifest = {
  defaultOperationId: 'addition',
  operations: [
    { id: 'addition', name: 'Addition', symbol: '+', arity: 2, operands: [{ id: 'first', label: 'First number', placeholder: '0' }, { id: 'second', label: 'Second number', placeholder: '0' }], validations: [] },
    { id: 'subtraction', name: 'Subtraction', symbol: '−', arity: 2, operands: [{ id: 'first', label: 'First number', placeholder: '0' }, { id: 'second', label: 'Second number', placeholder: '0' }], validations: [] },
    { id: 'multiplication', name: 'Multiplication', symbol: '×', arity: 2, operands: [{ id: 'first', label: 'First number', placeholder: '0' }, { id: 'second', label: 'Second number', placeholder: '0' }], validations: [] },
    { id: 'division', name: 'Division', symbol: '÷', arity: 2, operands: [{ id: 'first', label: 'Dividend', placeholder: '0' }, { id: 'second', label: 'Divisor', placeholder: '0' }], validations: [{ message: 'The divisor must not be zero.', expression: { kind: 'comparison', operand: 'second', operator: 'notEqual', value: 0 } }] },
    { id: 'exponentiation', name: 'Exponentiation', symbol: 'xʸ', arity: 2, operands: [{ id: 'first', label: 'Base', placeholder: '0' }, { id: 'second', label: 'Exponent', placeholder: '0' }], validations: [] },
    { id: 'square-root', name: 'Square Root', symbol: '√', arity: 1, operands: [{ id: 'first', label: 'Number', placeholder: '0' }], validations: [{ message: 'The number must be zero or greater.', expression: { kind: 'comparison', operand: 'first', operator: 'greaterThanOrEqual', value: 0 } }] },
    { id: 'percentage', name: 'Percentage', symbol: '%', arity: 2, operands: [{ id: 'first', label: 'Base value', placeholder: '0' }, { id: 'second', label: 'Percentage', placeholder: '0', suffix: '%' }], validations: [{ message: 'Percentage must be between 0 and 100.', expression: { kind: 'allOf', expressions: [{ kind: 'comparison', operand: 'second', operator: 'greaterThanOrEqual', value: 0 }, { kind: 'comparison', operand: 'second', operator: 'lessThanOrEqual', value: 100 }] } }] },
  ],
};

let scenario: MockScenario = 'healthy';

const delay = () => new Promise((resolve) => window.setTimeout(resolve, 15));

export const configureMockScenario = (nextScenario: MockScenario): void => { scenario = nextScenario; };

export const getOperations = async (): Promise<OperationManifest> => {
  await delay();
  if (scenario === 'manifestFailure') throw apiError('manifest_unavailable', 'Calculator capabilities are unavailable.');
  return operationManifest;
};

export const calculate = async (request: CalculationRequest): Promise<CalculationResponse> => {
  await delay();
  if (scenario === 'serviceFailure') throw apiError('service_unavailable', 'The calculation service is temporarily unavailable.');
  const operation = operationManifest.operations.find((candidate) => candidate.id === request.operationId);
  if (!operation) throw apiError('unknown_operation', 'That calculation operation is not available.');
  const validationMessage = validateOperation(operation, request.operands);
  if (validationMessage) throw apiError('validation_failed', validationMessage);
  const result = execute(operation, request.operands);
  return { expression: formatExpression(operation, request.operands), result };
};

const apiError = (code: ApiError['code'], message: string): ApiError => ({ code, message });

const execute = (operation: OperationDefinition, operands: number[]): number => {
  const [firstOperand, secondOperand] = operands;
  switch (operation.id) {
    case 'addition': return firstOperand + secondOperand;
    case 'subtraction': return firstOperand - secondOperand;
    case 'multiplication': return firstOperand * secondOperand;
    case 'division': return firstOperand / secondOperand;
    case 'exponentiation': return firstOperand ** secondOperand;
    case 'square-root': return Math.sqrt(firstOperand);
    case 'percentage': return firstOperand * (secondOperand / 100);
    default: throw apiError('unknown_operation', 'That calculation operation is not available.');
  }
};

const formatExpression = (operation: OperationDefinition, operands: number[]): string => operation.arity === 1
  ? `${operation.symbol}${operands[0]}`
  : `${operands[0]} ${operation.symbol} ${operands[1]}`;
