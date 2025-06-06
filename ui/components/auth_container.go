// Package components defines various UI components for the LazyPost application.
package components

import (
	"encoding/base64"
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// authTypeOptions lists the available authentication types for the AuthSelector.
var authTypeOptions = []string{"None", "Basic", "Bearer", "JWT", "OAuth2", "API Key"}

// AuthSelectorKeyMap defines keybindings for the AuthSelector component.
// These bindings are used when the AuthSelector is active and its dropdown is open or closed.
// (Placeholder for future more complex interactivity, currently uses simple string matching).
type AuthSelectorKeyMap struct {
	Open   key.Binding // Key to open the dropdown.
	Close  key.Binding // Key to close the dropdown.
	Next   key.Binding // Key to navigate to the next option in the dropdown.
	Prev   key.Binding // Key to navigate to the previous option in the dropdown.
	Select key.Binding // Key to select the highlighted option in the dropdown.
}

// DefaultAuthSelectorKeyMap provides default keybindings for the AuthSelector.
// These are standard keys like Enter, Space, Escape, and arrow keys.
var DefaultAuthSelectorKeyMap = AuthSelectorKeyMap{
	Open:   key.NewBinding(key.WithKeys("enter", " ")),
	Close:  key.NewBinding(key.WithKeys("esc")),
	Next:   key.NewBinding(key.WithKeys("down", "j")),
	Prev:   key.NewBinding(key.WithKeys("up", "k")),
	Select: key.NewBinding(key.WithKeys("enter")),
}

// AuthSelector manages the dropdown UI for selecting an authentication type.
// It handles opening/closing the dropdown, navigating options, and displaying the current selection.
type AuthSelector struct {
	options            []string         // options are the available authentication type strings.
	selectedIndex      int              // selectedIndex is the index of the currently chosen option.
	highlightedIndex   int              // highlightedIndex is the index of the option highlighted when the dropdown is open.
	isOpen             bool             // isOpen indicates whether the dropdown list is visible.
	active             bool             // active indicates whether the component is currently focused and interactive.
	width              int              // width is the rendering width of the component.
	activeStyle        lipgloss.Style   // activeStyle is the style applied when the component is active.
	inactiveStyle      lipgloss.Style   // inactiveStyle is the style applied when the component is inactive.
	dropdownTextStyle  lipgloss.Style   // dropdownTextStyle is the style for text within the dropdown.
	dropdownArrowStyle lipgloss.Style   // dropdownArrowStyle is the style for the dropdown arrow indicator.
	dropdownItemStyle  lipgloss.Style   // dropdownItemStyle is the style for individual items when the dropdown is open.
	keymap             AuthSelectorKeyMap // keymap holds the keybindings for interacting with the selector.
}

// NewAuthSelector creates and initializes a new AuthSelector component.
// It sets default options, styles, and keymap.
func NewAuthSelector() AuthSelector {
	return AuthSelector{
		options:            authTypeOptions,
		selectedIndex:      0,
		highlightedIndex:   0, // Initialize highlightedIndex
		isOpen:             false,
		active:             false,
		activeStyle:        styles.DefaultTheme.ActiveInputStyle.Copy(),
		inactiveStyle:      styles.DefaultTheme.InactiveInputStyle.Copy(),
		dropdownTextStyle:  styles.DefaultTheme.DropdownTextStyle.Copy(),
		dropdownArrowStyle: styles.DefaultTheme.DropdownArrowStyle.Copy(),
		dropdownItemStyle: styles.DefaultTheme.DropdownItemStyle.Copy(),             // Initialize new style
		keymap:             DefaultAuthSelectorKeyMap,
	}
}

// View renders the AuthSelector component.
// It displays either the currently selected option (if closed) or the list of options (if open).
// Styling is applied based on the active state and whether the dropdown is open.
func (as AuthSelector) View() string {
	var currentStyle lipgloss.Style
	if as.active {
		currentStyle = as.activeStyle
	} else {
		currentStyle = as.inactiveStyle
	}

	// Calculate effective width for content inside the style's padding
	effectiveContentWidth := as.width - currentStyle.GetHorizontalPadding()
	if effectiveContentWidth < 0 {
		effectiveContentWidth = 0
	}

	if !as.isOpen {
		selectedOptionText := as.dropdownTextStyle.Render(as.options[as.selectedIndex])
		arrow := as.dropdownArrowStyle.Render(" ▼")
		optionStrPaddedWidth := effectiveContentWidth - lipgloss.Width(arrow) - 2 // -2 for spaces around text and arrow
		if optionStrPaddedWidth < 0 {
			optionStrPaddedWidth = 0
		}
		viewString := fmt.Sprintf(" %-*s%s ", optionStrPaddedWidth, selectedOptionText, arrow)
		return currentStyle.Copy().Width(as.width).Render(viewString)
	}

	// Render open state
	var items []string
	for i, optionText := range as.options {
		var renderedText string
		if i == as.highlightedIndex {
			displayText := "▶ " + optionText
			renderedText = styles.DefaultTheme.SelectedItemStyle.Render(displayText)
		} else {
			displayText := "  " + optionText
			renderedText = as.dropdownTextStyle.Render(displayText)
		}
		// Each item line uses dropdownItemStyle for padding and width, then renders the specific text style inside.
		line := as.dropdownItemStyle.Copy().Width(effectiveContentWidth).Render(renderedText)
		items = append(items, line)
	}

	dropdownContent := lipgloss.JoinVertical(lipgloss.Left, items...)
	return currentStyle.Copy().Width(as.width).Render(dropdownContent)
}

// SetWidth sets the rendering width for the AuthSelector component.
func (as *AuthSelector) SetWidth(width int) {
	as.width = width
}

// SetActive sets the active state of the AuthSelector.
// An active selector can be interacted with via keybindings.
func (as *AuthSelector) SetActive(active bool) {
	as.active = active
}

// Update handles messages for the AuthSelector, primarily key presses.
// It manages opening/closing the dropdown, navigating options, and selecting an item.
// It only processes messages if the selector is active.
func (as *AuthSelector) Update(msg tea.Msg) tea.Cmd {
	if !as.active { // Only process messages if active and the component is supposed to be interactive
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if as.isOpen { // Handle keys when dropdown is open
			switch {
			case key.Matches(msg, as.keymap.Close):
				as.isOpen = false
				as.highlightedIndex = as.selectedIndex // Reset highlight to current selection
				return nil
			case key.Matches(msg, as.keymap.Next):
				as.highlightedIndex = (as.highlightedIndex + 1) % len(as.options)
				return nil
			case key.Matches(msg, as.keymap.Prev):
				as.highlightedIndex = (as.highlightedIndex - 1 + len(as.options)) % len(as.options)
				return nil
			case key.Matches(msg, as.keymap.Select):
				as.selectedIndex = as.highlightedIndex
				as.isOpen = false
				// TODO: Potentially send a message indicating selection changed, if other components need to react immediately.
				// For now, the change will be reflected in the next View() call.
				return nil
			}
		} else { // Handle keys when dropdown is closed
			switch {
			case key.Matches(msg, as.keymap.Open):
				as.isOpen = true
				as.highlightedIndex = as.selectedIndex // Start highlighting from current selection
				return nil
			}
		}
	}
	return nil
}

// AuthContainer encapsulates the AuthSelector and the various authentication detail components.
// It manages which auth detail view is shown based on the AuthSelector's choice
// and delegates updates and focus to the appropriate child component.
type AuthContainer struct {
	Width          int            // Width is the rendering width of the container.
	Height         int            // Height is the rendering height of the container.
	Active         bool           // Active indicates if the container (and potentially its children) is focused.
	authSelector   AuthSelector   // authSelector is the dropdown for choosing auth type.
	titleStyle     lipgloss.Style // titleStyle is used for the container's title (if any, currently unused).

	// Detail components for each authentication type.
	basicAuthDetails   BasicAuthDetailsComponent  // basicAuthDetails handles Basic authentication inputs.
	tokenAuthDetails   TokenAuthDetailsComponent  // tokenAuthDetails handles Bearer token input.
	jwtAuthDetails     JWTAuthDetailsComponent    // jwtAuthDetails handles JWT input.
	apiKeyAuthDetails  APIKeyAuthDetailsComponent // apiKeyAuthDetails handles API Key input.
	oauth2AuthDetails  OAuth2AuthDetailsComponent // oauth2AuthDetails handles OAuth2 details.
}

// NewAuthContainer creates and initializes a new AuthContainer.
// It creates an AuthSelector and instances of all auth detail components.
func NewAuthContainer() AuthContainer {
	selector := NewAuthSelector()
	return AuthContainer{
		Width:          0,
		Height:         0,
		Active:         false,
		authSelector:   selector,
		titleStyle:     styles.DefaultTheme.TitleStyle.Copy(),

		basicAuthDetails:  NewBasicAuthDetailsComponent(),
		tokenAuthDetails:  NewTokenAuthDetailsComponent(),
		jwtAuthDetails:    NewJWTAuthDetailsComponent(), // Initialize new component
		apiKeyAuthDetails: NewAPIKeyAuthDetailsComponent(),
		oauth2AuthDetails: NewOAuth2AuthDetailsComponent(),
	}
}

// SetWidth sets the rendering width for the AuthContainer and its children.
// The width is distributed to the AuthSelector and the active auth detail component.
func (ac *AuthContainer) SetWidth(width int) {
	ac.Width = width
	// Child components' widths will be set during View rendering or specific focus changes.
}

// SetHeight sets the rendering height for the AuthContainer and its children.
// The height is distributed to the AuthSelector and the active auth detail component.
func (ac *AuthContainer) SetHeight(height int) {
	ac.Height = height
	// Child components' heights will be set during View rendering.
}

// SetActive sets the active state of the AuthContainer.
// It also propagates the active state to the AuthSelector and the currently selected auth detail component.
func (ac *AuthContainer) SetActive(active bool) {
	ac.Active = active
	// The authSelector is always potentially interactive if the container is active.
	ac.authSelector.SetActive(active)

	// Deactivate all detail components first
	ac.basicAuthDetails.SetActive(false)
	ac.tokenAuthDetails.SetActive(false)
	ac.jwtAuthDetails.SetActive(false) // Deactivate new component
	ac.apiKeyAuthDetails.SetActive(false)
	ac.oauth2AuthDetails.SetActive(false)

	if active {
		// If the container is active, the selected detail component (if any) should also be marked active.
		// This doesn't mean it has primary focus, just that it's the one to interact with if focus moves there.
		selectedType := ac.authSelector.options[ac.authSelector.selectedIndex]
		switch selectedType {
		case "Basic":
			ac.basicAuthDetails.SetActive(true)
		case "Bearer": // Explicitly Bearer
			ac.tokenAuthDetails.SetActive(true)
		case "JWT": // New case for JWT
			ac.jwtAuthDetails.SetActive(true)
		case "API Key":
			ac.apiKeyAuthDetails.SetActive(true)
		case "OAuth2":
			ac.oauth2AuthDetails.SetActive(true)
		}
	}
}

// Update handles messages for the AuthContainer.
// It delegates messages to the AuthSelector and the currently active auth detail component.
// It also re-evaluates which detail component should be active if the AuthSelector's selection changes.
// It only processes messages if the container itself is active.
func (ac *AuthContainer) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	if !ac.Active {
		return nil
	}

	// Priority 1: AuthSelector, especially if open or specific key presses
	// The selector's Update method already checks its own 'active' state,
	// but we ensure it only gets messages if AuthContainer itself is active.
	selectorCmd := ac.authSelector.Update(msg)
	if selectorCmd != nil {
		cmds = append(cmds, selectorCmd)
	}

	// If selector made a selection, we might need to change active detail component
	// This is implicitly handled by SetActive being called from QueryTab or App,
	// or by re-evaluating in View. For now, direct selection reaction is minimal here.
	// The main thing is that authSelector.selectedIndex has changed.
	// We should ensure the correct detail component is marked active.
	ac.SetActive(ac.Active) // Re-evaluate active detail component

	// Detail component updates: Check which detail component should be active based on selection
	selectedType := ac.authSelector.options[ac.authSelector.selectedIndex]
	var detailCmd tea.Cmd
	switch selectedType {
	case "Basic":
		if ac.basicAuthDetails.active { // Check if it's supposed to be active
			detailCmd = ac.basicAuthDetails.Update(msg)
		}
	case "Bearer": // Explicitly Bearer
		if ac.tokenAuthDetails.active {
			detailCmd = ac.tokenAuthDetails.Update(msg)
		}
	case "JWT": // New case for JWT
		if ac.jwtAuthDetails.active {
			detailCmd = ac.jwtAuthDetails.Update(msg)
		}
	case "API Key":
		if ac.apiKeyAuthDetails.active {
			detailCmd = ac.apiKeyAuthDetails.Update(msg)
		}
	case "OAuth2":
		if ac.oauth2AuthDetails.active {
			detailCmd = ac.oauth2AuthDetails.Update(msg)
		}
	}
	if detailCmd != nil {
		cmds = append(cmds, detailCmd)
	}

	return tea.Batch(cmds...)
}

