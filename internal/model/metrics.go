package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func (m *Metrics) UnmarshalJSON(data []byte) error {
	type MetricsAlias Metrics
	aux := &struct {
		*MetricsAlias
		Delta any `json:"delta,omitempty"`
		Value any `json:"value,omitempty"`
	}{
		MetricsAlias: (*MetricsAlias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal metrics: %w", err)
	}

	if aux.Delta != nil {
		switch v := aux.Delta.(type) {
		case float64:
			delta := int64(v)
			m.Delta = &delta
		case string:
			delta, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing delta: %w", err)
			}
			m.Delta = &delta
		default:
			return fmt.Errorf("unexpected type for delta: %T", aux.Delta)
		}
	}

	if aux.Value != nil {
		switch v := aux.Value.(type) {
		case float64:
			m.Value = &v
		case string:
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("error parsing value: %w", err)
			}
			m.Value = &value
		default:
			return fmt.Errorf("unexpected type for value: %T", aux.Value)
		}
	}

	return nil
}
