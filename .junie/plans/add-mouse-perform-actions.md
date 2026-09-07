---
sessionId: session-260907-124925-msvq
---

# Requirements

### Overview & Goals
Compare the library’s command surface with the current Firefox Marionette command registry in `remote/marionette/driver.sys.mjs`, record the gaps, and implement the first missing feature: W3C `WebDriver:PerformActions` for mouse pointer input.

### Scope
#### In Scope
- Create a command-support inventory matching every current Firefox driver command to its existing Go API, missing status, or intentionally unsupported/internal status.
- Add a public, typed API for mouse action sources and the W3C actions `pointerMove`, `pointerDown`, `pointerUp`, and `pause`.
- Support pointer move origins of `viewport`, current `pointer`, and a `WebElement` reference.
- Send multiple synchronized action sources/ticks through `WebDriver:PerformActions`.
- Verify JSON payloads without Firefox and mouse behavior in the existing ordered Firefox integration suite.
- Document a minimal usage example.

#### Out of Scope
- Keyboard, touch/pen, and wheel action builders.
- Implementing all other commands discovered by the audit.
- Refactoring the Marionette transport or protocol framing.

### Acceptance Criteria
- A caller can compose a mouse sequence, including element-relative movement and button press/release, and submit it through `Client`.
- The wire payload is `{ "actions": [...] }` and follows the current W3C WebDriver action-source schema accepted by Firefox.
- Invalid client-side models that cannot form a valid mouse sequence are rejected before transport where practical.
- Existing APIs and tests remain compatible.
- The repository contains a dated/source-linked support matrix that makes future command work independently actionable.

# Technical Design

### Current Implementation
- `client.go` is a flat public façade whose methods build command payloads and call `Transporter.Send`; it currently contains 54 direct Marionette command calls.
- `transport.go:13–19` exposes a generic `Send(command string, values any)` boundary, and `proto.go:36–52` already serializes nested values as protocol-v3 JSON, so no wire-layer change is needed.
- `WebElement` stores its W3C ID privately in `webelement.go:30–37`; action origins therefore need an explicit serializer that emits `{ "element-6066-11e4-a52e-4f735466cecf": id }`.
- `client_test.go` runs one ordered live session and currently covers element clicks/keys but no actions command. There is no mock transport test pattern today.

### Key Decisions
- Add action models in a focused `actions.go` file while preserving the project’s single-package organization.
- Model the W3C wire contract directly: an input source has `type`, `id`, optional pointer `parameters`, and an ordered `actions` list. Mouse sources set `parameters.pointerType` to `mouse`.
- Represent pointer actions with typed constructors/models rather than exposing unstructured maps. Keep pointer-origin construction controlled so string origins and element origins serialize correctly.
- Keep `PerformActions` on `Client`; it validates the source-level invariants and calls `c.transport.Send("WebDriver:PerformActions", map[string]any{"actions": sources})`.
- Add a test-only recording `Transporter` in `actions_test.go` to assert command and payload independently of Firefox.

### Proposed API
```go
type ActionSequence struct { /* W3C source id/type/parameters/actions */ }
type PointerAction struct { /* pointerMove/down/up or pause fields */ }

func MouseActions(id string, actions ...PointerAction) ActionSequence
func PointerMove(x, y int, duration time.Duration, origin PointerOrigin) PointerAction
func PointerDown(button int) PointerAction
func PointerUp(button int) PointerAction
func Pause(duration time.Duration) PointerAction
func ViewportOrigin() PointerOrigin
func PointerOriginCurrent() PointerOrigin
func ElementOrigin(element *WebElement) PointerOrigin
func (c *Client) PerformActions(actions ...ActionSequence) (*Response, error)
```
Exact exported names may be adjusted to avoid collisions with existing package symbols, but the typed behavior and W3C JSON contract remain fixed.

### Command Gap Audit
Add `PROTOCOL_SUPPORT.md` using the current official `driver.sys.mjs` command table as source of truth. Group commands by session, browsing context/navigation, script, element, actions, cookies, screenshots/print, prompts, Marionette extensions, add-ons, WebAuthn, and platform-specific commands; for each entry record the Firefox command name, current Go method if present, status, and recommended implementation priority. This distinguishes genuine client gaps such as actions from commands that are internal, platform-specific, or outside the library’s intended public surface.

### Files
- Add `actions.go`: public action models, constructors, validation, and `Client.PerformActions`.
- Add `actions_test.go`: recording-transport payload and validation tests.
- Update `webelement.go` only if element-reference JSON marshaling is best shared with future commands; otherwise keep origin serialization private to `actions.go`.
- Update `client_test.go`: register a mouse integration subtest before session teardown.
- Update `testdata/html/form.html`: add deterministic pointer/click state observable by the test.
- Update `README.md`: mouse-actions example and support-matrix link.
- Add `PROTOCOL_SUPPORT.md`: exhaustive current command comparison.

