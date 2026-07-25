package calculator

import "fmt"

// Registry is an immutable, ordered collection of calculator operations.
// It is safe for concurrent reads without synchronization because no writes
// occur after successful construction; no mutex is required.
type Registry struct {
	// operationsByID enables O(1) lookup by stable operation ID.
	operationsByID map[string]Operation
	// operations preserves the deterministic product-defined ordering required
	// for manifest projection. Go map iteration is nondeterministic, so ordering
	// is maintained separately in this slice.
	operations []Operation
	// defaultID is the operation ID selected on initial client load.
	defaultID string
}

// NewRegistry constructs an immutable Registry from the given operations.
// It validates all configuration invariants at construction time and returns
// a descriptive error on the first violation.
//
// The constructor copies the operations slice and each operation's operand slice
// so that later mutations of the caller's input cannot affect registry state.
func NewRegistry(defaultOperationID string, operations []Operation) (*Registry, error) {
	if len(operations) == 0 {
		return nil, fmt.Errorf("registry requires at least one operation")
	}
	if defaultOperationID == "" {
		return nil, fmt.Errorf("registry requires a non-empty default operation ID")
	}

	operationsByID := make(map[string]Operation, len(operations))
	orderedCopy := make([]Operation, len(operations))

	for index, operation := range operations {
		if err := validateRegistryOperation(operation); err != nil {
			return nil, fmt.Errorf("operation at index %d: %w", index, err)
		}
		if _, exists := operationsByID[operation.Definition.ID]; exists {
			return nil, fmt.Errorf("duplicate operation ID %q", operation.Definition.ID)
		}
		// Deep-copy to prevent the caller's later mutations from aliasing registry internals.
		owned := copyOperation(operation)
		operationsByID[owned.Definition.ID] = owned
		orderedCopy[index] = owned
	}

	if _, exists := operationsByID[defaultOperationID]; !exists {
		return nil, fmt.Errorf("default operation ID %q is not registered", defaultOperationID)
	}

	return &Registry{
		operationsByID: operationsByID,
		operations:     orderedCopy,
		defaultID:      defaultOperationID,
	}, nil
}

// Find returns the Operation with the given ID and whether it was found.
// The returned Operation is a deep copy; modifying it does not affect the registry.
//
// Find is preferred over Lookup here because it reads naturally at call sites
// (registry.Find("addition")) and is consistent with the comma-ok idiom for
// lookup-with-existence semantics.
func (registry *Registry) Find(operationID string) (Operation, bool) {
	operation, exists := registry.operationsByID[operationID]
	if !exists {
		return Operation{}, false
	}
	return copyOperation(operation), true
}

// Operations returns all registered operations in deterministic product order.
// The returned slice is a defensive copy; modifying it does not affect the registry.
func (registry *Registry) Operations() []Operation {
	result := make([]Operation, len(registry.operations))
	for index, operation := range registry.operations {
		result[index] = copyOperation(operation)
	}
	return result
}

// DefaultOperationID returns the ID of the operation selected on initial client load.
func (registry *Registry) DefaultOperationID() string {
	return registry.defaultID
}

// validateRegistryOperation checks that a single operation satisfies all registry invariants.
func validateRegistryOperation(operation Operation) error {
	definition := operation.Definition

	if definition.ID == "" {
		return fmt.Errorf("operation ID must not be empty")
	}
	if definition.Name == "" {
		return fmt.Errorf("operation %q: name must not be empty", definition.ID)
	}
	if definition.Symbol == "" {
		return fmt.Errorf("operation %q: symbol must not be empty", definition.ID)
	}
	if definition.Shortcut == "" {
		return fmt.Errorf("operation %q: shortcut must not be empty", definition.ID)
	}
	if definition.Arity != Unary && definition.Arity != Binary {
		return fmt.Errorf("operation %q: unsupported arity %d; only 1 and 2 are supported", definition.ID, definition.Arity)
	}
	if len(definition.Operands) != int(definition.Arity) {
		return fmt.Errorf("operation %q: arity %d but %d operand(s) defined", definition.ID, definition.Arity, len(definition.Operands))
	}

	seenOperandIDs := make(map[string]bool, len(definition.Operands))
	for _, operand := range definition.Operands {
		if operand.ID == "" {
			return fmt.Errorf("operation %q: operand has empty ID", definition.ID)
		}
		if seenOperandIDs[operand.ID] {
			return fmt.Errorf("operation %q: duplicate operand ID %q", definition.ID, operand.ID)
		}
		seenOperandIDs[operand.ID] = true
	}

	if operation.Execute == nil {
		return fmt.Errorf("operation %q: executor must not be nil", definition.ID)
	}

	return nil
}
