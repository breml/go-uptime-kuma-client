package notification_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTypeMatchesUpstreamProviderName ensures that the notification types
// reported by this package are exactly the provider names known to Uptime Kuma.
//
// Uptime Kuma builds its provider registry from the `name` class field of the
// notification providers (server/notification-providers/*.js) and dispatches on
// the stored notification type (server/notification.js, Notification.send). A
// type absent from this registry is still accepted on create, so the mismatch
// only surfaces when a notification is sent, which fails with "Notification
// type is not supported".
//
// Only a Type() implemented as a single return of a string literal is
// collected. Wrapper types delegating to their details counterpart are covered
// transitively, Base and Generic report a type read from the notification
// itself and have none to check.
func TestTypeMatchesUpstreamProviderName(t *testing.T) {
	// Provider names of Uptime Kuma 2.5.0, the `name` field of every class in
	// server/notification-providers/*.js. Regenerate after an upstream bump by
	// running the following in a checkout of Uptime Kuma:
	//
	//	grep -hoP '^\s*name = "\K[^"]+' server/notification-providers/*.js | sort -f
	upstreamProviderNames := []string{
		"alerta",
		"AlertNow",
		"AliyunSMS",
		"apprise",
		"bale",
		"Bark",
		"Bitrix24",
		"Brevo",
		"CallMeBot",
		"Cellsynt",
		"clicksendsms",
		"DingDing",
		"discord",
		"egosms",
		"Elks",
		"evolution",
		"Feishu",
		"FlashDuty",
		"Flowtriq",
		"fluxer",
		"FreeMobile",
		"GoAlert",
		"GoogleChat",
		"GoogleSheets",
		"gorush",
		"gotify",
		"GrafanaOncall",
		"gtxmessaging",
		"HaloPSA",
		"HeiiOnCall",
		"HomeAssistant",
		"JiraServiceManagement",
		"Keep",
		"Kook",
		"line",
		"lunasea",
		"matrix",
		"mattermost",
		"max",
		"nextcloudtalk",
		"nostr",
		"notifery",
		"ntfy",
		"octopush",
		"OneBot",
		"OneChat",
		"Onesender",
		"Ooredoo",
		"Opsgenie",
		"PagerDuty",
		"PagerTree",
		"plivo",
		"promosms",
		"pumble",
		"pushbullet",
		"PushByTechulus",
		"PushDeer",
		"pushover",
		"PushPlus",
		"pushy",
		"Resend",
		"rocket.chat",
		"SendGrid",
		"ServerChan",
		"serwersms",
		"SevenIO",
		"signal",
		"SIGNL4",
		"slack",
		"smsc",
		"SMSEagle",
		"smsir",
		"SMSManager",
		"SMSPartner",
		"SMSPlanet",
		"smtp",
		"Splunk",
		"SpugPush",
		"squadcast",
		"stackfield",
		"teams",
		"telegram",
		"telnyx",
		"Teltonika",
		"threema",
		"twilio",
		"VK",
		"VKTeams",
		"waha",
		"webhook",
		"Webpush",
		"WeCom",
		"whapi",
		"Whatsapp360messenger",
		"WPush",
		"WxPusher",
		"YZJ",
		"ZohoCliq",
	}

	files, err := filepath.Glob("*.go")
	require.NoError(t, err)

	fset := token.NewFileSet()
	types := map[string]string{}

	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(fset, name, nil, 0)
		require.NoError(t, err)

		for _, decl := range file.Decls {
			funcDecl, isFunc := decl.(*ast.FuncDecl)
			if !isFunc || funcDecl.Recv == nil || funcDecl.Name.Name != "Type" {
				continue
			}

			value, isConstant := constantReturnValue(funcDecl)
			if !isConstant {
				continue
			}

			receiver := receiverTypeName(funcDecl)
			require.NotEmpty(t, receiver, "Type() in %s has an unsupported receiver", name)

			types[receiver] = value
		}
	}

	for receiver, notificationType := range types {
		t.Run(receiver, func(t *testing.T) {
			require.Contains(t, upstreamProviderNames, notificationType,
				"Type() returns %q, which is not a provider name known to Uptime Kuma", notificationType)
		})
	}

	// Assert the two sets are equal, not just that every collected type is
	// known upstream. This catches a provider whose Type() is no longer a
	// string literal and therefore silently escapes collection, two providers
	// reporting the same type, and a provider added upstream but not
	// implemented here.
	require.ElementsMatch(t, upstreamProviderNames, slices.Collect(maps.Values(types)),
		"the notification types of this package do not match the upstream provider names one to one")
}

// constantReturnValue returns the string returned by the given function, if its
// body consists of a single return statement returning a string literal.
func constantReturnValue(funcDecl *ast.FuncDecl) (string, bool) {
	if funcDecl.Body == nil || len(funcDecl.Body.List) != 1 {
		return "", false
	}

	returnStmt, ok := funcDecl.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returnStmt.Results) != 1 {
		return "", false
	}

	lit, ok := returnStmt.Results[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return value, true
}

// receiverTypeName returns the name of the receiver type of the given method,
// dereferencing a pointer receiver. It returns an empty string, if the receiver
// is not a plain identifier, for example a generic type.
func receiverTypeName(funcDecl *ast.FuncDecl) string {
	expr := funcDecl.Recv.List[0].Type

	starExpr, isStar := expr.(*ast.StarExpr)
	if isStar {
		expr = starExpr.X
	}

	ident, isIdent := expr.(*ast.Ident)
	if !isIdent {
		return ""
	}

	return ident.Name
}
