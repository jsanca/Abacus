import type { OperationDefinition, ValidationExpression } from '../api/types';

export const validateOperation = (operation: OperationDefinition, operands: number[]): string | undefined => {
  if (operands.length !== operation.arity || operands.some((operand) => !Number.isFinite(operand))) {
    return 'Enter a valid number for each required operand.';
  }
  return operation.validations.find((validation) => !evaluate(validation.expression, operands))?.message;
};

const evaluate = (expression: ValidationExpression, operands: number[]): boolean => {
  if (expression.kind === 'allOf') return expression.expressions.every((child) => evaluate(child, operands));
  const operand = operands[expression.operand === 'first' ? 0 : 1];
  switch (expression.operator) {
    case 'equal': return operand === expression.value;
    case 'notEqual': return operand !== expression.value;
    case 'greaterThan': return operand > expression.value;
    case 'greaterThanOrEqual': return operand >= expression.value;
    case 'lessThan': return operand < expression.value;
    case 'lessThanOrEqual': return operand <= expression.value;
  }
};
