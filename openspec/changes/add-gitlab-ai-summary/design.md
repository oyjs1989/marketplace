# Design Document: GitLab AI-Powered Comment Analysis

## Context

### Background
FUTU's Go code review process generates extensive GitLab merge request discussions with comments, replies, and code modifications. Reviewers and developers currently must:
- Manually read through all comment threads
- Correlate comments with code changes across multiple files
- Identify priority issues from discussions
- Track resolution status of each discussion

This manual process is time-consuming and prone to missing important feedback, especially in large MRs with 20+ comments across multiple files.

### Constraints
- GitLab API rate limits: 600 requests per minute per user
- Must work with GitLab CE and EE (self-hosted and gitlab.com)
- Claude Code execution environment limitations (no persistent storage between runs)
- Must align with existing FUTU Go code review standards (P0/P1/P2 priorities)
- Output must support Chinese (default for FUTU) and English

### Stakeholders
- **Developers**: Need quick understanding of review feedback
- **Reviewers**: Want to ensure all issues are addressed
- **Team Leads**: Need aggregated view of code quality patterns

## Goals / Non-Goals

### Goals
1. Automatically fetch and analyze GitLab MR comments and diffs
2. Generate intelligent summaries organized by priority and file
3. Identify unresolved discussions and actionable items
4. Integrate with existing Go code review skill for unified analysis
5. Support bilingual output (Chinese default, English optional)
6. Handle large MRs efficiently (100+ comments, 50+ files)

### Non-Goals
1. Direct GitLab comment posting (read-only integration)
2. Real-time monitoring of new comments (snapshot analysis only)
3. GitLab CI/CD pipeline integration (manual trigger via Claude Code)
4. Multi-repository analysis (single MR focus)
5. Historical trend analysis across multiple MRs

## Decisions

### Decision 1: Claude Code Skill Architecture
**Choice**: Implement as a standalone Claude Code skill with orchestrator + specialized agents pattern

**Rationale**:
- Aligns with existing marketplace architecture (see `go-code-review` skill)
- Enables parallel execution of analysis tasks (comment parsing, summarization, integration)
- Maintains modularity for future enhancements
- Reuses proven orchestration patterns from go-code-review skill

**Alternatives Considered**:
- Monolithic skill: Rejected due to complexity and lack of modularity
- External service: Rejected due to deployment overhead and Claude Code integration complexity
- Command-line tool: Rejected because Claude Code skills provide better context awareness

### Decision 2: GitLab API Access Method
**Choice**: Use GitLab REST API v4 with personal access tokens or project tokens

**Rationale**:
- REST API is stable and well-documented
- Personal access tokens provide fine-grained permissions
- Works with both self-hosted and gitlab.com
- No additional dependencies beyond HTTP client (available in Claude Code)

**Alternatives Considered**:
- GraphQL API: More efficient but more complex to implement, less mature
- SSH/Git protocol: Cannot access comments and discussion metadata
- Webhooks: Requires infrastructure, doesn't support historical analysis

### Decision 3: Comment Priority Classification
**Choice**: Use AI-based classification with FUTU Go standards as reference, supplemented by keyword detection

**Rationale**:
- AI can understand context and nuance in natural language comments
- Direct mapping to existing P0/P1/P2 system familiar to developers
- Can learn from explicit priority markers in comments (e.g., "P0:", "critical:")
- Flexible enough to handle both English and Chinese comments

**Classification Logic**:
```
P0 (Critical - Must Fix):
- Keywords: "error", "crash", "security", "nil pointer", "必须", "严重"
- FUTU rules: 1.1.*, 1.2.* (error handling, nil checks)
- Unresolved + blocking comments

P1 (Recommended - Should Fix):
- Keywords: "should", "recommend", "建议", "应该"
- FUTU rules: 2.1.*, 2.2.* (naming, logging)
- Design feedback, architectural concerns

P2 (Suggested - Nice to Have):
- Keywords: "consider", "maybe", "可以考虑"
- FUTU rules: 2.5.*, 3.* (style, formatting)
- Style improvements, optimization suggestions
```

**Alternatives Considered**:
- Keyword-only: Too simplistic, misses context
- Manual priority in comments: Requires reviewer discipline, inconsistent adoption
- Machine learning model: Overkill, requires training data

### Decision 4: Integration with Go Code Review Skill
**Choice**: Loose coupling with optional correlation via shared data format

**Rationale**:
- Both skills can run independently
- When both run, correlation enhances insights
- Shared JSON/markdown format for findings enables correlation
- Avoid tight coupling that makes testing and maintenance difficult

