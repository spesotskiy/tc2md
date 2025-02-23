package tc2mdc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputNil(t *testing.T) {
	// > Empty Input
	// # Convert() returns error on 'nil' input
	// ## WHEN Convert(nil)
	testData, err := Parse(nil)
	// ## THEN error message: 'nil input', data is empty
	require.Empty(t, testData, "data must be empty")
	require.ErrorContains(t, err, "nil input")
}

func TestInputEmpty(t *testing.T) {
	// > Empty Input
	// # Parse() returns error on empty string input
	// ## GIVEN Input is ""
	var input = []string{""}
	// ## WHEN Parse("")
	testData, err := Parse(input)
	// ## THEN error message: 'empty input', data is empty
	require.Empty(t, testData, "data must be empty")
	require.ErrorContains(t, err, "empty input")
}

func TestInputMultiNoComments(t *testing.T) {
	// > Empty Input
	// # Parse() returns empty data on strings without one line comments
	// ## GIVEN Input contains multi lines without one line comments
	var input = []string{" something", "", " something else"}
	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output data is '<empty>'
	require.Empty(t, err, "must be no error")
	require.Empty(t, testData, "output must be empty")
}

func TestGoPackageNameAsHeader(t *testing.T) {
	// > Header, Go
	// # Parse() returns data with "packageName" on input with package name only
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
	}

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output data has:
	require.Empty(t, err, "must be no error")
	// - "packageName" = 'somePackage', ("title", "TOC", "Methods") = '<empty>'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "must be empty title")
	require.Empty(t, testData.toc, "must be empty TOC")
	require.Empty(t, testData.methods, "must be empty methods")
}

func TestGoFuncNameAsMethodName(t *testing.T) {
	// > Header, Go
	// # Parse() returns data with 1 element in "Methods" on input with 1 test func and 1 non-test func.
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
		// - "func TestSomething(t *testing.T) {"
		"func TestSomething(t *testing.T) {",
		// - "  if x {"
		"  if x {",
		// - "  }" // must be skipped as indented
		"  }",
		// - "}" // the func end as not indented
		"}",
		// - "func DoSomething() {" // non-test function
		"func DoSomething() {",
		// - "}" // the func end
		"}",
	}

	// ## WHEN Convert()
	testData, err := Parse(input)
	// ## THEN output data has:
	require.Empty(t, err, "must be no error")
	// - "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "title must be empty")
	require.Empty(t, testData.toc, "TOC must be empty")
	// - "Methods" contains 1 element: "name" = 'TestSomething', other fields are empty
	require.Equal(t, 1, len(testData.methods))
	require.Equal(t, "TestSomething", testData.methods[0].name)
	require.Empty(t, testData.methods[0].scenario, "scenario must be empty")
	require.Empty(t, testData.methods[0].steps, "steps must be empty")
	require.Empty(t, testData.methods[0].tags, "tags must be empty")
}

func TestGoFuncNameScenario(t *testing.T) {
	// > Header, Go
	// # Parse() returns data with 1 element in "Methods" on input with 1 test func and 1 non-test func.
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
		// - "func TestSomething(t *testing.T) {"
		"func TestSomething(t *testing.T) {",
		// - "// # Scenario"
		OLC + " # Scenario",
		// - "}" // the func end as not indented
		"}",
	}

	// ## WHEN Convert()
	testData, err := Parse(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	// - "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "title must be empty")
	require.Empty(t, testData.toc, "TOC must be empty")
	// - "Methods" contains 1 element: "name" = 'TestSomething', "scenario" = 'Scenario', other fields are empty
	require.Equal(t, 1, len(testData.methods))
	require.Equal(t, "TestSomething", testData.methods[0].name)
	require.Equal(t, "Scenario", testData.methods[0].scenario)
	require.Empty(t, testData.methods[0].steps, "steps must be empty")
	require.Empty(t, testData.methods[0].tags, "tags must be empty")
}

/*

func TestInputScenarioHeader(t *testing.T) {
	// > MD Markers
	// # Convert() returns header(3) line on input started with '#'
	// ## GIVEN Input is
	// - "// # Scenario"
	var input = []string{OLC + " # Scenario"}
	// ## WHEN Convert()
	output, toc, err := Parse(input)
	// ## THEN output:
	// - "### Scenario"
	require.Empty(t, err, "must be no error")
	require.Empty(t, toc, "toc must be empty")
	require.Equal(t, []string{"### Scenario"}, output)
}

func TestInputStepHeader(t *testing.T) {
	// > MD Markers
	// # Convert() returns header(4) line on input started with '##'
	// ## GIVEN Input is
	// - "// ## GIVEN"
	var input = []string{OLC + " ## GIVEN"}
	// ## WHEN Convert()
	output, toc, err := Parse(input)
	// ## THEN output:
	// - "#### GIVEN"
	require.Empty(t, err, "must be no error")
	require.Empty(t, toc, "toc must be empty")
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
	output, toc, err := Parse(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Empty(t, toc, "toc must be empty")
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
	output, toc, err := Parse(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Empty(t, toc, "toc must be empty")
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
	output, toc, err := Parse(input)
	// ## THEN output is:
	require.Empty(t, err, "must be no error")
	require.Empty(t, toc, "toc must be empty")
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
	output, toc, err := Parse(input)
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
	// - table of content is '<empty>'
	require.Empty(t, toc, "toc must be empty")
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
	output, toc, err := Parse(input)
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
	// - table of content is:
	require.Empty(t, toc, "toc must be empty")
	/*
		require.Equal(t, []string{
			// - "'1. [TestSomething1](#TestSomething1)"
			"'1. [Scenario1](#TestSomething1)",
			// - "'1. [TestSomething2](#TestSomething2)"
			"'1. [Scenario2](#TestSomething2)",
		}, toc)

}
*/
