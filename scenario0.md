## `tc2mdc`
---
#### `TestInputNil`
> Empty Input
### Parse() returns error on 'nil' input
#### WHEN Convert(nil)
#### THEN error message: 'nil input', data is 'nil'

[top](#tc2mdc)
---
#### `TestInputEmpty`
> Empty Input
### Parse() returns error on 'empty' string input
#### GIVEN Input is ""
#### WHEN Parse("")
#### THEN error message: 'empty input', data is 'nil'

[top](#tc2mdc)
---
#### `TestInputMultiNoComments`
> Empty Input
### Parse() returns empty data on strings without one line comments
#### GIVEN Input contains multi lines without one line comments
#### WHEN Parse()
#### THEN no error, but output data is '<empty>'

[top](#tc2mdc)
---
#### `TestGoPackageName`
> Package, Go
### Parse() returns data with "packageName" on input with package name only
#### GIVEN Input is
- "package somePackage"
#### WHEN Parse()
#### THEN output data has:
- "packageName" = 'somePackage', "title" = '<empty>', "TOC", "Methods" are 'nil'

[top](#tc2mdc)
---
#### `TestGoFuncNameAsMethodName`
> Methods, Go
### Parse() returns data with 1 element in "Methods" on input with 1 test func and 1 non-test func.
#### GIVEN Input is
- "package somePackage"
- "func TestSomething(t *testing.T) {"
- "  if x {"
- "  }" // must be skipped as indented
- "}" // the func end as not indented
- "func DoSomething() {" // non-test function
- "}" // the func end
#### WHEN Parse()
#### THEN output data has:
- "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
- "Methods" contains 1 element: "name" = 'TestSomething', other fields are empty

[top](#tc2mdc)
---
#### `TestGoFuncNameScenario`
> Methods, Go
### Parse() returns data with an element in "Methods" with "Name" and "Scenario"
#### GIVEN Input is
- "package somePackage"
- "func TestSomething(t *testing.T) {"
- "// # Scenario"
- "}" // the func end as not indented
#### WHEN Parse()
#### THEN output is:
- "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
- "Methods" contains 1 element: "name" = 'TestSomething', "scenario" = 'Scenario', other fields are empty

[top](#tc2mdc)
---
#### `TestGoFuncNameTags`
> Methods, Go
### Parse() returns data with an element in "Methods" with multiple tags
#### GIVEN Input is
- "package somePackage"
- "func TestSomething(t *testing.T) {"
- "// > Tag1, Tag2"
- "}" // the func end as not indented
#### WHEN Parse()
#### THEN output is:
- "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
- "Methods" contains 1 element:
  - "name" = 'TestSomething', "scenario" = '<empty>', other fields are empty
  - 2 "tags" = 'Tag1', 'Tag2'

[top](#tc2mdc)
---
#### `TestGoFuncNameGWT`
> Comments, Go
### Parse() returns data with an element in "Methods" with 3 steps - 'GWT' comments
#### GIVEN Input is
- "func TestSomething(t *testing.T) {"
- "// ## GIVEN set"
- "// ## WHEN act"
- "// ## THEN check"
- "}" // the func end as not indented
#### WHEN Parse()
#### THEN output is:
- "Methods" contains 1 element:
  - 3 "steps" :
    - {0, 'GIVEN set'}
    - {0, 'WHEN act'}
    - {0, 'THEN check'}

[top](#tc2mdc)
---
#### `TestGoFuncNameSteps123`
> Comments, Go
### Parse() returns data with an element in "Methods" with 3 steps - comments, common and indented
#### GIVEN Input is
- "func TestSomething(t *testing.T) {"
- "// - common comment"
- "// -- indented comment"
- "// --- indented twice comment"
- "}" // the func end as not indented
#### WHEN Parse()
#### THEN output is:
- "Methods" contains 1 element:
  - 3 "steps" :
    - {1, 'common comment'}
    - {2, 'indented comment'}
    - {3, 'indented2 comment'}

[top](#tc2mdc)
