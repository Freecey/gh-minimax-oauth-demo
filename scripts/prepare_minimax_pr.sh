#!/bin/bash

# Script to prepare and submit MiniMax OAuth feature to GitHub CLI
# This automates the process of creating a feature branch and PR

set -e

echo "🚀 Preparing MiniMax OAuth Feature for GitHub CLI"
echo "=============================================="

# Configuration
REPO="cli/cli"
BRANCH_NAME="feature/minimax-oauth"
PR_TITLE="feat: Add MiniMax OAuth Provider Support"
PR_BODY_FILE="/tmp/gh-cli-minimax-oauth-pr.md"
ISSUE_BODY_FILE="/tmp/gh-cli-minimax-oauth-issue.md"
SCRIPT_DIR="/tmp"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    local color=$1
    shift
    local message=$@
    echo -e "${color}${message}${NC}"
}

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    print_color $RED "❌ Error: This is not a git repository"
    print_color $YELLOW "Please run this script from the GitHub CLI root directory or fork"
    exit 1
fi

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    print_color $RED "❌ Error: GitHub CLI (gh) is not installed"
    print_color $YELLOW "Please install gh CLI first: https://cli.github.com/"
    exit 1
fi

# Check if we're authenticated with gh
if ! gh auth status &> /dev/null; then
    print_color $RED "❌ Error: Not authenticated with GitHub CLI"
    print_color $YELLOW "Please run 'gh auth login' first"
    exit 1
fi

print_color $GREEN "✅ Environment checks passed"

# Check current remote
current_remote=$(git remote get-url origin 2>/dev/null || echo "")
if [[ $current_remote != *"$REPO"* ]]; then
    print_color $YELLOW "⚠️  Current remote is: $current_remote"
    print_color $YELLOW "Make sure you're working with the correct repository"
fi

# Create and switch to feature branch
print_color $BLUE "🌿 Creating feature branch: $BRANCH_NAME"

if git show-ref --verify --quiet refs/heads/"$BRANCH_NAME"; then
    print_color $YELLOW "⚠️  Branch $BRANCH_NAME already exists"
    read -p "Do you want to use existing branch? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_color $RED "❌ Aborted"
        exit 1
    fi
    git checkout "$BRANCH_NAME"
else
    git checkout -b "$BRANCH_NAME"
    print_color $GREEN "✅ Created and switched to branch: $BRANCH_NAME"
fi

# Ensure we have the latest changes
print_color $BLUE "📥 Fetching latest changes..."
git fetch origin

# Create necessary directories
print_color $BLUE "📁 Creating directory structure..."
mkdir -p pkg/cmd/auth/minimax

# Copy the implementation files
print_color $BLUE "📄 Copying implementation files..."

# Copy auth flow extension
if [ -f "$SCRIPT_DIR/minimax_auth_flow.go" ]; then
    cp "$SCRIPT_DIR/minimax_auth_flow.go" internal/authflow/minimax_auth_flow.go
    print_color $GREEN "✅ Copied auth flow extension"
else
    print_color $RED "❌ Missing file: minimax_auth_flow.go"
    exit 1
fi

# Copy command files
for cmd_file in login status logout refresh; do
    if [ -f "$SCRIPT_DIR/minimax_$cmd_file.go" ]; then
        cp "$SCRIPT_DIR/minimax_$cmd_file.go" pkg/cmd/auth/minimax/$cmd_file.go
        print_color $GREEN "✅ Copied $cmd_file command"
    else
        print_color $RED "❌ Missing file: minimax_$cmd_file.go"
        exit 1
    fi
done

# Create the main minimax command file
cat > pkg/cmd/auth/minimax/minimax.go << 'EOF'
package minimax

import (
	"github.com/cli/cli/v2/pkg/cmd/auth/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	IO              *iostreams.IOStreams
	Config          func() (gh.Config, error)
	HttpClient      func() (*http.Client, error)
	PlainHttpClient func() (*http.Client, error)
	Prompter        shared.Prompt
	Browser         browser.Browser
}

