# Whalelint <img width="13px" src="logo.svg"> <img align="right" style="position: relative; top: 10px;" src="https://github.com/cremindes/whalelint/workflows/build/badge.svg" />

Dockerfile linter written in Go

*Disclaimer: this is a pet-project while I'm learning Golang.*

<p align="center">
  <img width="550px" src="illustration.svg"> 
</p>

## Feature list

- [x] extendable ruleset
- [ ] cli
- [x] output as json
- [ ] rule escaping per line
- [ ] output coloring
- [ ] config file

## Rules

Each Dockerfile AST element has a corresponding set of rules. 

### Naming convention:
- Rule ID: 
  ```3 uppercase letters abbreviation of the Dockerfile AST element and then 3 digits```
  ```regexp
  [A-Z]{3}[0-9]{3}, e.g. RUN007 or EXP042
  ```
- Filename of single rule: 
  ```3 lowercase letter abbreviation of the Dockerfile AST element and then 3 digits```
  ```regexp 
  ruleID.toLower() + ".go", i.e. [a-z]{3}[0-9]{3}.go, e.g. run007.go or exp042.go
  ```
- ValidationFn name: 
  ```Validation prefix and then the CamelCase version of the Rule ID```
  ```regexp
  "Validate" + rule name as [A-Z][A-Z]{2}[0-9]{3}, e.g. ValidateRun007 or ValidateEp042
  ```

TODO

## Sample output

TODO

## Plugins

- JetBrains
- VSCode