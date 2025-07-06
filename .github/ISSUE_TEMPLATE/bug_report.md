name: Bug Report
description: File a bug report to help us improve
title: "[BUG] "
labels: ["bug", "triage"]
assignees: []

body:
  - type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this bug report! 🐛

  - type: input
    id: version
    attributes:
      label: Version
      description: What version of SoxyChecker GUI are you running?
      placeholder: e.g., v1.0.0
    validations:
      required: true

  - type: dropdown
    id: os
    attributes:
      label: Operating System
      description: What operating system are you using?
      options:
        - Windows 10
        - Windows 11
        - macOS (Intel)
        - macOS (Apple Silicon)
        - Linux (Ubuntu)
        - Linux (Other)
        - Other
    validations:
      required: true

  - type: textarea
    id: description
    attributes:
      label: Bug Description
      description: A clear and concise description of what the bug is.
      placeholder: Describe the bug...
    validations:
      required: true

  - type: textarea
    id: steps
    attributes:
      label: Steps to Reproduce
      description: Steps to reproduce the behavior
      placeholder: |
        1. Go to '...'
        2. Click on '...'
        3. Scroll down to '...'
        4. See error
    validations:
      required: true

  - type: textarea
    id: expected
    attributes:
      label: Expected Behavior
      description: A clear and concise description of what you expected to happen.
      placeholder: What should have happened?
    validations:      required: true

  - type: textarea
    id: actual
    attributes:
      label: Actual Behavior
      description: A clear and concise description of what actually happened.
      placeholder: What actually happened?
    validations:
      required: true

  - type: textarea
    id: screenshots
    attributes:
      label: Screenshots
      description: If applicable, add screenshots to help explain your problem.
      placeholder: Drag and drop screenshots here or paste image URLs

  - type: textarea
    id: logs
    attributes:
      label: Logs
      description: Please copy and paste any relevant log output. This will be automatically formatted into code, so no need for backticks.
      render: shell

  - type: textarea
    id: additional
    attributes:
      label: Additional Context
      description: Add any other context about the problem here.
      placeholder: Any additional information that might be helpful...

  - type: checkboxes
    id: terms
    attributes:
      label: Code of Conduct
      description: By submitting this issue, you agree to follow our Code of Conduct
      options:
        - label: I agree to follow this project's Code of Conduct
          required: true
