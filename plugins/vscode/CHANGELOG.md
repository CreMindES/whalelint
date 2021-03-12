# Change Log

## v0.0.7 - 12th March 2021

- Support, i.e exception for special empty image "scratch".
- Bump versions for dependencies.
- Docker image
- Docs grooming

## v0.0.6 - 11th March 2021

- Bug fixes:
  - RUN009 wrong target - apt-get update instead of apt-get install for assume yes flag check |
    [GitHub Issue #53](https://github.com/CreMindES/whalelint/issues/53)
  - Wrong location on RUN009 for bash commands with repeating patterns |
    [GitHub Issue #53](https://github.com/CreMindES/whalelint/issues/53)

## v0.0.5 - 8th March 2021

- Bug fix:
  - Runtime error on ARG without value | [GitHub Issue #46](https://github.com/CreMindES/whalelint/issues/46)

## v0.0.4 - 8th March 2021

- Bug fixes:
  - LSP crash on empty Dockerfile | [GitHub Issue #42](https://github.com/CreMindES/whalelint/issues/42)
  - LSP crash during live editing on RUN command

## v0.0.3 - 8th March 2021

- Bug fixes:
  - false positive on RUN002 for pip install -r/--requirements [file] | [GitHub Issue #36](https://github.com/CreMindES/whalelint/issues/36).
  - false positive on EXP001 in case of arg variables [GitHub Issue #40](https://github.com/CreMindES/whalelint/issues/40)
    - first naive PoC for handling arg variable resolution for EXPOSE commands.

## v0.0.2 - 6th March 2021

- Bug fix:
  - false positive on CPY006 for basic COPY command | [GitHub Issue #34](https://github.com/CreMindES/whalelint/issues/34)

## v0.0.1 - Feb. 2021

- Initial Preview release.