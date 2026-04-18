# MiniMax OAuth Integration for GitHub CLI

## 🎯 Description

This repository demonstrates a complete implementation of MiniMax OAuth support for GitHub CLI. This is a proof-of-concept showing how to add MiniMax AI platform authentication to the official GitHub CLI tool.

## 🚀 Features

- ✅ **Complete OAuth 2.0 Flow**: Device Code Flow implementation
- ✅ **Multi-region Support**: Global (`api.minimax.io`) and China (`api.minimaxi.com`)
- ✅ **Full Command Set**: `login`, `status`, `logout`, `refresh`
- ✅ **Security**: OAuth 2.0 + PKCE implementation
- ✅ **Production Ready**: Follows GitHub CLI patterns and conventions

## 📋 Usage

### Basic Usage
```bash
# Login with MiniMax (global)
gh auth login minimax --hostname api.minimax.io

# Login with MiniMax China
gh auth login minimax --hostname api.minimaxi.com --region cn

# Check authentication status
gh auth status

# Logout from MiniMax
gh auth logout --hostname api.minimax.io
```

## 🔧 Technical Details

### OAuth Configuration
- **Client ID**: `78257093-7e40-4613-99e0-527b14b39113`
- **Scopes**: `openid profile model.completion`
- **Flow**: Device Code Flow
- **Regions**: `global` (default) and `cn` (China)

### Endpoints
- **Global**: 
  - Device Auth: `https://api.minimax.io/v1/oauth/device_authorize`
  - Token: `https://api.minimax.io/v1/oauth/token`

- **China**: 
  - Device Auth: `https://api.minimaxi.com/v1/oauth/device_authorize`
  - Token: `https://api.minimaxi.com/v1/oauth/token`

## 🏗️ Integration Architecture

This implementation follows the existing GitHub CLI patterns:

```
pkg/cmd/auth/minimax/    # Command implementations
├── minimax_login.go     # Login command
├── minimax_status.go    # Status command  
├── minimax_logout.go    # Logout command
└── minimax_refresh.go   # Refresh command

internal/authflow/minimax/  # OAuth flow implementation
└── minimax_auth_flow.go   # OAuth 2.0 device code flow
```

## 🧪 Testing

The implementation includes comprehensive tests:
- Unit tests for all OAuth flows
- Integration tests with real endpoints
- Regression tests for existing providers
- Manual testing procedures

## 📦 Installation

This is a demonstration implementation. For actual use, the code should be integrated into the official GitHub CLI repository.

## 🔒 Security

- Follows OAuth 2.0 RFC standards
- Implements PKCE for enhanced security
- Uses HTTPS endpoints only
- Secure token storage following existing patterns

## 🤝 Contributing

This code is intended as a reference implementation for contributing MiniMax OAuth support to the official GitHub CLI project.

## 📄 License

This demonstration follows the same license as the original GitHub CLI project.
