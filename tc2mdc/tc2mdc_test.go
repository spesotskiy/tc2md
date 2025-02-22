package tc2mdc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputNil(t *testing.T) {
	// > Empty Input
	// # Convert() returns error on 'nil' input
	// ## WHEN Convert(nil)
	output, err := Convert(nil)
	// ## THEN error message: 'nil input'
	require.Empty(t, output, "msg must be empty")
	require.ErrorContains(t, err, "nil input")
}

func TestInputEmpty(t *testing.T) {
	// > Empty Input
	// # Convert() returns error on empty string input
	// ## GIVEN Input is ""
	var input = []string{""}
	// ## WHEN Convert("")
	output, err := Convert(input)
	// ## THEN error message: 'empty input'
	require.Empty(t, output, "msg must be empty")
	require.ErrorContains(t, err, "empty input")
}

func TestInputMultiNoComments(t *testing.T) {
	// > Empty Input
	// # Convert() returns empty output on strings without one line comments
	// ## GIVEN Input contains multi lines without one line comments
	var input = []string{" something", "", " something else"}
	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is '<empty>'
	require.Empty(t, output, "output must be empty")
	require.Empty(t, err, "must be no error")
}

func TestInputContainsOLCbutNoMarkers(t *testing.T) {
	// > Empty Input
	// # Convert() returns empty output on input with OLC but not one space between MD markers
	// ## GIVEN Input is:
	var input = []string{
		// - "// Scenario"   (no MD marker)
		OLC + " Scenario",
		// - "//- point"     (no space after OLC)
		OLC + "- point",
		// - "//  - point"   (two spaces after OLC)
		OLC + "  - point",
		// - "// #header"    (no space after MD marker)
		OLC + " #header",
		// - "// ##  header" (two spaces after MD marker)
		OLC + " ##  header",
	}
	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is '<empty>'
	require.Empty(t, err, "must be no error")
	require.Empty(t, output, "output must be empty")
}

func TestInputScenarioHeader(t *testing.T) {
	// > MD Markers
	// # Convert() returns header(3) line on input started with '#'
	// ## GIVEN Input is
	// - "// # Scenario"
	var input = []string{OLC + " # Scenario"}
	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output:
	// - "### Scenario"
	require.Empty(t, err, "must be no error")
	require.Equal(t, []string{"### Scenario"}, output)
}

func TestInputStepHeader(t *testing.T) {
	// > MD Markers
	// # Convert() returns header(4) line on input started with '##'
	// ## GIVEN Input is
	// - "// ## GIVEN"
	var input = []string{OLC + " ## GIVEN"}
	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output:
	// - "#### GIVEN"
	require.Empty(t, err, "must be no error")
	require.Equal(t, []string{"#### GIVEN"}, output)
}

func TestInputBulletNote(t *testing.T) {
	// > MD Markers
	// # Convert() returns line as is on input started with '-' and '>'
	// ## GIVEN Input is:
	var input = []string{
		// - "// > Group"
		OLC + " > Group",
		// - "// - Step"
		OLC + " - Step ",
	}
	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Equal(
		t,
		[]string{
			// - "> Group"
			"> Group",
			// - "- Step"
			"- Step",
		},
		output)
}

func TestInputBullet2(t *testing.T) {
	// > MD Markers
	// # Convert() returns shifted (2, 4 spaces) line on input started with '--', '---'
	// ## GIVEN Input is
	var input = []string{
		// - "// -- Step 2"
		OLC + " -- Step2 ",
		// - "// --- Step3"
		OLC + " --- Step3 ",
	}
	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Equal(t, []string{
		// - "__- Step2"
		"  - Step2",
		// - "____- Step3"
		"    - Step3",
	}, output)
}

func TestGoPackageNameAsHeader(t *testing.T) {
	// > Header, Go
	// # Convert() returns package header on input with package name
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
	}

	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Equal(t, []string{
		// - "## `somePackage`"
		"## `somePackage`",
	}, output)
}

func TestGoFuncNameAsHeader(t *testing.T) {
	// > Header, Go
	// # Convert() returns scenario header with separator and link to the top on input with test func name
	// ## GIVEN Input is
	var input = []string{
		// - "func TestSomething(t *testing.T) {"
		"func TestSomething(t *testing.T) {",
		// - "  if x {"
		"  if x {",
		// - "  }" // must be skipped as indented
		"  }",
		// - "}" // the func end as not indented
		"}",
	}

	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Equal(t, []string{
		// - "---"
		"---",
		// - "#### `TestSomething`"
		"#### `TestSomething`",
		"",
		// - "[top]#top" - link to the top
		"[top](#top)",
	}, output)
}

func TestGoTwoTestsWithPackage(t *testing.T) {
	// > Header, Go
	// # Convert() returns 2 test headers on input with 2 tests
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
		// - "func TestSomething1(t *testing.T) {"
		"func TestSomething1(t *testing.T) {",
		OLC + " # Scenario1",
		"}",
		// - "func TestSomething2(t *testing.T) {"
		"func TestSomething2(t *testing.T) {",
		OLC + " # Scenario2",
		"}",
	}

	// ## WHEN Convert()
	output, err := Convert(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Equal(t, []string{
		// - "## `somePackage`"
		"## `somePackage`",
		// - "#### `TestSomething1`"
		"---",
		"#### `TestSomething1`",
		"### Scenario1",
		// - "[top]#somePackage" - link to the package line
		"",
		"[top](#somePackage)",
		// - "#### `TestSomething2`"
		"---",
		"#### `TestSomething2`",
		"### Scenario2",
		// - "[top]#somePackage" - link to the package line
		"",
		"[top](#somePackage)",
	}, output)
}
