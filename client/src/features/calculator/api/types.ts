export type ComparisonOperator = 'equal' | 'notEqual' | 'greaterThan' | 'greaterThanOrEqual' | 'lessThan' | 'lessThanOrEqual';

export type ValidationExpression =
  | { kind: 'comparison'; operand: 'first' | 'second'; operator: ComparisonOperator; value: number }
  | { kind: 'allOf'; expressions: ValidationExpression[] };

export interface ValidationDefinition {
  message: string;
  expression: ValidationExpression;
}

export interface OperandDefinition {
  id: 'first' | 'second';
  label: string;
  placeholder: string;
  suffix?: string;
}

export interface OperationDefinition {
  id: string;
  name: string;
  symbol: string;
  shortcut: string;
  arity: 1 | 2;
  operands: OperandDefinition[];
  validations: ValidationDefinition[];
}

export interface OperationManifest {
  version: string;
  operations: OperationDefinition[];
  defaultOperationId: string;
}

export interface CalculationRequest {
  operationId: string;
  operands: number[];
}

export interface CalculationResponse {
  expression: string;
  result: number;
}

export interface ApiError {
  code: 'manifest_unavailable' | 'unknown_operation' | 'validation_failed' | 'calculation_failed' | 'service_unavailable';
  message: string;
}
