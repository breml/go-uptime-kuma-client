package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Postgres struct {
	Base
	PostgresDetails
}

func (p Postgres) Type() string {
	return p.PostgresDetails.Type()
}

func (p Postgres) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(p.Base, false), formatMonitor(p.PostgresDetails, true))
}

func (p *Postgres) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	details := PostgresDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	*p = Postgres{
		Base:            base,
		PostgresDetails: details,
	}

	return nil
}

func (p Postgres) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = p.ID
	raw["type"] = "postgres"
	raw["name"] = p.Name
	raw["description"] = p.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = p.PathName
	raw["parent"] = p.Parent
	raw["interval"] = p.Interval
	raw["retryInterval"] = p.RetryInterval
	raw["resendInterval"] = p.ResendInterval
	raw["maxretries"] = p.MaxRetries
	raw["upsideDown"] = p.UpsideDown
	raw["active"] = p.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range p.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}
	raw["notificationIDList"] = ids

	// Always override with current Postgres-specific field values.
	raw["databaseConnectionString"] = p.DatabaseConnectionString
	raw["databaseQuery"] = p.DatabaseQuery

	// Server expects these fields to be arrays and not null.
	raw["accepted_statuscodes"] = []string{}

	return json.Marshal(raw)
}

type PostgresDetails struct {
	DatabaseConnectionString string `json:"databaseConnectionString"`
	DatabaseQuery            string `json:"databaseQuery"`
}

func (p PostgresDetails) Type() string {
	return "postgres"
}