### Risks
- W3C `origin` is a union of strings and element references; custom marshaling and exact payload tests prevent malformed element origins.
- Durations are milliseconds on the wire; reject negative values and convert `time.Duration` explicitly.
- Button-down state can leak after a failed sequence. Integration coverage will use balanced down/up actions; broader release semantics can be prioritized in the support matrix.

# Testing

### Validation Approach
Use fast unit tests for command/payload correctness and one deterministic live Firefox test for end-to-end behavior. Run `go test ./...` and the project’s existing formatting/static checks from the `Makefile` where available.

### Key Scenarios
- Serialize mouse source metadata and ordered move/down/pause/up actions with millisecond durations.
- Serialize viewport, pointer, and W3C element-reference origins exactly.
- Submit multiple sources without changing their order or tick alignment.
- Navigate to local `form.html`, move to a target element, press/release the primary button, and assert fixture state changed.

### Edge Cases
- Empty source ID or action list.
- Nil element origin.
- Negative duration and invalid mouse button values.
- Transport errors propagate unchanged.
- Existing element, navigation, and protocol tests continue to pass.

# Delivery Steps

### ✓ Step 1: Publish the current Marionette command support inventory
The repository contains an exhaustive, source-linked matrix of current Firefox Marionette commands and library gaps.

- Compare every command registered at the bottom of current `remote/marionette/driver.sys.mjs` with `Client` and `WebElement` methods in `client.go` and `webelement.go`.
- Add `PROTOCOL_SUPPORT.md` with supported, missing, extension/platform-specific, and intentionally excluded statuses.
- Group and prioritize missing commands, identifying `WebDriver:PerformActions` as the first implementation slice.
- Link the matrix from `README.md` so later protocol updates have a maintained baseline.

### ✓ Step 2: Implement typed mouse action modeling and command serialization
Callers can construct valid W3C mouse sequences and send them through `WebDriver:PerformActions`.

- Add `actions.go` with typed mouse source, pointer action, pause, and origin models.
- Encode viewport/current-pointer origins as strings and element origins with `WebdriverElementKey` from `client.go`.
- Add `Client.PerformActions` using the existing `Transporter.Send` boundary; leave `transport.go` and protocol-v3 framing unchanged.
- Validate malformed IDs, origins, durations, buttons, and empty sequences before sending.
- Add recording-transport tests in `actions_test.go` for exact commands, nested payloads, multiple sources, validation failures, and transport-error propagation.

### ✓ Step 3: Verify mouse actions against a live Firefox session
The ordered integration suite proves a composed mouse move/down/up sequence affects a deterministic local page.

- Extend `testdata/html/form.html` with a stable mouse target and observable event/click state.
- Add and register a `PerformActions` subtest in `client_test.go` before session teardown.
- Exercise an element-origin move followed by balanced primary-button down/up actions and assert the resulting DOM state through existing client APIs.
- Run formatting and the complete Go test suite, correcting regressions while preserving the suite’s required sequential execution.
- Add a concise `README.md` usage example matching the tested public API.

### ✓ Step 4: Stabilize all tests on Firefox 141.0.3
The complete unit and ordered integration suite passes against Firefox 141.0.3 without relying on external websites.

- Audit existing unit and integration failures and replace network-dependent fixtures with deterministic local pages where practical.
- Update tests or library behavior for Firefox 141.0.3 protocol and prompt semantics without weakening assertions.
- Add a test helper that launches Firefox with a fresh profile, required Marionette flags, readiness checking, and cleanup for every complete test run.
- Run transport-independent tests first, then perform one complete self-contained Firefox-backed run.

### ✓ Step 5: Add Firefox 141.0.3 to GitHub Actions
GitHub Actions installs and launches Firefox 141.0.3 with the flags required by the ordered Marionette integration suite.

- Inspect the current workflow and preserve its supported Go matrix.
- Add Firefox 141.0.3 installation and configure the test helper to use the installed binary.
- Validate workflow syntax and rerun local transport-independent checks after the CI change.

### ✓ Step 6: Refresh usage and test documentation
The README documents mouse actions and makes the self-managed Firefox test lifecycle explicit.

- Add or revise the `PerformActions` example to match the public typed API.
- Review existing documentation for stale setup and test instructions.
- Explain that integration tests launch Firefox with a fresh temporary profile and the required command-line flags.