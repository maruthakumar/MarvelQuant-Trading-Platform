# Documentation Style Guide

This document provides a comprehensive style guide for all documentation in the Trading Platform project. It ensures consistency in formatting, terminology, and structure across all documentation files.

This style guide should be used in conjunction with the [Documentation Finalization Plan](documentation_finalization_plan.md), [Documentation Finalization Checklist](documentation_finalization_checklist.md), and [Documentation Template Library](documentation_template_library.md) to ensure all documentation meets the project standards.

## General Formatting Guidelines

### Document Structure

All documentation files should follow this structure:

1. **Title**: Use a single `#` for the main title
2. **Overview**: Brief introduction (1-2 paragraphs) explaining the document's purpose
3. **Table of Contents**: For documents longer than 1000 words
4. **Main Content**: Organized in logical sections with clear headings
5. **References**: Links to related documentation
6. **Revision History**: Table tracking document changes

### Headings

- Use ATX-style headings with # symbols
- Use title case for all headings
- Maximum heading depth: 4 levels
- Include a space after the # symbol
- No trailing # symbols

```markdown
# Level 1 Heading
## Level 2 Heading
### Level 3 Heading
#### Level 4 Heading
```

### Text Formatting

- Use **bold** for emphasis on important terms or concepts
- Use *italics* for slight emphasis or to introduce new terms
- Use `code formatting` for code snippets, file names, and technical terms
- Use > for blockquotes or important notes
- Use horizontal rules (---) to separate major sections

### Lists

- Use - for unordered lists
- Use 1. for ordered lists
- Indent nested lists with 2 spaces
- Include a space after the list marker
- Maintain consistent capitalization in list items

### Code Blocks

- Use triple backticks (```) for code blocks
- Specify the language for syntax highlighting
- Indent code properly within the code block
- Include comments for complex code sections

```go
// This is a Go code example
func example() string {
    return "Hello, World!"
}
```

### Tables

- Use standard markdown table syntax
- Include header row and separator row
- Align columns appropriately (left for text, right for numbers)
- Keep tables simple and readable

```markdown
| Name | Type | Description |
|------|------|-------------|
| id | string | Unique identifier |
| value | number | Numeric value |
```

### Links

- Use descriptive link text
- Use relative links for internal documentation
- Use absolute links for external resources
- Group related links in a References section

```markdown
[Trading Platform Architecture](../architecture/overview.md)
```

## Terminology and Language

### Technical Terminology

Consistent terminology is crucial. Always use the official terms from this glossary:

- **Trading Platform**: The complete system (not "platform" or "system" alone)
- **Module**: A major functional component (not "component" or "section")
- **C++ Execution Engine**: The C++ component for order execution (not "C++ engine" or "execution component")
- **WebSocket**: One word, capital W and S (not "websocket" or "web socket")
- **Backend**: One word (not "back-end" or "back end")
- **Frontend**: One word (not "front-end" or "front end")
- **API**: Application Programming Interface (always use the acronym)
- **SIM**: Simulation environment (always use the acronym)

### Writing Style

- Use present tense
- Use active voice
- Be concise and direct
- Avoid jargon and overly technical language
- Define acronyms on first use
- Use second person ("you") for user guides
- Use third person for technical documentation

### Code Examples

- Include complete, working examples
- Use consistent naming conventions
- Include comments explaining key concepts
- Show both correct usage and common errors
- Include expected output where applicable

## Documentation Types

### User Documentation

- Focus on how to use the system
- Include step-by-step instructions
- Use screenshots for clarity
- Avoid technical implementation details
- Include troubleshooting sections

### Developer Documentation

- Focus on how the system works
- Include architecture diagrams
- Provide detailed API references
- Explain design decisions
- Include code examples

### Project Documentation

- Focus on project status and plans
- Include clear timelines
- Track progress against goals
- Document decisions and their rationale

## File Organization

### File Naming

- Use lowercase for file names
- Use underscores to separate words
- Use descriptive names that indicate content
- Include category prefixes for related files
- Use .md extension for all markdown files

Examples:
- `user_guide.md`
- `api_reference.md`
- `module_1_documentation.md`

### Directory Structure

- Group related documentation in directories
- Use lowercase for directory names
- Use descriptive directory names
- Maintain a flat hierarchy where possible

```
docs/
  ├── user/
  │   ├── user_guide.md
  │   ├── quick_start.md
  │   └── api_reference.md
  ├── developer/
  │   ├── architecture.md
  │   ├── modules/
  │   │   ├── module_1.md
  │   │   └── module_2.md
  │   └── integration.md
  └── project/
      ├── status.md
      ├── plans.md
      └── progress.md
```

## Images and Diagrams

### Image Guidelines

- Use PNG format for screenshots and diagrams
- Use SVG format for vector graphics
- Keep image file sizes under 500KB
- Use descriptive file names
- Include alt text for all images
- Store images in an `images` directory

### Diagram Standards

- Use consistent colors and shapes
- Include a legend for complex diagrams
- Keep diagrams simple and focused
- Use directional arrows to show flow
- Include clear labels

## Document Revision History

All documents must include a revision history table at the end:

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | YYYY-MM-DD | Name | Initial document creation |
| 1.1 | YYYY-MM-DD | Name | Description of changes |

## Document Metadata

Include metadata at the top of each document:

```markdown
---
title: Document Title
version: 1.0
date: YYYY-MM-DD
author: Author Name
category: User/Developer/Project
status: Draft/Review/Approved
---
```

## Conclusion

Following this style guide ensures consistency across all Trading Platform documentation, making it more professional, usable, and maintainable. All documentation contributors should adhere to these guidelines.

## Document Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | April 4, 2025 | Trading Platform Team | Initial document creation |
