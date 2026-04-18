package minimax

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/Freecey/gh-minimax-oauth-demo/internal/authflow"
	"github.com/Freecey/gh-minimax-oauth-demo/internal/browser"
	"github.com/Freecey/gh-minimax-oauth-demo/internal/gh"
	"github.com/Freecey/gh-minimax-oauth-demo/pkg/cmd/auth/shared"
	"github.com/Freecey/gh-minimax-oauth-demo/pkg/cmdutil"
	"github.com/Freecey/gh-minimax-oauth-demo/pkg/iostreams"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	IO              *iostreams.IOStreams
	Config          func() (gh.Config, error)
	HttpClient      func() (*http.Client, error)
	PlainHttpClient func() (*http.Client, error)
	Prompter        shared.Prompt
	Browser         browser.Browser

	Region          string
	Interactive     bool
	Web             bool
	Clipboard       bool
}

func NewCmdMinimax(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{
		IO:              f.IOStreams,
		Config:          f.Config,
		HttpClient:      f.HttpClient,
		PlainHttpClient: f.PlainHttpClient,
		Prompter:        f.Prompter,
		Browser:         f.Browser,
	}

	cmd := &cobra.Command{
		Use:   "minimax <command>",
		Short: "Authenticate with MiniMax AI services",
		Long: heredoc.Doc(`
			Authenticate with MiniMax AI services using OAuth.
			
			This command allows you to securely authenticate with MiniMax's AI services
			using OAuth 2.0 device code flow. Your credentials are stored securely
			and can be used across different MiniMax tools and services.
		`),
		Example: heredoc.Doc(`
			# Start interactive authentication with global MiniMax
			$ gh auth minimax login
			
			# Authenticate with China MiniMax region
			$ gh auth minimax login --region cn
			
			# Check authentication status
			$ gh auth minimax status
			
			# Logout from MiniMax
			$ gh auth minimax logout
		`),
	}

	cmd.AddCommand(NewCmdMinimaxLogin(f, opts))
	cmd.AddCommand(NewCmdMinimaxStatus(f, opts))
	cmd.AddCommand(NewCmdMinimaxLogout(f, opts))
	cmd.AddCommand(NewCmdMinimaxRefresh(f, opts))

	return cmd
}

func NewCmdMinimaxLogin(f *cmdutil.Factory, opts *LoginOptions) *cobra.Command {
	var region string

	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Log in to MiniMax",
		Long: heredoc.Doc(`
			Log in to MiniMax using OAuth 2.0 device code flow.
			
			The default region is global (api.minimax.io). You can specify
			--region cn to use the China endpoint (api.minimaxi.com).
		`),
		Example: heredoc.Doc(`
			# Login to global MiniMax (default)
			$ gh auth minimax login
			
			# Login to China MiniMax
			$ gh auth minimax login --region cn
			
			# Login with web browser flow
			$ gh auth minimax login --web
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Region = region
			opts.Interactive = true
			opts.Web = f.IOStreams.IsStdoutTTY()
			opts.Clipboard = f.IOStreams.IsStdoutTTY()
			
			if !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf("authentication required when running in non-interactive mode")
			}

			return runLogin(opts)
		},
	}

	cmd.Flags().StringVarP(&region, "region", "r", "global", "MiniMax region (global or cn)")
	cmd.Flags().BoolVar(&opts.Web, "web", false, "Use web browser for authentication")

	return cmd
}

func runLogin(opts *LoginOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	httpClient, err := opts.PlainHttpClient()
	if err != nil {
		return err
	}

	// Normalize region
	region := strings.ToLower(opts.Region)
	if region != "global" && region != "cn" {
		return fmt.Errorf("invalid region: %s (must be 'global' or 'cn')", opts.Region)
	}

	// Get MiniMax provider configuration
	provider := authflow.GetMiniMaxProvider(region)

	// Determine hostname for storage
	hostname := fmt.Sprintf("minimax.%s", region)

	// Check if already logged in
	usersForHost := cfg.UsersForHost(hostname)
	if len(usersForHost) > 0 {
		confirm, err := opts.Prompter.Confirm(
			fmt.Sprintf("You are already logged in to %s as %s. Log in again?", hostname, usersForHost[0]),
			false,
		)
		if err != nil {
			return err
		}
		if !confirm {
			return fmt.Errorf("authentication cancelled")
		}
	}

	// Perform OAuth flow
	notice := fmt.Sprintf("Logging in to MiniMax (%s region)", strings.ToUpper(region))
	
	authToken, username, err := authflow.ProviderAuthFlow(
		httpClient,
		provider,
		hostname,
		opts.IO,
		notice,
		nil, // additional scopes
		opts.Interactive,
		opts.Browser,
		opts.Clipboard,
	)
	if err != nil {
		return fmt.Errorf("failed to authenticate with MiniMax: %w", err)
	}

	// Store the authentication
	insecureStorageUsed, err := cfg.Login(hostname, username, authToken, "https", true)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	// Set region in config
	cfg.Set(hostname, "region", region)

	// Show success message
	cs := opts.IO.ColorScheme()
	fmt.Fprintf(opts.IO.ErrOut, "%s Authentication complete.\n", cs.SuccessIcon())

	if insecureStorageUsed {
		fmt.Fprintf(opts.IO.ErrOut, "%s Your credentials have been saved to plain text file.\n", cs.Yellow("!"))
		fmt.Fprintf(opts.IO.ErrOut, "  Note: GitHub CLI will prefer a system credential helper when available.\n")
	}

	fmt.Fprintf(opts.IO.ErrOut, "%s Ready to use MiniMax models from %s region.\n", cs.SuccessIcon(), strings.ToUpper(region))

	return nil
}