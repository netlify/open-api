package models

import "fmt"

// ValidateCrossFieldConstraints checks cross-field invariants that cannot be
// expressed in the OpenAPI 2.0 schema and are therefore not covered by the
// generated Validate method.  Callers should invoke this after Validate.
//
// Rules enforced:
//   - When both min_cu and max_cu are provided, max_cu >= min_cu.
//   - When both min_cu and max_cu are provided, max_cu - min_cu <= 8.0.
func (m *DatabaseComputeSettingsRequest) ValidateCrossFieldConstraints() error {
	if m.MinCu == nil || m.MaxCu == nil {
		return nil
	}

	if *m.MaxCu < *m.MinCu {
		return fmt.Errorf("max_cu (%g) must be greater than or equal to min_cu (%g)", *m.MaxCu, *m.MinCu)
	}

	if *m.MaxCu-*m.MinCu > 8.0 {
		return fmt.Errorf("max_cu - min_cu must not exceed 8.0 (got %g)", *m.MaxCu-*m.MinCu)
	}

	return nil
}
