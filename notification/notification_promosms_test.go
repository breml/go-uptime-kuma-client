package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationPromoSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.PromoSMS
		wantJSON string
	}{
		{
			name: "success with long SMS",
			data: []byte(`{"id":1,"name":"Test PromoSMS","active":true,"userId":42,"isDefault":false,"config":"{\"type\":\"promosms\",\"promosmsLogin\":\"user@example.com\",\"promosmsPassword\":\"password123\",\"promosmsPhoneNumber\":\"+48123456789\",\"promosmsSenderName\":\"UptimeKuma\",\"promosmsSMSType\":\"1\",\"promosmsAllowLongSMS\":true}"}`),

			want: notification.PromoSMS{
				Base: notification.Base{
					ID:        1,
					Name:      "Test PromoSMS",
					IsActive:  true,
					UserID:    42,
					IsDefault: false,
				},
				PromoSMSDetails: notification.PromoSMSDetails{
					Login:        "user@example.com",
					Password:     "password123",
					PhoneNumber:  "+48123456789",
					SenderName:   "UptimeKuma",
					SMSType:      "1",
					AllowLongSMS: true,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":1,"isDefault":false,"name":"Test PromoSMS","promosmsAllowLongSMS":true,"promosmsLogin":"user@example.com","promosmsPassword":"password123","promosmsPhoneNumber":"+48123456789","promosmsSenderName":"UptimeKuma","promosmsSMSType":"1","type":"promosms","userId":42}`,
		},
		{
			name: "success without long SMS",
			data: []byte(`{"id":2,"name":"Test PromoSMS Short","active":true,"userId":42,"isDefault":true,"config":"{\"type\":\"promosms\",\"promosmsLogin\":\"admin@example.com\",\"promosmsPassword\":\"secret\",\"promosmsPhoneNumber\":\"+48987654321\",\"promosmsSenderName\":\"Monitoring\",\"promosmsSMSType\":\"0\",\"promosmsAllowLongSMS\":false}"}`),

			want: notification.PromoSMS{
				Base: notification.Base{
					ID:            2,
					Name:          "Test PromoSMS Short",
					IsActive:      true,
					UserID:        42,
					IsDefault:     true,
					ApplyExisting: false,
				},
				PromoSMSDetails: notification.PromoSMSDetails{
					Login:        "admin@example.com",
					Password:     "secret",
					PhoneNumber:  "+48987654321",
					SenderName:   "Monitoring",
					SMSType:      "0",
					AllowLongSMS: false,
				},
			},
			wantJSON: `{"active":true,"applyExisting":false,"id":2,"isDefault":true,"name":"Test PromoSMS Short","promosmsAllowLongSMS":false,"promosmsLogin":"admin@example.com","promosmsPassword":"secret","promosmsPhoneNumber":"+48987654321","promosmsSenderName":"Monitoring","promosmsSMSType":"0","type":"promosms","userId":42}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":3,"name":"Test PromoSMS Minimal","active":false,"userId":10,"isDefault":false,"config":"{\"type\":\"promosms\",\"promosmsLogin\":\"\",\"promosmsPassword\":\"\",\"promosmsPhoneNumber\":\"\",\"promosmsSenderName\":\"\",\"promosmsSMSType\":\"\",\"promosmsAllowLongSMS\":false}"}`),

			want: notification.PromoSMS{
				Base: notification.Base{
					ID:            3,
					Name:          "Test PromoSMS Minimal",
					IsActive:      false,
					UserID:        10,
					IsDefault:     false,
					ApplyExisting: false,
				},
				PromoSMSDetails: notification.PromoSMSDetails{
					Login:        "",
					Password:     "",
					PhoneNumber:  "",
					SenderName:   "",
					SMSType:      "",
					AllowLongSMS: false,
				},
			},
			wantJSON: `{"active":false,"applyExisting":false,"id":3,"isDefault":false,"name":"Test PromoSMS Minimal","promosmsAllowLongSMS":false,"promosmsLogin":"","promosmsPassword":"","promosmsPhoneNumber":"","promosmsSenderName":"","promosmsSMSType":"","type":"promosms","userId":10}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			promosms := notification.PromoSMS{}

			err := json.Unmarshal(tc.data, &promosms)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, promosms)

			data, err := json.Marshal(promosms)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