**Integration Points**:
1. **Input**: Both skills analyze same files/MR
2. **Output**: Markdown reports with structured sections
3. **Correlation**: Match line numbers and file paths
4. **Reporting**: Unified report shows: automated finding → GitLab comment status

**Alternatives Considered**:
- Tight integration: Rejected due to complexity and reduced flexibility
- No integration: Rejected because correlation provides significant value
- Go code review depends on GitLab skill: Rejected to maintain independence

### Decision 5: Output Format
**Choice**: Structured Markdown with hierarchical organization (Priority > File > Location)

**Format**:
```markdown
# GitLab Review Summary

## Executive Summary
- Total Comments: 45
- Unresolved: 12
- Priority Breakdown: P0: 3, P1: 8, P2: 1

## P0 Issues (Critical)
### File: services/user_service.go
#### Line 156: Nil pointer check missing
**Comment Thread** (3 messages):
- @reviewer: "这里需要检查 nil 指针" (2024-01-09)
- @developer: "好的,我会修复" (2024-01-09)
- Status: ✅ Resolved

**Action**: Add nil check before dereferencing

## P1 Issues (Recommended)
...

## Actionable Items Checklist
- [ ] services/user_service.go:156 - Add nil pointer check (P0)
- [ ] handlers/auth.go:89 - Improve error message (P1)
...

## Patterns Identified
- **Missing nil checks**: Found in 5 files
  - services/user_service.go:156
  - services/order_service.go:234
  ...
```

**Rationale**:
- Priority-first organization focuses attention on critical issues
- File grouping provides context
- Line numbers enable quick navigation
- Checklist format is actionable
- Pattern section identifies systemic issues

**Alternatives Considered**:
- File-first organization: Less effective for prioritization
- JSON output: Less readable for developers
- Plain text: Lacks structure and formatting

### Decision 6: Bilingual Support Strategy
**Choice**: Template-based approach with language-specific strings, preserve original comment language

**Strategy**:
```
1. Detect output language (default: Chinese for Go reviews)
2. Use language-specific templates for structure
3. Preserve original comment text (quoted)
4. Generate summaries/analysis in target language
```

**Example**:
```markdown
## P0 问题 (严重)
### 文件: services/user_service.go
#### 第156行: 缺少 nil 指针检查
**评论内容**:
> 这里需要检查 nil 指针
**分析**: 在解引用之前必须添加 nil 检查以防止运行时崩溃
```

**Rationale**:
- Maintains authenticity of original reviewer comments
- Provides analysis in preferred language
- Simple to implement and maintain
- Clear visual distinction (quotes for original, plain text for analysis)

**Alternatives Considered**:
- Full translation: Loses nuance, introduces errors
- English only: Doesn't match FUTU workflow
- Dual language side-by-side: Too verbose

## Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│         Claude Code Session (User)                   │
└───────────────────┬─────────────────────────────────┘
                    │ Invokes skill
                    ▼
┌─────────────────────────────────────────────────────┐
│   GitLab AI Summary Skill (Orchestrator)            │
│   - Parse MR URL/ID                                  │
│   - Coordinate agents                                │
│   - Generate final report                            │
└─────┬───────────────────┬───────────────────┬───────┘
      │                   │                   │
      ▼                   ▼                   ▼
┌──────────┐      ┌──────────────┐    ┌─────────────┐
│ GitLab   │      │   Comment    │    │  Summary    │
│ Fetcher  │      │  Analyzer    │    │ Generator   │
│ Agent    │      │  Agent       │    │ Agent       │
└────┬─────┘      └──────┬───────┘    └──────┬──────┘
     │                   │                    │
     │ API calls         │ Analysis           │ AI Summary
     ▼                   ▼                    ▼
┌─────────────────────────────────────────────────────┐
│              GitLab API (REST v4)                    │
│  - Merge Requests                                    │
│  - Discussions / Notes                               │
│  - Diffs / Changes                                   │
└─────────────────────────────────────────────────────┘

Optional Integration:
┌─────────────────────────────────────────────────────┐
│       Go Code Review Skill                           │
│       - Automated static analysis                    │
│       - FUTU standards validation                    │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼ Correlation layer
          ┌────────────────────┐
          │  Unified Report    │
          │  - Automated + MR  │
          │  - Gap analysis    │
          └────────────────────┘
