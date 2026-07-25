package calculator

import "encoding/json"

// ManifestVersion is the stable schema version included in every manifest response.
const ManifestVersion = "1"

// Manifest is the frontend-safe projection of the operation registry.
// It is produced by Registry.Manifest and safe to serialize as JSON for the client.
//
// Contract note: the TypeScript OperationManifest does not currently define a
// version field; the extra field is ignored by the existing client and is
// forward-compatible. The shortcut field in each ManifestOperation is similarly
// extra. Both will be aligned with the client in a future integration task.
type Manifest struct {
	// Version identifies the manifest schema for future compatibility checks.
	Version string `json:"version"`
	// DefaultOperationID is the ID of the operation selected on initial load.
	DefaultOperationID string `json:"defaultOperationId"`
	// Operations is the ordered list of available calculator operations.
	Operations []ManifestOperation `json:"operations"`
}

// ManifestOperation is the frontend-safe projection of a single Operation.
// It excludes the executor function and any backend-only runtime details.
type ManifestOperation struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Symbol   string            `json:"symbol"`
	Shortcut string            `json:"shortcut"`
	Arity    int               `json:"arity"`
	Operands []ManifestOperand `json:"operands"`
	// Validations is always an empty JSON array in this manifest version.
	// The declarative validation contract (ABA-006) will populate this field.
	// An empty array is used rather than null to match the TypeScript
	// ValidationDefinition[] contract, which the client reads as an iterable.
	Validations json.RawMessage `json:"validations"`
}

// ManifestOperand is the frontend-safe projection of a single OperandDefinition.
type ManifestOperand struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Placeholder string `json:"placeholder"`
	Suffix      string `json:"suffix,omitempty"`
}

// Manifest generates a frontend-safe projection of the registry.
// The returned value is independent of registry internals; mutating it
// does not affect registered operations.
func (registry *Registry) Manifest() Manifest {
	manifestOperations := make([]ManifestOperation, len(registry.operations))
	for index, operation := range registry.operations {
		manifestOperations[index] = projectOperation(operation)
	}
	return Manifest{
		Version:            ManifestVersion,
		DefaultOperationID: registry.defaultID,
		Operations:         manifestOperations,
	}
}

// projectOperation produces a ManifestOperation from a registered Operation.
// It allocates new slices so the result does not alias registry-owned storage.
func projectOperation(operation Operation) ManifestOperation {
	definition := operation.Definition

	manifestOperands := make([]ManifestOperand, len(definition.Operands))
	for index, operand := range definition.Operands {
		manifestOperands[index] = ManifestOperand(operand)
	}

	return ManifestOperation{
		ID:          definition.ID,
		Name:        definition.Name,
		Symbol:      definition.Symbol,
		Shortcut:    definition.Shortcut,
		Arity:       int(definition.Arity),
		Operands:    manifestOperands,
		Validations: json.RawMessage("[]"),
	}
}
