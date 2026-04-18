package main

import (
	"fmt"
	"os"

	"github.com/Freecey/gh-minimax-oauth-demo/pkg/cmd/auth/minimax"
	"github.com/spf13/cobra"
)

func main() {
	// Commande racine
	var rootCmd = &cobra.Command{
		Use:   "gh",
		Short: "GitHub CLI with MiniMax OAuth Demo",
		Long: `GitHub CLI enhanced with MiniMax OAuth support.
This is a demo implementation showing MiniMax OAuth integration.`,
	}

	// Ajouter la commande auth
	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with GitHub and MiniMax",
		Long:  "Manage authentication with GitHub and MiniMax AI platform",
	}

	// Ajouter les commandes MiniMax
	authCmd.AddCommand(minimax.NewLoginCmd())
	authCmd.AddCommand(minimax.NewStatusCmd()) 
	authCmd.AddCommand(minimax.NewLogoutCmd())
	authCmd.AddCommand(minimax.NewRefreshCmd())

	rootCmd.AddCommand(authCmd)

	// Exécuter
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}