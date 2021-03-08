# Change Log

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