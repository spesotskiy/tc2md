## `tc2mdc`
---
#### `TestWriteNil`
> Write to MD
### Write() returns 'nil' data as 'nil'
#### WHEN Write(nil)
#### THEN - MD text is 'nil'

[top](#tc2mdc)
---
#### `TestWritePackageName`
> Write to MD
### Write() returns "packageName" as header(2)
#### GIVEN - testData: "packageName" = 'SomePackage'
#### WHEN Write()
#### THEN - MD text includes 1 line:
- "## `SomePackage`"

[top](#tc2mdc)
---
#### `TestWriteOneMethod`
> Write to MD
### Write() returns test method name as header(4) with separator before name and a link to the top after
#### GIVEN - testData: "packageName" = ''
- 1 element in "methods": "name" = 'TestSomething'
#### WHEN Write()
#### THEN - MD text includes 4 lines:
- "---"
- "#### `TestSomething`"
- "" // to separate the link and the last step
- "[top]#top" - link to the top

[top](#tc2mdc)
---
#### `TestWritePackageTwoMethods`
> Write to MD
### Write() returns test method name as header(4) with separator before name and a link to the top after
#### GIVEN - testData: "packageName" = 'somePackage'
- 2 elements in "methods": "name" = 'TestSomething1' and 'TestSomething2'
#### WHEN Write()
#### THEN - MD text includes 9 lines - 1 package line and 2 sets of 4 lines per a test:
- "## `somePackage`"
- "---"
- "#### `TestSomething1`"
- "", "[top]#somePackage" - link to the package line
- "---"
- "#### `TestSomething2`"
- "", "[top]#somePackage" - link to the package line

[top](#tc2mdc)
---
#### `TestWriteMethodTags`
> Write to MD
### Write() returns "tags" as one line - comma separated list prefixed with '>'
#### GIVEN - testData: "packageName" = ''
- 1 element in "methods": "name" = 'TestSomething'
- two "tags": "Tag", "Complex Tag"
#### WHEN Write()
#### THEN - MD text includes one tags line after method name line:
- "#### `TestSomething`"
- "> Tag, Complex Tag"

[top](#tc2mdc)
---
#### `TestWriteMethodScenario`
> Write to MD
### Write() returns "scenario" as header(3) line
#### GIVEN - testData: "packageName" = ''
- 1 element in "methods": "name" = 'TestSomething'
- "scenario" = 'Something happens'
#### WHEN Write()
#### THEN - MD text includes a scenario line after method name line:
- "#### `TestSomething`"
- "### Something happens"

[top](#tc2mdc)
---
#### `TestWriteMethodGWT`
> Write to MD
### Write() returns 'GWT' comments as header(4) lines
#### GIVEN - testData: "packageName" = ''
- 1 element in "methods": "name" = 'TestSomething'
- 3 steps of 'GWT' kind: "GIVEN set", "WHEN act", "THEN check"
#### WHEN Write()
#### THEN - MD text includes 3 'GWT' step lines as headers(4):
- "#### `TestSomething`"
- "#### GIVEN set"
- "#### WHEN act"
- "#### THEN check"

[top](#tc2mdc)
---
#### `TestWriteMethodIndentedSteps`
> Write to MD
### Write() returns indented comments as indented lines prefixed with '-'
#### GIVEN - testData: "packageName" = ''
- 1 element in "methods": "name" = 'TestSomething'
- 3 steps: common 'Step1', indented 'Step2', indented twice 'Step3'
#### WHEN Write()
#### THEN - MD text includes 3 indented step lines of 3 levels:
- "#### `TestSomething`"
- "- Step1"
- "__- Step2"
- "____- Step3"

[top](#tc2mdc)
