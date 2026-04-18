# Feature Request: Add MiniMax OAuth Provider Support

## Summary

This proposal aims to add support for MiniMax OAuth authentication to GitHub CLI (gh), enabling users to authenticate with MiniMax's AI services directly through the `gh auth` command. This would complement the existing GitHub OAuth flow and provide a unified CLI experience across different AI providers.

## Motivation

### User Benefits
- **Unified CLI Experience**: Users can manage multiple AI provider authentications through a single, familiar CLI interface
- **Consistent Authentication**: Leverage the battle-tested OAuth device code flow already implemented in `gh auth`
- **Developer Productivity**: Seamless integration between GitHub workflows and MiniMax AI services
- **Enterprise Ready**: Follows the same security patterns and storage mechanisms as existing GitHub auth

### Technical Alignment
- **Extensible Architecture**: The current `internal/authflow` package is designed to support multiple OAuth providers
- **Provider Pattern**: GitHub CLI already supports multiple GitHub instances (GitHub.com, GitHub Enterprise), demonstrating the provider pattern
- **OAuth Standards**: MiniMax follows standard OAuth 2.0 device code flow, compatible with existing implementation

## Proposed Implementation

### 1. Core Changes

#### A. Add MiniMax Provider Configuration
**File: `internal/authflow/flow.go`**

```go
// Add MiniMax OAuth configuration
var (
    // The "GitHub CLI" OAuth app (existing)
    oauthClientID = "178c6fc778ccc68e1d6a"
    oauthClientSecret = "34ddeff2b558a23d38fba8a6de74f086ede1cc0b"
    
    // MiniMax OAuth app configuration
    minimaxOAuthClientID = "78257093-7e40-4613-99e0-527b14b39113"
    minimaxOAuthScopes = []string{"openid", "profile", "model.completion"}
    minimaxOAuthBaseURL = "https://api.minimax.io/oauth"
    minimaxOAuthCNBaseURL = "https://api.minimaxi.com/oauth"
)
```

#### B. Extend AuthFlow to Support Multiple Providers
**File: `internal/authflow/flow.go`**

```go
type OAuthProvider struct {
    Name         string
    ClientID     string
    ClientSecret string
    BaseURL      string
    Scopes       []string
}

type ProviderOAuthFlow struct {
    Provider OAuthProvider
    // ... existing flow configuration
}

func NewProviderOAuthFlow(provider OAuthProvider, config *oauth.FlowConfig) *ProviderOAuthFlow {
    // Initialize flow with provider-specific configuration
}

func (f *ProviderOAuthFlow) DetectFlow() (*oauth.Token, error) {
    // Reuse existing device code flow logic with provider-specific URLs
}
```

#### C. Add MiniMax CLI Commands
**File: `pkg/cmd/auth/minimax/`** (new directory)

```
pkg/cmd/auth/minimax/
├── login.go        # gh auth minimax login
├── status.go       # gh auth minimax status
├── logout.go       # gh auth minimax logout
└── refresh.go      # gh auth minimax refresh
```

#### D. Update Root Auth Command
**File: `pkg/cmd/auth/auth.go`**

```go
import (
    authMinimaxCmd "github.com/cli/cli/v2/pkg/cmd/auth/minimax"
    // ... existing imports
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:     "auth <command>",
        Short:   "Authenticate gh and git with GitHub and AI providers",
        GroupID: "core",
    }

    // ... existing commands
    
    // Add MiniMax commands
    cmd.AddCommand(authMinimaxCmd.NewCmdMinimax(f))
    
    return cmd
}
```

### 2. User Interface

#### A. New Command Structure
```bash
gh auth minimax login          # Authenticate with MiniMax
gh auth minimax status         # Check MiniMax auth status
gh auth minimax logout         # Logout from MiniMax
gh auth minimax refresh        # Refresh MiniMax token

# With region support
gh auth minimax login --region global    # Default global MiniMax
gh auth minimax login --region cn        # China MiniMax
```

#### B. Interactive Login Experience
```bash
$ gh auth minimax login
? Which MiniMax region? [Global/China]: Global
? First copy your one-time code: ABC-DEF-GHI
! Press Enter to open https://api.minimax.io/oauth/device/auth in your browser...
✓ Authentication complete.
✓ Configured git protocol
✓ Ready to use MiniMax models
```

### 3. Configuration Storage