// View renders the AuthContainer.
// It displays the AuthSelector and the view of the currently selected auth detail component.
// The layout includes spacing between the selector and the detail view.
// The container is enclosed in a border styled according to its active state.
func (ac AuthContainer) View() string {
	if ac.Width == 0 || ac.Height == 0 {
		return ""
	}

	var currentFrameStyle lipgloss.Style
	if ac.Active {
		currentFrameStyle = styles.DefaultTheme.ActiveBorderStyle.Copy()
	} else {
		currentFrameStyle = styles.DefaultTheme.BorderStyle.Copy()
	}

	outerFrame := currentFrameStyle.
		Width(ac.Width).
		Height(ac.Height).
		Padding(0, 1)

	trueInnerWidth := ac.Width - outerFrame.GetHorizontalFrameSize()
	trueInnerHeight := ac.Height - outerFrame.GetVerticalFrameSize()
	if trueInnerWidth < 0 {
		trueInnerWidth = 0
	}
	if trueInnerHeight < 0 {
		trueInnerHeight = 0
	}

	var contentLines []string

	// Part 1: AuthSelector
	// Use a mutable copy of the selector to set width for rendering
	tempSelector := ac.authSelector
	tempSelector.SetWidth(30) // Fixed width for AuthSelector
	// The active state of ac.authSelector is managed by AuthContainer.SetActive
	selectorView := tempSelector.View() // This can be a multi-line block if dropdown is open
	
	// Render the selectorView.
	contentLines = append(contentLines, lipgloss.NewStyle().Width(trueInnerWidth).Render(selectorView))
	
	currentContentHeight := lipgloss.Height(selectorView)
	
	// Part 2: Spacing (3 lines)
	spacingHeight := 3
	if trueInnerHeight > currentContentHeight && spacingHeight > 0 {
		if currentContentHeight+spacingHeight > trueInnerHeight {
			spacingHeight = trueInnerHeight - currentContentHeight
		}
		if spacingHeight > 0 {
			spacingBlock := lipgloss.NewStyle().Width(trueInnerWidth).Height(spacingHeight).Render("")
			contentLines = append(contentLines, spacingBlock)
			currentContentHeight += spacingHeight
		}
	}

	// Part 3: Auth Detail Sub-Container
	detailViewContent := ""
	selectedType := ac.authSelector.options[ac.authSelector.selectedIndex]
	
	detailComponentHeight := trueInnerHeight - currentContentHeight
	if detailComponentHeight < 0 {
		detailComponentHeight = 0
	}

	// Create mutable copies of detail components to set size and get view
	// This is a bit clunky; ideally, SetSize would be called less frequently,
	// or View would take size parameters. For now, this matches the pattern.
	// The active state is already set by ac.SetActive().
	
	// Make a non-pointer copy for view rendering if needed, or ensure methods are value receivers
	// For components like BasicAuthDetailsComponent, since SetSize modifies them,
	// we need to be careful if ac is a value receiver in View.
	// Let's assume these components are simple enough for now.
	// To be safe, we should use pointers or ensure methods handle this.
	// For this iteration, we'll proceed with direct field access/modification on ac's fields.
	// This means AuthContainer methods that modify children (like SetSize on them) should take *AuthContainer.

	// To ensure `SetSize` calls modify the actual components within `ac`,
	// we'll call them on `ac.basicAuthDetails` etc. directly.
	// The `View` methods of these components are value receivers, so they won't modify.

	if detailComponentHeight > 0 {
		switch selectedType {
		case "Basic":
			// ac.basicAuthDetails.SetActive(ac.Active) // Active state set in AuthContainer.SetActive
			ac.basicAuthDetails.SetSize(trueInnerWidth, detailComponentHeight)
			detailViewContent = ac.basicAuthDetails.View()
		case "Bearer": // Explicitly Bearer
			// ac.tokenAuthDetails.SetActive(ac.Active)
			ac.tokenAuthDetails.SetSize(trueInnerWidth, detailComponentHeight)
			detailViewContent = ac.tokenAuthDetails.View()
		case "JWT": // New case for JWT
			// ac.jwtAuthDetails.SetActive(ac.Active)
			ac.jwtAuthDetails.SetSize(trueInnerWidth, detailComponentHeight)
			detailViewContent = ac.jwtAuthDetails.View()
		case "API Key":
			// ac.apiKeyAuthDetails.SetActive(ac.Active)
			ac.apiKeyAuthDetails.SetSize(trueInnerWidth, detailComponentHeight)
			detailViewContent = ac.apiKeyAuthDetails.View()
		case "OAuth2":
			// ac.oauth2AuthDetails.SetActive(ac.Active)
			ac.oauth2AuthDetails.SetSize(trueInnerWidth, detailComponentHeight)
			detailViewContent = ac.oauth2AuthDetails.View()
		case "None":
			// No detail view for "None"
			detailViewContent = ""
		}
		if detailViewContent != "" {
			contentLines = append(contentLines, detailViewContent)
		}
	}
	
	innerContentBlock := lipgloss.JoinVertical(lipgloss.Left, contentLines...)

	// Final padding for the entire container if needed
	paddingHeight := trueInnerHeight - lipgloss.Height(innerContentBlock)
	if paddingHeight < 0 {
		paddingHeight = 0
	}

	var finalInnerContent string
	if paddingHeight > 0 {
		paddingBlock := lipgloss.NewStyle().Width(trueInnerWidth).Height(paddingHeight).Render("")
		finalInnerContent = lipgloss.JoinVertical(lipgloss.Left, innerContentBlock, paddingBlock)
	} else {
		finalInnerContent = innerContentBlock
	}
	
	return outerFrame.Render(finalInnerContent)
}

