# Claude Code Marketplace

A curated collection of reusable skills, commands, agents, and workflows for Claude Code.

## Overview

This repository provides a marketplace of pre-built components that extend Claude Code's capabilities. Whether you're looking to automate common tasks, implement best practices, or accelerate your development workflow, you'll find ready-to-use solutions here.

## Project Structure

```
marketplace/
├── skills/          # Individual reusable skills
│   └── go-code-review/     # Go code review skill (v2.0.0)
├── commands/        # Custom slash commands (coming soon)
├── agents/          # Specialized agent configurations (coming soon)
├── workflows/       # Complex multi-step workflows (coming soon)
├── test-cases/      # Test files for validating skills
├── package.json     # Project metadata
└── README.md        # This file
```

## Core Concepts

**Skills**: Domain knowledge and expertise (e.g., coding standards, design patterns)
**Commands**: Quick shortcuts with `/command` syntax
**Agents**: Autonomous task executors with tool access
**Workflows**: Multi-step processes orchestrating the above

Note: Claude Code has built-in support for Skills, Commands, and Agents. Workflows can be implemented through these components.

## Getting Started

### Installation

1. Clone this repository:
```bash
git clone https://gitlab.futunn.com/jasonouyang/marketplace.git
cd marketplace
```

2. Browse the available components in each directory

3. Copy the desired components to your Claude Code project

### Usage

Each component type has its own usage pattern:

- **Skills**: Invoke using skill notation in your Claude Code sessions
- **Commands**: Use with slash notation (e.g., `/command-name`)
- **Agents**: Configure as specialized assistants for specific tasks
- **Workflows**: Follow the documented process for complex operations

## Available Components

### Skills

#### Go Code Review v2.0.0

Comprehensive Go code review with 5 specialized sub-skills and 73+ coding rules.

**Quick Start**:
```
Review my Go code
```

**Features**:
- GORM database operations validation
- Error handling and concurrency safety
- Naming conventions and logging standards
- Code organization and quality checks
- Priority-based issue reporting (P0/P1/P2)

**Documentation**: See [skills/go-code-review/SKILL.md](skills/go-code-review/SKILL.md)

**Test**: Run `Review test-cases/go-code-review/bad/user_service_bad.go` to see it in action.

### Commands
Coming soon. See [commands/README.md](commands/README.md) for planned commands.

### Agents
Coming soon. See [agents/README.md](agents/README.md) for planned agents.

### Workflows
Coming soon. See [workflows/README.md](workflows/README.md) for planned workflows.

## Contributing

Contributions are welcome! To contribute:

1. Fork this repository
2. Create a new branch for your component
3. Add your skill/command/agent/workflow with proper documentation
4. Test thoroughly
5. Submit a merge request

### Contribution Guidelines

- Follow the existing structure and documentation patterns
- Include clear examples and usage instructions
- Test your components before submitting
- Document any dependencies
- Follow best practices for Claude Code development

## Best Practices

When creating or using components from this marketplace:

1. **Documentation**: Always include clear, comprehensive documentation
2. **Testing**: Thoroughly test components before use in production
3. **Dependencies**: Clearly document any external dependencies
4. **Versioning**: Use semantic versioning for tracking changes
5. **Security**: Review code for security implications
6. **Maintenance**: Keep components up-to-date with Claude Code updates

## Support

For questions, issues, or suggestions:

- Open an issue in the repository
- Review existing documentation in each directory
- Check the examples provided with each component

## License

MIT License - feel free to use, modify, and distribute these components.

## Project Status

Active development. New skills, commands, agents, and workflows are being added regularly.

## Testing

Test cases are available in the `test-cases/` directory to validate skill functionality:

- **Go Code Review Tests**: `test-cases/go-code-review/`
  - Comprehensive test files covering P0/P1/P2 violations
  - Good vs bad code examples
  - Detailed test documentation and quick start guide
  - See [test-cases/go-code-review/README.md](test-cases/go-code-review/README.md)

## Roadmap

- [x] Add initial set of common skills (Go Code Review v2.0.0)
- [x] Create testing framework (Go Code Review test cases)
- [ ] Create useful slash commands
- [ ] Develop specialized agents
- [ ] Document complex workflows
- [ ] Build example integrations
- [ ] Add CI/CD pipeline

## Acknowledgments

Built for the Claude Code community to share and reuse valuable development tools.
