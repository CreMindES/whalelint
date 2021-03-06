# WhaleLint Rules

<p align="center">
  <img width="500px" src="ruletable/ruletable.png" />
</p>

## Description

WhaleLint has a total of 27 rules at the moment.

Each rule's validation function tries to catch a developer mistake, a bad habbit or advise a better solution.
As such, each of them is assigned one of the common severity levels:
`Error`, `Warning`, `Info`, `Deprecation`.

## Rule List


  - <a href="../../linter/ruleset/cmd001.md">`CMD001`</a> - Prefer JSON notation array format for CMD and ENTRYPOINT
  - <a href="../../linter/ruleset/cpy001.md">`CPY001`</a> - Flag format validation | COPY --[chmod|chown|from]=... srcList... dest|destDir
  - <a href="../../linter/ruleset/cpy002.md">`CPY002`</a> - COPY --chmod=XXXX where XXXX should be a valid permission set value.
  - <a href="../../linter/ruleset/cpy003.md">`CPY003`</a> - COPY chown flag should be in --chown=${USER}:${GROUP} format.
  - <a href="../../linter/ruleset/cpy004.md">`CPY004`</a> - COPY with more than one source requires the destination to end with &#34;/&#34;.
  - <a href="../../linter/ruleset/cpy005.md">`CPY005`</a> - Prefer ADD over COPY for extracting local archives into an image.
  - <a href="../../linter/ruleset/cpy006.md">`CPY006`</a> - COPY --from value should not be the same as the stage.
  - <a href="../../linter/ruleset/ent001.md">`ENT001`</a> - Prefer JSON notation array format for CMD and ENTRYPOINT
  - <a href="../../linter/ruleset/exp001.md">`EXP001`</a> - Expose a valid UNIX port.
  - <a href="../../linter/ruleset/mtr001.md">`MTR001`</a> - MAINTAINER is deprecated. Use a LABEL instead.
  - <a href="../../linter/ruleset/run001.md">`RUN001`</a> - Some bash commands make no sense in an ordinary Docker container.
  - <a href="../../linter/ruleset/run002.md">`RUN002`</a> - Consider pinning versions of packages
  - <a href="../../linter/ruleset/run003.md">`RUN003`</a> - Operators &#34;&amp;&amp;, ||, |&#34; has no affect after semicolon.
  - <a href="../../linter/ruleset/run004.md">`RUN004`</a> - Do not use sudo as it leads to unpredictable behavior. Use a tool like gosu to enforce root.
  - <a href="../../linter/ruleset/run005.md">`RUN005`</a> - Do not upgrade or dist-upgrade the base image
  - <a href="../../linter/ruleset/run006.md">`RUN006`</a> - Clean cache after package manager operation.
  - <a href="../../linter/ruleset/run007.md">`RUN007`</a> - Use &#39;WORKDIR&#39; to switch to a directory.
  - <a href="../../linter/ruleset/run008.md">`RUN008`</a> - Prefer apt-get over apt as the latter does not have a stable CLI.
  - <a href="../../linter/ruleset/run009.md">`RUN009`</a> - Pass -y|--yes|--assume-yes flag to apt-get in order to be headless.
  - <a href="../../linter/ruleset/run010.md">`RUN010`</a> - Pass --no-install-recommends to avoid installing unnecessary packages.
  - <a href="../../linter/ruleset/stl001.md">`STL001`</a> - Stage name alias must be unique.
  - <a href="../../linter/ruleset/sts001.md">`STS001`</a> - Stage name should have an explicit tag..
  - <a href="../../linter/ruleset/sts002.md">`STS002`</a> - Stage name &#34;latest&#34; is prone to future errors.
  - <a href="../../linter/ruleset/sts003.md">`STS003`</a> - Platform should be specified in build tool and not FROM.
  - <a href="../../linter/ruleset/sts004.md">`STS004`</a> - There should only be 1 CMD and/or ENTRYPOINT command.
  - <a href="../../linter/ruleset/usr001.md">`USR001`</a> - Last USER should not be root.
  - <a href="../../linter/ruleset/wkd001.md">`WKD001`</a> - WORKDIR should be an absolute path for clarity and reliability.

## Naming convention:
   - Rule ID

     >3 uppercase letter abbreviation of the Dockerfile AST element and 3 digits
     ```regexp
     [A-Z]{3}[0-9]{3}, e.g. RUN007 or EXP042
     ```

   - Filename of single rule:

     >3 lowercase letter abbreviation of the Dockerfile AST element and 3 digits
     ```regexp
     ruleID.toLower() + ".go", i.e. [a-z]{3}[0-9]{3}.go, e.g. run007.go or exp042.go
     ```

   - ValidationFn name:

     >Validation prefix and the CamelCase version of the Rule ID
     ```regexp
     "Validate" + rule name as [A-Z][A-Z]{2}[0-9]{3}, e.g. ValidateRun007 or ValidateEp042
     ```

   TODO


<br><p align="right">Back to <a href="../../README">README</a></p>