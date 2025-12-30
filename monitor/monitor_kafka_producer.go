package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// KafkaProducer represents a Kafka Producer monitor for testing Kafka connectivity.
type KafkaProducer struct {
	Base
	KafkaProducerDetails
}

// Type returns the monitor type.
func (k KafkaProducer) Type() string {
	return k.KafkaProducerDetails.Type()
}

// String returns a string representation of the KafkaProducer monitor.
func (k KafkaProducer) String() string {
	return fmt.Sprintf("%s, %s", formatMonitor(k.Base, false), formatMonitor(k.KafkaProducerDetails, true))
}

// UnmarshalJSON unmarshals a KafkaProducer monitor from JSON data.
func (k *KafkaProducer) UnmarshalJSON(data []byte) error {
	base := Base{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	details := KafkaProducerDetails{}
	err = json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	*k = KafkaProducer{
		Base:                 base,
		KafkaProducerDetails: details,
	}

	return nil
}

// MarshalJSON marshals a KafkaProducer monitor to JSON data.
func (k KafkaProducer) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	raw["id"] = k.ID
	raw["type"] = "kafka-producer"
	raw["name"] = k.Name
	raw["description"] = k.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = k.PathName
	raw["parent"] = k.Parent
	raw["interval"] = k.Interval
	raw["retryInterval"] = k.RetryInterval
	raw["resendInterval"] = k.ResendInterval
	raw["maxretries"] = k.MaxRetries
	raw["upsideDown"] = k.UpsideDown
	raw["active"] = k.IsActive

	// Update notification IDs.
	ids := map[string]bool{}
	for _, id := range k.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	// Always override with current KafkaProducer-specific field values.
	raw["kafkaProducerBrokers"] = k.Brokers
	raw["kafkaProducerTopic"] = k.Topic
	raw["kafkaProducerMessage"] = k.Message
	raw["kafkaProducerSsl"] = k.SSL
	raw["kafkaProducerAllowAutoTopicCreation"] = k.AllowAutoTopicCreation
	raw["kafkaProducerSaslOptions"] = k.SASLOptions

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

// KafkaProducerDetails contains Kafka Producer-specific monitor configuration.
type KafkaProducerDetails struct {
	// Brokers is a JSON array of broker addresses.
	Brokers []string `json:"kafkaProducerBrokers"`
	// Topic is the topic to publish to.
	Topic string `json:"kafkaProducerTopic"`
	// Message is the test message to send.
	Message string `json:"kafkaProducerMessage"`
	// SSL enables SSL/TLS connection.
	SSL bool `json:"kafkaProducerSsl"`
	// AllowAutoTopicCreation allows automatic topic creation.
	AllowAutoTopicCreation bool `json:"kafkaProducerAllowAutoTopicCreation"`
	// SASLOptions is an optional map containing string with SASL configuration.
	SASLOptions *map[string]any `json:"kafkaProducerSaslOptions"`
}

// Type returns the monitor type.
func (KafkaProducerDetails) Type() string {
	return "kafka-producer"
}
