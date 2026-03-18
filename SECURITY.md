# Security reporting

Report suspected security issues privately to the repository maintainer. **Do not open a public issue for security vulnerabilities.**

This workspace's pre-commit and CI pipelines run private-key detection (`detect-private-key`) and secret scanning (gitleaks, `.gitleaks.toml`). A committed secret triggers CI failure; rotate the secret and rewrite history if one is exposed — do not rely on deleting the commit alone.

<!-- checkpoint: context(trust-anchors): clarify trust anchors -->
