package monitor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/monitor"
)

func TestCondition_Marshal(t *testing.T) {
	c := monitor.Condition{
		Variable: "result",
		Operator: "==",
		Value:    "OK",
		AndOr:    monitor.ConditionAnd,
	}

	data, err := json.Marshal(c)
	require.NoError(t, err)

	require.JSONEq(
		t,
		`{"type":"expression","variable":"result","operator":"==","value":"OK","andOr":"and"}`,
		string(data),
	)
}

func TestCondition_Unmarshal(t *testing.T) {
	data := []byte(`{"type":"expression","variable":"result","operator":"contains","value":"foo","andOr":"or"}`)

	var c monitor.Condition

	err := json.Unmarshal(data, &c)
	require.NoError(t, err)

	require.Equal(t, monitor.Condition{
		Variable: "result",
		Operator: "contains",
		Value:    "foo",
		AndOr:    monitor.ConditionOr,
	}, c)
}