// GetAuthHeaders constructs and returns a map of HTTP headers based on the selected authentication type
// and the values entered in the corresponding auth detail component.
// For "None", it returns an empty map. For other types, it retrieves credentials/tokens
// and formats them into the appropriate "Authorization" header (or other headers for API Key, if applicable).
// Placeholder comments indicate where logic for JWT, API Key, and OAuth2 needs to be fully implemented.
func (ac AuthContainer) GetAuthHeaders() map[string]string {
	headers := make(map[string]string)
	selectedType := ac.authSelector.options[ac.authSelector.selectedIndex]

	switch selectedType {
	case "Basic":
		username, password := ac.basicAuthDetails.GetValues()
		if username != "" || password != "" { // Only add header if there's a username or password
			auth := username + ":" + password
			headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		}
	case "Bearer":
		// TODO: Implement Bearer token retrieval from tokenAuthDetails
		// token := ac.tokenAuthDetails.GetValue()
		// if token != "" {
		// 	headers["Authorization"] = "Bearer " + token
		// }
	case "JWT":
		// TODO: Implement JWT retrieval from jwtAuthDetails
		// jwt := ac.jwtAuthDetails.GetValue()
		// if jwt != "" {
		// 	headers["Authorization"] = "Bearer " + jwt // Typically Bearer for JWT too
		// }
	case "API Key":
		// TODO: Implement API Key retrieval and header construction from apiKeyAuthDetails
		// e.g., headerName, headerValue, addTo := ac.apiKeyAuthDetails.GetValues()
		// if headerName != "" && headerValue != "" {
		// 	 if addTo == "header" { headers[headerName] = headerValue } ... else if query etc.
		// }
	case "OAuth2":
		// TODO: Implement OAuth2 token retrieval from oauth2AuthDetails
		// This will likely be more complex, involving a token that might be stored
		// accessToken := ac.oauth2AuthDetails.GetAccessToken()
		// if accessToken != "" {
		// 	headers["Authorization"] = "Bearer " + accessToken
		// }
	case "None":
		// No headers to add
	}
	return headers
}

// IsFocused checks if the AuthContainer itself is considered to be in a focused state.
// Currently, this is equivalent to its Active state.
// (Placeholder for potentially more complex focus logic).
func (ac AuthContainer) IsFocused() bool {
	return ac.Active // Simple check for now, might need to check authSelector.active
}
