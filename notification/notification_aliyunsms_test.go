package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/notification"
)

func TestNotificationAliyunSMS_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte

		want     notification.AliyunSMS
		wantJSON string
	}{
		{
			name: "success",
			data: []byte(`{"id":1,"name":"My Aliyun SMS Alert","active":true,"userId":1,"isDefault":true,"config":"{\"applyExisting\":true,\"isDefault\":true,\"name\":\"My Aliyun SMS Alert\",\"accessKeyId\":\"AKIAIOSFODNN7EXAMPLE\",\"secretAccessKey\":\"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\",\"phonenumber\":\"8613800000001\",\"signName\":\"Uptime Kuma\",\"templateCode\":\"SMS_1234567890\",\"type\":\"AliyunSMS\"}"}`),

			want: notification.AliyunSMS{
				Base: notification.Base{
					ID:            1,
					Name:          "My Aliyun SMS Alert",
					IsActive:      true,
					UserID:        1,
					IsDefault:     true,
					ApplyExisting: true,
				},
				AliyunSMSDetails: notification.AliyunSMSDetails{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					PhoneNumber:     "8613800000001",
					SignName:        "Uptime Kuma",
					TemplateCode:    "SMS_1234567890",
				},
			},
			wantJSON: `{"accessKeyId":"AKIAIOSFODNN7EXAMPLE","active":true,"applyExisting":true,"id":1,"isDefault":true,"name":"My Aliyun SMS Alert","phonenumber":"8613800000001","secretAccessKey":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY","signName":"Uptime Kuma","templateCode":"SMS_1234567890","type":"AliyunSMS","userId":1}`,
		},
		{
			name: "minimal",
			data: []byte(`{"id":2,"name":"Simple Aliyun SMS","active":true,"userId":1,"isDefault":false,"config":"{\"applyExisting\":false,\"isDefault\":false,\"name\":\"Simple Aliyun SMS\",\"accessKeyId\":\"AKIA123\",\"secretAccessKey\":\"secret123\",\"phonenumber\":\"8613800000002\",\"signName\":\"Alert\",\"templateCode\":\"SMS_0000000001\",\"type\":\"AliyunSMS\"}"}`),

			want: notification.AliyunSMS{
				Base: notification.Base{
					ID:            2,
					Name:          "Simple Aliyun SMS",
					IsActive:      true,
					UserID:        1,
					IsDefault:     false,
					ApplyExisting: false,
				},
				AliyunSMSDetails: notification.AliyunSMSDetails{
					AccessKeyID:     "AKIA123",
					SecretAccessKey: "secret123",
					PhoneNumber:     "8613800000002",
					SignName:        "Alert",
					TemplateCode:    "SMS_0000000001",
				},
			},
			wantJSON: `{"accessKeyId":"AKIA123","active":true,"applyExisting":false,"id":2,"isDefault":false,"name":"Simple Aliyun SMS","phonenumber":"8613800000002","secretAccessKey":"secret123","signName":"Alert","templateCode":"SMS_0000000001","type":"AliyunSMS","userId":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			aliyunsms := notification.AliyunSMS{}

			err := json.Unmarshal(tc.data, &aliyunsms)
			require.NoError(t, err)

			require.EqualExportedValues(t, tc.want, aliyunsms)

			data, err := json.Marshal(aliyunsms)
			require.NoError(t, err)

			require.JSONEq(t, tc.wantJSON, string(data))
		})
	}
}
