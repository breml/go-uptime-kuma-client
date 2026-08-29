package notification_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTypeMatchesUpstreamProviderName ensures that every Type() implementation
// of this package, which returns a constant string, reports a notification type
// known to Uptime Kuma.
//
// Uptime Kuma builds its provider registry from the `name` class field of the
// notification providers (server/notification-providers/*.js) and dispatches on
// the stored notification type. A type not present in this registry results in
// a notification, which can not be sent.
func TestTypeMatchesUpstreamProviderName(t *testing.T) {
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

			types[receiverTypeName(funcDecl)] = value
		}
	}

	require.NotEmpty(t, types)

	for receiver, notificationType := range types {
		t.Run(receiver, func(t *testing.T) {
			require.True(t, slices.Contains(upstreamProviderNames, notificationType),
				"Type() returns %q, which is not a provider name known to Uptime Kuma", notificationType)
		})
	}
}

// constantReturnValue returns the constant string returned by the given
// function, if its body consists of a single return statement returning a
// string literal.
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

// receiverTypeName returns the name of the receiver type of the given method.
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
