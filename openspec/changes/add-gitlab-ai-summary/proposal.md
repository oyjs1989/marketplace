# Change: Add AI-Powered GitLab Comment Analysis and Summarization

## Why

Currently, Go code reviews generate numerous GitLab comments, replies, and code modifications that reviewers and developers must manually read through to understand the full context. This creates inefficiencies:

- Reviewers spend significant time reading through long comment threads to understand issues
- Developers must manually correlate comments with code changes across multiple files
- Key insights and patterns are buried in lengthy discussions
- No aggregated view of review priorities or common issues

AI-powered analysis and summarization will automatically extract key insights, prioritize issues, and provide concise summaries of review discussions, significantly reducing review overhead and improving code quality.

## What Changes

- Add AI capability to fetch and analyze GitLab merge request comments and code changes
- Implement intelligent summarization of comment threads grouped by topic/file
- Generate priority-based summaries (P0/P1/P2) aligned with FUTU Go standards
- Extract actionable items and unresolved discussions
- Provide aggregated review insights and patterns across the MR
- Integrate with existing Go code review skill for comprehensive analysis
- Support both Chinese and English output for GitLab summaries

## Impact

### Affected Specs
- **NEW**: `gitlab-integration` - Core capability for GitLab API integration and comment analysis
- **MODIFIED**: `go-code-review` (future) - Will integrate with GitLab summary data

### Affected Code
- New skill: `skills/gitlab-ai-summary/SKILL.md`
- New GitLab API integration utilities
- New AI summarization logic for comments and code diffs
- Integration points with existing `go-code-review` skill
- Test cases: `test-cases/gitlab-ai-summary/`

### Key Files
- `skills/gitlab-ai-summary/SKILL.md` - Main skill definition
- `skills/gitlab-ai-summary/shared/GITLAB_INTEGRATION_GUIDE.md` - Integration patterns
- `agents/gitlab-ai-summary/comment-analyzer.md` - Comment analysis agent
- `agents/gitlab-ai-summary/summary-generator.md` - Summary generation agent

## Breaking Changes
None. This is a new capability that doesn't modify existing functionality.