```

## Data Flow

### Phase 1: Data Collection
```
1. User provides MR URL → Parse project ID + MR ID
2. Authenticate with GitLab API (token from env or prompt)
3. Fetch MR metadata (title, description, author, dates)
4. Fetch all discussions (comments + replies)
5. Fetch diffs (changed files, lines, context)
```

### Phase 2: Analysis
```
6. Parse comment threads → Extract metadata
7. Classify comments by priority (P0/P1/P2)
8. Identify resolved vs unresolved
9. Extract actionable items
10. Detect patterns across files
```

### Phase 3: Summarization
```
11. Generate executive summary (metrics)
12. Summarize long threads (>5 messages)
13. Group similar issues (patterns)
14. Create priority-organized report
```

### Phase 4: Output
```
15. Format markdown report
16. Apply language templates (Chinese/English)
17. Generate actionable checklist
18. Save to gitlab_review_summary.md
19. Display in Claude Code session
```

## Risks / Trade-offs

### Risk 1: GitLab API Rate Limiting
**Risk**: Hitting rate limits on large MRs or frequent analysis

**Mitigation**:
- Implement exponential backoff retry logic
- Cache MR data for 5 minutes (within session)
- Batch API requests where possible
- Inform user of rate limit delays

**Trade-off**: Slower analysis vs API rate limit compliance

### Risk 2: AI Summarization Quality
**Risk**: AI may misinterpret comments or generate inaccurate summaries

**Mitigation**:
- Always preserve original comments (quoted) for reference
- Use clear prompting with FUTU standards as context
- Validate summaries against structured data (e.g., resolution status from API)
- Allow users to reference original comments easily

**Trade-off**: Summary conciseness vs accuracy/completeness

### Risk 3: Large MR Performance
**Risk**: Very large MRs (1000+ comments, 100+ files) may be slow or exceed token limits

**Mitigation**:
- Implement pagination and streaming where possible
- Prioritize critical sections (unresolved, P0 issues)
- Provide progress indicators
- Support filtering by file path or priority

**Trade-off**: Comprehensive analysis vs performance

### Risk 4: Authentication and Security
**Risk**: Storing/handling GitLab tokens securely

**Mitigation**:
- Use environment variables (GITLAB_TOKEN) as primary method
- Never log or persist tokens
- Support project-level tokens with minimal permissions (read-only)
- Clear documentation on token creation and permissions

**Trade-off**: Security vs convenience

### Risk 5: Language Detection Accuracy
**Risk**: Mixed-language comments may cause formatting issues

**Mitigation**:
- Use simple heuristics (character set detection)
- Default to Chinese for FUTU projects
- Allow explicit language override
- Preserve original language in all quoted text

**Trade-off**: Automatic detection convenience vs user control

## Migration Plan

N/A - This is a new capability with no existing implementation to migrate from.

### Rollout Plan
1. **Phase 1**: Implement core functionality (GitLab fetch + basic summarization)
2. **Phase 2**: Add priority classification and FUTU standards integration
3. **Phase 3**: Implement Go code review correlation
4. **Phase 4**: Add advanced features (pattern detection, bilingual support)

### Rollback
If issues arise, users can:
- Continue using manual GitLab review
- Use Go code review skill independently
- Report issues for rapid fixes (skill updates are immediate, no deployment required)

## Open Questions

1. **Q**: Should we support GitLab merge request comparison (MR A vs MR B)?
   **A**: Defer to future enhancement. Current focus is single MR analysis.

2. **Q**: Should we analyze commit messages in addition to comments?
   **A**: Yes, include commit messages in the analysis as they provide context for changes.

3. **Q**: How to handle confidential comments or sensitive data in MRs?
   **A**: Skill operates locally in Claude Code session. Data is not persisted or transmitted beyond GitLab API calls. User controls data access via token permissions.

4. **Q**: Should we support MR approval status and merge conflicts in the summary?
   **A**: Yes, include MR status metadata in executive summary (approval status, merge conflicts, pipeline status).

5. **Q**: How to handle very old MRs with 100+ comments spanning weeks?
   **A**: Default to most recent comments first. Support date filtering in future enhancement.

## Testing Strategy

### Unit Testing
- Mock GitLab API responses for various scenarios
- Test priority classification logic with sample comments
- Validate markdown generation with edge cases

### Integration Testing
- Test with real GitLab MRs (small, medium, large)
- Validate bilingual output with mixed-language comments
- Test integration with Go code review skill

### Performance Testing
- Measure analysis time for MRs of varying sizes:
  - Small: <10 comments, <5 files
  - Medium: 10-50 comments, 5-20 files
  - Large: 50+ comments, 20+ files
- Set performance targets:
  - Small MRs: <5 seconds
  - Medium MRs: <30 seconds
  - Large MRs: <2 minutes

### User Acceptance Testing
- Test with real FUTU Go code review MRs
- Validate summary quality and accuracy
- Gather feedback on output format and priorities