#### A. Extend Config Structure
**File: `internal/gh/gh.go`**

```go
type Config interface {
    // ... existing methods
    
    // MiniMax-specific methods
    MiniMaxToken(hostname string) ConfigEntry
    SetMiniMaxToken(hostname, token string)
    MiniMaxRegion(hostname string) ConfigEntry
    SetMiniMaxRegion(hostname, region string)
}
```

#### B. Config File Structure
```yaml
# ~/.config/gh/hosts.yml
github.com:
    user: octocat
    oauth_token: gho_...
    git_protocol: https
minimax.global:
    user: user_id
    oauth_token: mmx_...
    region: global
minimax.cn:
    user: user_id
    oauth_token: mmx_...
    region: cn
```

## Implementation Details

### 1. OAuth Flow Integration
- Reuse existing `internal/authflow.AuthFlow` infrastructure
- Add provider-specific URL handling
- Maintain existing security patterns (PKCE, token validation)
- Support both global and China MiniMax endpoints

### 2. Token Management
- Leverage existing secure storage mechanisms
- Implement refresh token support
- Add token validation and auto-refresh
- Support token revocation

### 3. Error Handling
- Extend existing error mapping for MiniMax-specific errors
- Provide clear user guidance for authentication failures
- Maintain consistent error messaging across providers

### 4. Testing Strategy
- Unit tests for new OAuth provider logic
- Integration tests with mock OAuth server
- End-to-end tests for complete authentication flow
- Compatibility testing with existing GitHub auth

## Security Considerations

### 1. Token Security
- Use existing secure storage mechanisms
- Implement proper token encryption
- Support system credential stores (keychain, etc.)
- Prevent token leakage in logs

### 2. OAuth Security
- Maintain PKCE (Proof Key for Code Exchange) implementation
- Validate OAuth responses
- Implement proper state management
- Support HTTPS-only connections

### 3. User Privacy
- Clear consent prompts
- Transparent data collection notice
- Option to opt-out of telemetry
- Data minimization principles

## Backward Compatibility

### 1. Existing Functionality
- No changes to existing GitHub auth flow
- All current commands remain unchanged
- Existing configurations unaffected
- Gradual migration path for users

### 2. Default Behavior
- `gh auth login` continues to authenticate with GitHub
- MiniMax auth requires explicit command (`gh auth minimax login`)
- No breaking changes to existing APIs
- Optional feature activation

## Success Metrics

### 1. Adoption Metrics
- Number of successful MiniMax authentications
- Usage frequency of MiniMax-specific commands
- User retention rate for MiniMax integration
- Community feedback and issues

### 2. Technical Metrics
- Authentication success rate (>99%)
- Token refresh success rate (>99%)
- API response times (<2s)
- Error rate (<0.1%)

### 3. User Satisfaction
- User survey results
- GitHub star contributions
- Community engagement
- Documentation feedback

## Alternatives Considered

### 1. Extension Approach
- **Pros**: Separate package, less risk to core CLI
- **Cons**: Inconsistent user experience, duplicate code
- **Decision**: Integrate directly for unified experience

### 2. Environment Variable Only
- **Pros**: Simple implementation
- **Cons**: Poor user experience, no secure storage
- **Decision**: Implement full OAuth flow with secure storage

### 3. Third-Party Plugin
- **Pros**: Decoupled from core CLI
- **Cons**: Fragmented ecosystem, maintenance burden
- **Decision**: First-party integration for reliability

## Open Questions

1. **Configuration Management**: Should MiniMax config be stored in the same file as GitHub config?
2. **Command Namespace**: Should we use `gh auth minimax` or `gh minimax auth`?
3. **Default Region**: Should we default to global or ask users to specify?
4. **Telemetry**: What metrics should we collect for MiniMax usage?

## Next Steps

1. **Feedback Phase**: Gather community feedback on this proposal
2. **Design Phase**: Create detailed design documents and mockups
3. **Implementation Phase**: Start with core OAuth flow integration
4. **Testing Phase**: Comprehensive testing and security review
5. **Release Phase**: Gradual rollout with documentation

## Labels

- `enhancement`
- `help wanted` (pending feedback)
- `good first issue` (for individual components)

---

This proposal aims to extend GitHub CLI's authentication capabilities while maintaining the simplicity, security, and reliability that users expect. The implementation leverages existing patterns and infrastructure to minimize risk while maximizing user value.