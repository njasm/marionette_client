# Marionette protocol support

Inventory date: **2026-09-07**. Source of truth: Firefox `main`,
[`remote/marionette/driver.sys.mjs`](https://searchfox.org/firefox-main/source/remote/marionette/driver.sys.mjs),
`GeckoDriver.#commandHandlers`.

Statuses: **supported** has a public Go API; **missing** is a public protocol gap;
**extension** is Firefox-specific and not currently exposed; **internal** is test or
browser-internal infrastructure; **legacy** is sent by this client but is absent from
the current registry. Priorities apply only to missing public commands.

## Session, navigation, and browsing contexts

| Firefox command | Go API | Status | Priority |
|---|---|---|---|
| `WebDriver:NewSession` | `Client.NewSession` | supported | — |
| `WebDriver:DeleteSession` | `Client.DeleteSession` | supported | — |
| `WebDriver:SetTimeouts` | `Client.SetTimeouts` and timeout helpers | supported | — |
| `WebDriver:GetTimeouts` | `Client.GetTimeouts` | supported | — |
| `WebDriver:Navigate` | `Client.Navigate` | supported | — |
| `WebDriver:GetCurrentURL` | `Client.Url` | supported | — |
| `WebDriver:GetTitle` | `Client.Title` | supported | — |
| `WebDriver:Back` | `Client.Back` | supported | — |
| `WebDriver:Forward` | `Client.Forward` | supported | — |
| `WebDriver:Refresh` | `Client.Refresh` | supported | — |
| `WebDriver:GetPageSource` | `Client.PageSource` | supported | — |
| `WebDriver:GetWindowHandle` | `Client.GetWindowHandle` | supported | — |
| `WebDriver:GetWindowHandles` | `Client.GetWindowHandles` | supported | — |
| `WebDriver:SwitchToWindow` | `Client.SwitchToWindow` | supported | — |
| `WebDriver:NewWindow` | `Client.NewWindow` | supported | — |
| `WebDriver:CloseWindow` | `Client.CloseWindow` | supported | — |
| `WebDriver:CloseChromeWindow` | `Client.CloseChromeWindow` | supported | — |
| `WebDriver:GetWindowRect` | `Client.GetWindowRect` | supported | — |
| `WebDriver:SetWindowRect` | `Client.SetWindowRect` | supported | — |
| `WebDriver:MaximizeWindow` | `Client.MaximizeWindow` | supported | — |
| `WebDriver:MinimizeWindow` | `Client.MinimizeWindow` | supported | — |
| `WebDriver:FullscreenWindow` | `Client.FullscreenWindow` | supported | — |
| `WebDriver:SwitchToFrame` | `Client.SwitchToFrame` | supported | — |
| `WebDriver:SwitchToParentFrame` | `Client.SwitchToParentFrame` | supported | — |

## Scripts and elements

| Firefox command | Go API | Status | Priority |
|---|---|---|---|
| `WebDriver:ExecuteScript` | `Client.ExecuteScript` | supported | — |
| `WebDriver:ExecuteAsyncScript` | `Client.ExecuteAsyncScript` | supported | — |
| `WebDriver:FindElement` | `Client.FindElement`, `WebElement.FindElement` | supported | — |
| `WebDriver:FindElements` | `Client.FindElements`, `WebElement.FindElements` | supported | — |
| `WebDriver:GetActiveElement` | `Client.GetActiveElement`, `WebElement.GetActiveElement` | supported | — |
| `WebDriver:FindElementFromShadowRoot` | — | missing | P2 |
| `WebDriver:FindElementsFromShadowRoot` | — | missing | P2 |
| `WebDriver:GetShadowRoot` | — | missing | P2 |
| `WebDriver:ElementClick` | `WebElement.Click` | supported | — |
| `WebDriver:ElementClear` | `WebElement.Clear` | supported | — |
| `WebDriver:ElementSendKeys` | `WebElement.SendKeys` (text and complete special-key set) | supported | — |
| `WebDriver:GetElementAttribute` | `WebElement.Attribute` | supported | — |
| `WebDriver:GetElementCSSValue` | `WebElement.CssValue` | supported | — |
| `WebDriver:GetElementProperty` | `WebElement.Property` | supported | — |
| `WebDriver:GetElementRect` | `WebElement.Rect`, `Location`, `Size` | supported | — |
| `WebDriver:GetElementTagName` | `WebElement.TagName` | supported | — |
| `WebDriver:GetElementText` | `WebElement.Text` | supported | — |
| `WebDriver:IsElementDisplayed` | `WebElement.Displayed` | supported | — |
| `WebDriver:IsElementEnabled` | `WebElement.Enabled` | supported | — |
| `WebDriver:IsElementSelected` | `WebElement.Selected` | supported | — |
| `WebDriver:GetComputedLabel` | — | missing | P2 |
| `WebDriver:GetComputedRole` | — | missing | P2 |

## Actions, cookies, capture, and prompts

| Firefox command | Go API | Status | Priority |
|---|---|---|---|
| `WebDriver:PerformActions` | `Client.PerformActions` (mouse and keyboard) | supported | — |
| `WebDriver:ReleaseActions` | `Client.ReleaseActions` | supported | — |
| `WebDriver:AddCookie` | `Client.AddCookie` | supported | — |
| `WebDriver:GetCookies` | `Client.GetCookies` | supported | — |
| `WebDriver:DeleteCookie` | `Client.DeleteCookie` | supported | — |
| `WebDriver:DeleteAllCookies` | `Client.DeleteAllCookies` | supported | — |
| `WebDriver:TakeScreenshot` | `Client.Screenshot`, `WebElement.Screenshot` | supported | — |
| `WebDriver:Print` | — | missing | P2 |
| `WebDriver:DismissAlert` | `Client.DismissAlert` | supported | — |
| `WebDriver:AcceptAlert` | `Client.AcceptAlert` | supported | — |
| `WebDriver:GetAlertText` | `Client.TextFromAlert` | supported | — |
| `WebDriver:SendAlertText` | `Client.SendAlertText` | supported | — |

## Firefox extensions and platform facilities

| Firefox command | Go API | Status | Priority |
|---|---|---|---|
| `Addon:Install` | — | extension | — |
| `Addon:Uninstall` | — | extension | — |
| `L10n:LocalizeProperty` | — | extension | — |
| `Marionette:AcceptConnections` | — | internal | — |
| `Marionette:GetAccessibilityPropertiesForAccessibilityNode` | — | extension | — |
| `Marionette:GetAccessibilityPropertiesForElement` | — | extension | — |
| `Marionette:GetContext` | `Client.Context` | supported | — |
| `Marionette:GetScreenOrientation` | — | extension | — |
| `Marionette:GetWindowType` | — | extension | — |
| `Marionette:Quit` | `Client.Quit` | supported | — |
| `Marionette:RegisterChromeHandler` | — | internal | — |
| `Marionette:SetContext` | `Client.SetContext` | supported | — |
| `Marionette:SetScreenOrientation` | — | extension | — |
| `Marionette:UnregisterChromeHandler` | — | internal | — |
| `GPC:GetGlobalPrivacyControl` | — | extension | — |
| `GPC:SetGlobalPrivacyControl` | — | extension | — |
| `WebDriver:SetPermission` | — | extension (compatibility alias) | — |
| `Permissions:SetPermission` | — | extension | — |
| `Reporting:GenerateTestReport` | — | internal | — |

## Reftest and WebAuthn

| Firefox command | Go API | Status | Priority |
|---|---|---|---|
| `reftest:run` | — | internal | — |
| `reftest:setup` | — | internal | — |
| `reftest:teardown` | — | internal | — |
| `WebAuthn:AddCredential` | — | missing | P3 |
| `WebAuthn:AddVirtualAuthenticator` | — | missing | P3 |
| `WebAuthn:GetCredentials` | — | missing | P3 |
| `WebAuthn:RemoveCredential` | — | missing | P3 |
| `WebAuthn:RemoveAllCredentials` | — | missing | P3 |
| `WebAuthn:RemoveVirtualAuthenticator` | — | missing | P3 |
| `WebAuthn:SetUserVerified` | — | missing | P3 |

## Legacy client command

`Client.GetCapabilities` sends `WebDriver:GetCapabilities`, which is not present in
the current Firefox registry. It is retained for compatibility but should not be used
as the basis for new APIs; capabilities are returned by `WebDriver:NewSession`.

## Recommended implementation order

Completed baseline: `WebDriver:PerformActions` with typed mouse and keyboard input, the complete WebDriver special-key
set for element entry, and `WebDriver:ReleaseActions`.

1. P2: shadow-root access, computed accessibility properties, and print.
2. P3: WebAuthn virtual-authenticator support.
