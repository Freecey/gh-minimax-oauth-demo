package minimax

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/Freecey/gh-minimax-oauth-demo/internal/gh"
	"github.com/Freecey/gh-minimax-oauth-demo/pkg/cmdutil"
	"github.com/Freecey/gh-minimax-oauth-demo/pkg/iostreams"
	"github.com/spf13/cobra"
)

type LogoutOptions struct {
	IO     *iostreams.IOStreams
	Config func() (gh.Config, error)
	Host   string
}

func NewCmdMinimaxLogout(f *cmdutil.Factory, opts *LoginOptions) *cobra.Command {
	logoutOpts := &LogoutOptions{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "logout",
		Args:  cobra.ExactArgs(0),
		Short: "Log out of MiniMax",
		Long: heredoc.Doc(`
			Log out of MiniMax and remove stored credentials.
			
			This command removes stored authentication tokens for MiniMax
			services. You will need to authenticate again to use MiniMax.
		`),
		Example: heredoc.Doc(`
			# Log out of all MiniMax instances
			$ gh auth minimax logout
			
			# Log out of specific MiniMax instance
			$ gh auth minimax logout --host minimax.global
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(logoutOpts)
		},
	}

	cmd.Flags().StringVarP(&logoutOpts.Host, "host", "H", "", "MiniMax host to log out from")

	return cmd
}

func runLogout(opts *LogoutOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	// Determine which hosts to logout from
	var hosts []string
	if opts.Host != "" {
		hosts = []string{opts.Host}
	} else {
		// Logout from all MiniMax hosts
		hosts = []string{"minimax.global", "minimax.cn"}
	}

	// Check which hosts have active sessions
	var activeHosts []string
	for _, host := range hosts {
		users := cfg.UsersForHost(host)
		if len(users) > 0 {
			activeHosts = append(activeHosts, host)
		}
	}

	if len(activeHosts) == 0 {
		fmt.Fprintf(opts.IO.ErrOut, "You are not logged in to any MiniMax instance.\n")
		return nil
	}

	// Confirm logout
	if opts.IO.CanPrompt() {
		hostsStr := strings.Join(activeHosts, ", ")
		confirmed, err := opts.IO.Prompter.Confirm(
			fmt.Sprintf("Log out of MiniMax? (%s)", hostsStr),
			false,
		)
		if err != nil {
			return err
		}
		if !confirmed {
			return fmt.Errorf("logout cancelled")
		}
	}

	// Perform logout
	cs := opts.IO.ColorScheme()
	for _, host := range activeHosts {
		users := cfg.UsersForHost(host)
		for _, user := range users {
			err := cfg.Logout(host, user)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to log out of %s: %v\n", cs.Red("!"), host, err)
				continue
			}
			
			// Also remove region setting
			cfg.Set(host, "region", "")
			
			fmt.Fprintf(opts.IO.ErrOut, "%s Logged out of %s\n", cs.SuccessIcon(), host)
		}
	}

	return nil
}