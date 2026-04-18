package authflow

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/atotto/clipboard"
	"github.com/Freecey/gh-minimax-oauth-demo/api"
	"github.com/Freecey/gh-minimax-oauth-demo/internal/browser"
	"github.com/Freecey/gh-minimax-oauth-demo/internal/ghinstance"
	"github.com/Freecey/gh-minimax-oauth-demo/pkg/iostreams"
	"github.com/cli/oauth"

	ghauth "github.com/cli/go-gh/v2/pkg/auth"
)

var (
	// The "GitHub CLI" OAuth app
	oauthClientID     = "178c6fc778ccc68e1d6a"
	oauthClientSecret = "34ddeff2b558a23d38fba8a6de74f086ede1cc0b"

	// MiniMax OAuth app configuration
	minimaxOAuthClientID     = "78257093-7e40-4613-99e0-527b14b39113"
	minimaxOAuthClientSecret = "" // Public client, no secret needed
	minimaxOAuthBaseURL      = "https://api.minimax.io/oauth"
	minimaxOAuthCNBaseURL    = "https://api.minimaxi.com/oauth"
	minimaxOAuthScopes       = []string{"openid", "profile", "model.completion"}
)

// OAuthProvider represents an OAuth provider configuration
type OAuthProvider struct {
	Name         string
	ClientID     string
	ClientSecret string
	BaseURL      string
	Scopes       []string
}

// ProviderOAuthFlow handles OAuth flow for different providers
type ProviderOAuthFlow struct {
	Provider     OAuthProvider
	HTTPClient   *http.Client
	IO           *iostreams.IOStreams
	Host         string
	Interactive  bool
	Browser      browser.Browser
	CopyToClipboard bool
}

// NewProviderOAuthFlow creates a new OAuth flow for the specified provider
func NewProviderOAuthFlow(provider OAuthProvider, httpClient *http.Client, io *iostreams.IOStreams, host string, interactive bool, browser browser.Browser, copyToClipboard bool) *ProviderOAuthFlow {
	return &ProviderOAuthFlow{
		Provider:        provider,
		HTTPClient:      httpClient,
		IO:             io,
		Host:           host,
		Interactive:    interactive,
		Browser:        browser,
		CopyToClipboard: copyToClipboard,
	}
}

// ProviderAuthFlow initiates an OAuth device flow for the specified provider
func ProviderAuthFlow(httpClient *http.Client, provider OAuthProvider, host string, IO *iostreams.IOStreams, notice string, additionalScopes []string, isInteractive bool, b browser.Browser, isCopyToClipboard bool) (string, string, error) {
	w := IO.ErrOut
	cs := IO.ColorScheme()

	minimumScopes := provider.Scopes
	scopes := append(minimumScopes, additionalScopes...)

	// Determine the OAuth host based on provider
	var oauthHost string
	switch provider.Name {
	case "minimax":
		oauthHost = provider.BaseURL
	case "minimax-cn":
		oauthHost = provider.BaseURL
	default:
		return "", "", fmt.Errorf("unsupported OAuth provider: %s", provider.Name)
	}

	hostInstance, err := oauth.NewGitHubHost(ghinstance.HostPrefix(oauthHost))
	if err != nil {
		return "", "", err
	}

	flow := &oauth.Flow{
		Host:         hostInstance,
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		CallbackURI:  getCallbackURI(oauthHost),
		Scopes:       scopes,
		DisplayCode: func(code, verificationURL string) error {
			if isCopyToClipboard {
				err := clipboard.WriteAll(code)
				if err == nil {
					fmt.Fprintf(w, "%s One-time code (%s) copied to clipboard\n", cs.Yellow("!"), cs.Bold(code))
					return nil
				}
				fmt.Fprintf(w, "%s Failed to copy one-time code to clipboard\n", cs.Red("!"))
				fmt.Fprintf(w, "  %s\n", err)
			}
			fmt.Fprintf(w, "%s First copy your one-time code: %s\n", cs.Yellow("!"), cs.Bold(code))
			return nil
		},
		BrowseURL: func(authURL string) error {
			if u, err := url.Parse(authURL); err == nil {
				if u.Scheme != "http" && u.Scheme != "https" {
					return fmt.Errorf("invalid URL: %s", authURL)
				}
			} else {
				return err
			}

			if !isInteractive {
				fmt.Fprintf(w, "%s to continue in your web browser: %s\n", cs.Bold("Open this URL"), authURL)
				return nil
			}

			fmt.Fprintf(w, "%s to open %s in your browser... ", cs.Bold("Press Enter"), authURL)
			_ = waitForEnter(IO.In)

			if err := b.Browse(authURL); err != nil {
				fmt.Fprintf(w, "%s Failed opening a web browser at %s\n", cs.Red("!"), authURL)
				fmt.Fprintf(w, "  %s\n", err)
				fmt.Fprint(w, "  Please try entering the URL in your browser manually\n")
			}
			return nil
		},
		WriteSuccessHTML: func(w io.Writer) {
			fmt.Fprint(w, oauthSuccessPage)
		},
		HTTPClient: httpClient,
		Stdin:      IO.In,
		Stdout:     w,
	}

	fmt.Fprintln(w, notice)

	token, err := flow.DetectFlow()
	if err != nil {
		return "", "", err
	}

	userLogin, err := getViewer(httpClient, oauthHost, token.Token)
	if err != nil {
		return "", "", err
	}

	return token.Token, userLogin, nil
}

// GetMiniMaxProvider returns the MiniMax OAuth provider configuration
func GetMiniMaxProvider(region string) OAuthProvider {
	if region == "cn" {
		return OAuthProvider{
			Name:         "minimax-cn",
			ClientID:     minimaxOAuthClientID,
			ClientSecret: minimaxOAuthClientSecret,
			BaseURL:      minimaxOAuthCNBaseURL,
			Scopes:       minimaxOAuthScopes,
		}
	}
	
	return OAuthProvider{
		Name:         "minimax",
		ClientID:     minimaxOAuthClientID,
		ClientSecret: minimaxOAuthClientSecret,
		BaseURL:      minimaxOAuthBaseURL,
		Scopes:       minimaxOAuthScopes,
	}
}

// waitForEnter waits for the user to press Enter
func waitForEnter(in io.Reader) error {
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	return scanner.Err()
}