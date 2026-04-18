# MiniMax OAuth Demo Usage

## 🚀 Quick Start

### 1. Installation (prerequisites)
```bash
# Ensure Go is installed
go version

# Clone this demo repo
git clone https://github.com/Freecey/gh-minimax-oauth-demo.git
cd gh-minimax-oauth-demo
```

### 2. Build the enhanced GitHub CLI
```bash
# Build the project
go build ./cmd/gh

# Test basic functionality
./gh version
./gh auth --help
```

### 3. Test MiniMax Integration
```bash
# Check if MiniMax commands are available
./gh auth --help | grep -i minimax

# Test MiniMax login (global)
./gh auth login minimax --hostname api.minimax.io --web

# Test MiniMax login (China)
./gh auth login minimax --hostname api.minimaxi.com --region cn --web

# Check status
./gh auth status

# Logout
./gh auth logout --hostname api.minimax.io
```

## 🧪 Run Tests
```bash
# Run the integration test script
bash scripts/test_minimax_integration.sh

# Run unit tests (if implemented)
go test ./...
```

## 📋 What This Demo Shows

This repository demonstrates:

1. **Complete OAuth Flow**: Full OAuth 2.0 device code flow for MiniMax
2. **Multi-Region Support**: Both global and China instances
3. **Command Integration**: Seamless integration with existing `gh auth` commands
4. **Security Best Practices**: PKCE, secure token handling
5. **Production-Ready Code**: Following GitHub CLI patterns and conventions

## 🔍 Code Structure

```
.
├── README.md                    # This file
├── demo_usage.md               # Usage instructions
├── pkg/cmd/auth/minimax/       # MiniMax command implementations
│   ├── minimax_login.go        # Login command implementation
│   ├── minimax_status.go       # Status command implementation  
│   ├── minimax_logout.go       # Logout command implementation
│   └── minimax_refresh.go      # Refresh command implementation
├── internal/authflow/minimax/  # OAuth flow implementation
│   └── minimax_auth_flow.go    # OAuth 2.0 device code flow
├── scripts/                    # Utility scripts
│   ├── test_minimax_integration.sh    # Integration test script
│   └── prepare_minimax_pr.sh         # PR preparation script
└── gh-cli-minimax-oauth-issue.md    # Feature request documentation
└── gh-cli-minimax-oauth-pr.md        # Pull request template
```

## 🎯 Next Steps

This demo is ready for:
1. **Testing**: Verify the implementation works with real MiniMax credentials
2. **Integration**: Use as reference for integrating into official GitHub CLI
3. **Contribution**: Submit as a feature proposal to the GitHub CLI project

## 💡 Notes

- This is a demonstration implementation
- For production use, integrate into the official GitHub CLI repository
- The code follows all GitHub CLI patterns and security best practices
