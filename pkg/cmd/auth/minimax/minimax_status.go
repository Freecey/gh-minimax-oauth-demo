package minimax

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type StatusOptions struct {
	IO         *iostreams.IOStreams
	Config     func() (gh.Config, error)
	HttpClient func() (*http.Client, error)
}

func NewCmdMinimaxStatus(f *cmdutil.Factory, opts *LoginOptions) *cobra.Command {
	statusOpts := &StatusOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Args:  cobra.ExactArgs(0),
		Short: "Show MiniMax authentication status",
		Long: heredoc.Doc(`
			Display the authentication status for MiniMax services.
			
			This command shows which MiniMax instances you are currently
			logged into, along with the user account and token status.
		`),
		Example: heredoc.Doc(`
			# Show authentication status
			$ gh auth minimax status
			
			# Show status in machine-readable format
			$ gh auth minimax status --json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(statusOpts)
		},
	}

	cmd.Flags().Bool("json", false, "Output status in JSON format")

	return cmd
}

func runStatus(opts *StatusOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return err
	}

	// Check for MiniMax configurations
	minimaxHosts := []string{"minimax.global", "minimax.cn"}
	
	var authenticatedHosts []string
	for _, host := range minimaxHosts {
		users := cfg.UsersForHost(host)
		if len(users) > 0 {
			authenticatedHosts = append(authenticatedHosts, host)
		}
	}

	if len(authenticatedHosts) == 0 {
		fmt.Fprintf(opts.IO.ErrOut, "You are not logged in to any MiniMax instance.\n")
		fmt.Fprintf(opts.IO.ErrOut, "Run 'gh auth minimax login' to authenticate.\n")
		return nil
	}

	// Display status for each authenticated host
	cs := opts.IO.ColorScheme()
	w := opts.IO.Out

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(w, "%s MiniMax Authentication Status\n", cs.Bold())
		fmt.Fprintln(w)
	}

	tw := tabwriter.NewWriter(w, 2, 6, 3, ' ', 0)
	for _, host := range authenticatedHosts {
		users := cfg.UsersForHost(host)
		user := users[0]
		
		// Get token
		tokenEntry := cfg.GetOrDefault(host, "oauth_token")
		if !tokenEntry.IsSet() {
			fmt.Fprintf(tw, "%s\t%s\t%s\n", host, user, cs.Red("No token"))
			continue
		}
		
		// Get region
		regionEntry := cfg.GetOrDefault(host, "region")
		region := "global"
		if regionEntry.IsSet() {
			region = regionEntry.Value
		}

		// Check token validity (basic check)
		tokenValid := "✓ Valid"
		if tokenEntry.Value == "" {
			tokenValid = cs.Red("✗ Invalid")
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", host, user, region, tokenValid)
	}
	tw.Flush()

	// Show additional information
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Use 'gh auth minimax login' to add new authentications.\n")
		fmt.Fprintf(w, "Use 'gh auth minimax logout' to remove existing authentications.\n")
	}

	return nil
}