package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

// TestLegacyTypeIsRewrittenOnMarshal pins the migration behavior for a
// notification stored with a notification type this package reported before it
// was aligned with the upstream provider names.
//
// The stored type is preserved as read from the server and is reported by
// notification.Base. Marshaling always emits the type of the details struct, so
// writing such a notification back repairs the stored type.
func TestLegacyTypeIsRewrittenOnMarshal(t *testing.T) {
	legacy := []byte(
		`{"id":1,"name":"My 46elks Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My 46elks Alert\",\"elksUsername\":\"username@example.com\",\"elksAuthToken\":\"auth_token_secret\",\"elksFromNumber\":\"1234\",\"elksToNumber\":\"0701234567\",\"type\":\"46elks\"}"}`,
	)

	elks := notification.FortySixElks{}

	err := json.Unmarshal(legacy, &elks)
	require.NoError(t, err)

	require.Equal(t, "46elks", elks.Base.Type(),
		"the stored notification type is preserved until the notification is written back")
	require.Equal(t, "Elks", elks.Type(),
		"the notification itself reports the upstream provider name")

	data, err := json.Marshal(elks)
	require.NoError(t, err)

	require.JSONEq(
		t,
		`{"active":true,"applyExisting":true,"elksAuthToken":"auth_token_secret","elksFromNumber":"1234","elksToNumber":"0701234567","elksUsername":"username@example.com","id":1,"isDefault":true,"name":"My 46elks Alert","type":"Elks","userId":1}`,
		string(data),
	)
}
