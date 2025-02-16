> Empty Input
### Convert() returns error on 'nil' input
#### WHEN Convert(nil)
#### THEN error message: 'nil input'
> Empty Input
### Convert() returns error on empty string input
#### GIVEN Input is ""
#### WHEN Convert("")
#### THEN error message: 'empty input'
> Empty Input
### Convert() returns empty output on strings without one line comments
#### GIVEN Input contains multi lines without one line comments
#### WHEN Convert()
#### THEN output is '<empty>'
> Empty Input
### Convert() returns empty output on input with OLC but not one space between MD markers
#### GIVEN Input is:
- "// Scenario"   (no MD marker)
- "//- point"     (no space after OLC)
- "//  - point"   (two spaces after OLC)
- "// #header"    (no space after MD marker)
- "// ##  header" (two spaces after MD marker)
#### WHEN Convert()
#### THEN output is '<empty>'
> MD Markers
### Convert() returns header(3) line on input started with '#'
#### GIVEN Input is
- "// # Scenario"
#### WHEN Convert()
#### THEN output:
- "### Scenario"
> MD Markers
### Convert() returns header(4) line on input started with '##'
#### GIVEN Input is
- "// ## GIVEN"
#### WHEN Convert()
#### THEN output:
- "#### GIVEN"
> MD Markers
### Convert() returns line as is on input started with '-' and '>'
#### GIVEN Input is:
- "// > Group"
- "// - Step"
#### WHEN Convert()
#### THEN output is:
- "> Group"
- "- Step"
> MD Markers
### Convert() returns shifted (2, 4 spaces) line on input started with '--', '---'
#### GIVEN Input is
- "// -- Step 2"
- "// --- Step3"
#### WHEN Convert()
#### THEN output is:
- "__- Step2"
- "____- Step3"
