package monitor

import (
	"encoding/json"
	"fmt"
)

// ConditionOperator represents the logical operator that joins a condition with
// the previous one in the chain.
type ConditionOperator string

const (
	// ConditionAnd combines a condition with the previous one using a logical AND.
	ConditionAnd ConditionOperator = "and"
	// ConditionOr combines a condition with the previous one using a logical OR.
	ConditionOr ConditionOperator = "or"
)

// conditionTypeExpression is the wire-format value upstream Uptime Kuma uses for
// flat expression-type conditions. Group-type conditions (with nested children)
// are not modeled by this client.
const conditionTypeExpression = "expression"

// Condition represents a single assertion clause that an Uptime Kuma monitor
// evaluates against its parsed result (query result, MQTT payload, JSON-query
// value, etc.). Conditions are chained together using the AndOr field.
//
// The set of allowed Variable and Operator values is monitor-type specific and
// validated server-side. See the upstream `EditMonitor.vue` for the available
// variables per monitor type.
type Condition struct {
	// Variable is the name of the field to test against (monitor-type specific).
	Variable string `json:"variable"`
	// Operator is the comparison operator (e.g. "==", "!=", "<", ">", "contains").
	Operator string `json:"operator"`
	// Value is the value to compare against.
	Value string `json:"value"`
	// AndOr chains this condition with the previous one ("and" or "or").
	AndOr ConditionOperator `json:"andOr"`
}

// MarshalJSON marshals a Condition to JSON, adding the wire-format
// `type: "expression"` field expected by the upstream evaluator.
func (c Condition) MarshalJSON() ([]byte, error) {
	raw := struct {
		Type     string            `json:"type"`
		Variable string            `json:"variable"`
		Operator string            `json:"operator"`
		Value    string            `json:"value"`
		AndOr    ConditionOperator `json:"andOr"`
	}{
		Type:     conditionTypeExpression,
		Variable: c.Variable,
		Operator: c.Operator,
		Value:    c.Value,
		AndOr:    c.AndOr,
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal condition: %w", err)
	}

	return data, nil
}

// conditionsForWire normalizes a Conditions slice so that nil is encoded as an
// empty JSON array, matching the upstream wire-format expectation.
func conditionsForWire(conditions []Condition) any {
	if conditions == nil {
		return []Condition{}
	}

	return conditions
}
