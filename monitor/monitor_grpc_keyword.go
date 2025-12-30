package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GrpcKeyword ...
type GrpcKeyword struct {
	Base
	GrpcKeywordDetails
}

// Type ...
func (g GrpcKeyword) Type() string {
	return g.GrpcKeywordDetails.Type()
}

// String ...
func (g GrpcKeyword) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(g.Base, false), formatMonitor(g.GrpcKeywordDetails, true))
}

// UnmarshalJSON ...
func (g *GrpcKeyword) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := GrpcKeywordDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*g = GrpcKeyword{
		Base:               base,
		GrpcKeywordDetails: details,
	}

	return nil
}

// MarshalJSON ...
func (g GrpcKeyword) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = g.ID
	raw["type"] = "grpc-keyword"
	raw["name"] = g.Name
	raw["description"] = g.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = g.PathName
	raw["parent"] = g.Parent
	raw["interval"] = g.Interval
	raw["retryInterval"] = g.RetryInterval
	raw["resendInterval"] = g.ResendInterval
	raw["maxretries"] = g.MaxRetries
	raw["upsideDown"] = g.UpsideDown
	raw["active"] = g.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range g.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current gRPC-specific field values.
	raw["grpcUrl"] = g.GrpcURL
	raw["grpcProtobuf"] = g.GrpcProtobuf
	raw["grpcServiceName"] = g.GrpcServiceName
	raw["grpcMethod"] = g.GrpcMethod
	raw["grpcEnableTls"] = g.GrpcEnableTLS
	raw["grpcBody"] = g.GrpcBody
	raw["keyword"] = g.Keyword
	raw["invertKeyword"] = g.InvertKeyword

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	// Uptime Kuma v2 requires conditions field (empty array by default)
	raw["conditions"] = []any{}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return data, nil
}

// GrpcKeywordDetails ...
type GrpcKeywordDetails struct {
	GrpcURL         string `json:"grpcUrl"`
	GrpcProtobuf    string `json:"grpcProtobuf"`
	GrpcServiceName string `json:"grpcServiceName"`
	GrpcMethod      string `json:"grpcMethod"`
	GrpcEnableTLS   bool   `json:"grpcEnableTls"`
	GrpcBody        string `json:"grpcBody"`
	Keyword         string `json:"keyword"`
	InvertKeyword   bool   `json:"invertKeyword"`
}

// Type ...
func (GrpcKeywordDetails) Type() string {
	return "grpc-keyword"
}
