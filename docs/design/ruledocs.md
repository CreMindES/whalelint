# Design Decisions

[Docs](../../README.md) > [Design](readme.md)

---

## Rule Docs

A superb documentation and clear examples are essential for lint rules. This needs to be (semi-)automatically
generated with minimal - in an ideal world zero - extra mental effort from the developer.

> 22.11.2020

As of now, a lint rule documentation has two main parts:
- the rule itself, like ID, description, etc.
- the test cases of the rule, which needs to be parsed from the corresponding test file.

Two options were considered for the latter:
- parse from AST
    - pro:
        - test data is defined right inside the test function
        - test data can be in any format, inside an anonymous struct
        - no global variable
    - con:
        - AST parsing can be tedious
        - need additional conventions
- define test data and add it to a global test data map
    - pro:
        - no need to parse AST
        - the data is readily available for consumption
    - con:
        - global test data object
        - testcase structs need to implement a common interface -> no anonymous structs

> 25.11.2020

After implementing both as an at least 80% ready PoC, a better solution was devised:
- Instead of an interface, just mandate, that each rule's test case struct has the following fields, next to any other freely chosen ones:
  - `IsViolation` `bool`
  - `ExampleName` `string`
  - `DocsContext` `string`

Here the last one is the interesting one. It's serves as a template string, and the fields are coming from it's parent, the test case struct itself.
This makes it possible, to have any other field that may partially serve as a test case example, while not duplicating code which would easily lead to code and docs being out-of-sync.

Example:

```go
IsViolation:   false,
ExampleName:   "One stage FROM ubuntu:20.04.",
StageBaseName: "ubuntu:20.04",
DocsContext:   "`FROM` {{ .StageBaseName }}",
```

Here the test docs uses the the base stage's name, while the test docs reuses it for it's `FROM` statement.

This concept is further extenable to the ExampleName field as well. Due to available time and keeping the test names human readable, this is postponed.

> 26.11.2020

The only remaining question is how to build this efficiently, since these anonymus typed test cases only available from test scopes.

> 17.01.2020

Using Go's build tags