// NewCmdMinimax creates the root command for MiniMax authentication
func NewCmdMinimax(f *cmdutil.Factory, opts *LoginOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minimax <command>",
		Short: "Authenticate with MiniMax AI services",
		Long: `Authenticate with MiniMax AI services using OAuth.

This command group allows you to manage MiniMax authentication
through GitHub CLI's familiar interface.`,
	}

	return cmd
}
EOF

# Update the main auth command to include minimax
print_color $BLUE "🔧 Updating main auth command..."

# Backup the original file
cp pkg/cmd/auth/auth.go pkg/cmd/auth/auth.go.backup

# Add minimax import and command
# This is a simple approach - in practice, you'd need more careful patching
if ! grep -q "authMinimaxCmd" pkg/cmd/auth/auth.go; then
    # Add import
    sed -i '/import (/a\\
	authMinimaxCmd "github.com/cli/cli/v2/pkg/cmd/auth/minimax"' pkg/cmd/auth/auth.go
    
    # Add command
    sed -i '/cmd.AddCommand(authTokenCmd.NewCmdToken(f, nil))/a\\
	cmd.AddCommand(authMinimaxCmd.NewCmdMinimax(f))' pkg/cmd/auth/auth.go
    
    print_color $GREEN "✅ Updated main auth command"
else
    print_color $YELLOW "⚠️  Main auth command already updated"
fi

# Add files to git
print_color $BLUE "📤 Adding files to git..."
git add internal/authflow/minimax_auth_flow.go
git add pkg/cmd/auth/minimax/
git add pkg/cmd/auth/auth.go

# Create commit
print_color $BLUE "💾 Creating commit..."
git commit -m "feat: Add MiniMax OAuth provider support

- Add MiniMax OAuth device code flow implementation
- Support for global and China regions
- Add login, status, logout, and refresh commands
- Integrate with existing auth infrastructure
- Maintain backward compatibility

This enables users to authenticate with MiniMax AI services
using the familiar GitHub CLI interface."

# Check for test compilation
print_color $BLUE "🧪 Testing compilation..."
if command -v go &> /dev/null; then
    if go build ./cmd/gh; then
        print_color $GREEN "✅ Build successful"
    else
        print_color $RED "❌ Build failed"
        print_color $YELLOW "Please fix compilation errors before proceeding"
        exit 1
    fi
else
    print_color $YELLOW "⚠️  Go not available - skipping build test"
fi

# Create PR
print_color $BLUE "🚀 Creating Pull Request..."

# Check if PR already exists
if gh pr view --json title | grep -q "$PR_TITLE"; then
    print_color $YELLOW "⚠️  PR with this title already exists"
    read -p "Do you want to view existing PR? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        gh pr view --web
        exit 0
    fi
else
    # Push branch
    print_color $BLUE "📤 Pushing branch to remote..."
    git push -u origin "$BRANCH_NAME"
    
    # Create PR
    if [ -f "$PR_BODY_FILE" ]; then
        gh pr create --title "$PR_TITLE" --body "$(cat $PR_BODY_FILE)" --web
    else
        gh pr create --title "$PR_TITLE" --body "Add MiniMax OAuth provider support to GitHub CLI" --web
    fi
    
    print_color $GREEN "✅ Pull Request created successfully!"
fi

# Summary
print_color $BLUE "📋 Summary"
print_color $BLUE "========="
print_color $GREEN "✅ Feature branch created: $BRANCH_NAME"
print_color $GREEN "✅ Implementation files added"
print_color $GREEN "✅ Integration completed"
print_color $GREEN "✅ Pull Request created"

print_color $YELLOW ""
print_color $YELLOW "Next steps:"
print_color $YELLOW "1. Monitor the PR for review comments"
print_color $YELLOW "2. Address any feedback from maintainers"
print_color $YELLOW "3. Ensure all tests pass"
print_color $YELLOW "4. Wait for merge approval"

print_color $GREEN ""
print_color $GREEN "🎉 MiniMax OAuth feature submitted for review!"
print_color $GREEN "Thank you for contributing to GitHub CLI!"