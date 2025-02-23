package tc2mdc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputNil(t *testing.T) {
	// > Empty Input
	// # Parse() returns error on 'nil' input
	// ## WHEN Convert(nil)
	testData, err := Parse(nil)
	// ## THEN error message: 'nil input', data is 'nil'
	require.Nil(t, testData, "data must be nil")
	require.ErrorContains(t, err, "nil input")
}

func TestInputEmpty(t *testing.T) {
	// > Empty Input
	// # Parse() returns error on 'empty' string input
	// ## GIVEN Input is ""
	var input = []string{""}
	// ## WHEN Parse("")
	testData, err := Parse(input)
	// ## THEN error message: 'empty input', data is 'nil'
	require.Nil(t, testData, "data must be nil")
	require.ErrorContains(t, err, "empty input")
}

func TestInputMultiNoComments(t *testing.T) {
	// > Empty Input
	// # Parse() returns empty data on strings without one line comments
	// ## GIVEN Input contains multi lines without one line comments
	var input = []string{" something", "", " something else"}
	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN no error, but output data is '<empty>'
	require.Nil(t, err, "must be no error")
	require.Empty(t, testData, "data must be empty")
}

func TestGoPackageName(t *testing.T) {
	// > Package, Go
	// # Parse() returns data with "packageName" on input with package name only
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
	}

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output data has:
	require.Nil(t, err, "must be no error")
	// - "packageName" = 'somePackage', "title" = '<empty>', "TOC", "Methods" are 'nil'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "must be empty title")
	require.Nil(t, testData.toc, "TOC must be nil")
	require.Nil(t, testData.methods, "methodsm must be nil")
}

func TestGoFuncNameAsMethodName(t *testing.T) {
	// > Methods, Go
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

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output data has:
	require.Nil(t, err, "must be no error")
	// - "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "title must be empty")
	require.Nil(t, testData.toc, "TOC must be empty")
	// - "Methods" contains 1 element: "name" = 'TestSomething', other fields are empty
	require.Equal(t, 1, len(testData.methods))
	require.Equal(t, "TestSomething", testData.methods[0].name)
	require.Empty(t, testData.methods[0].scenario, "scenario must be empty")
	require.Nil(t, testData.methods[0].steps, "steps must be empty")
	require.Nil(t, testData.methods[0].tags, "tags must be empty")
}

func TestGoFuncNameScenario(t *testing.T) {
	// > Methods, Go
	// # Parse() returns data with an element in "Methods" with "Name" and "Scenario"
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

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output is:
	require.Nil(t, err, "must be no error")
	// - "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "title must be empty")
	require.Empty(t, testData.toc, "TOC must be empty")
	// - "Methods" contains 1 element: "name" = 'TestSomething', "scenario" = 'Scenario', other fields are empty
	require.Equal(t, 1, len(testData.methods))
	require.Equal(t, "TestSomething", testData.methods[0].name)
	require.Equal(t, "Scenario", testData.methods[0].scenario)
	require.Nil(t, testData.methods[0].steps, "steps must be nil")
	require.Nil(t, testData.methods[0].tags, "tags must be nil")
}

func TestGoFuncNameTags(t *testing.T) {
	// > Methods, Go
	// # Parse() returns data with an element in "Methods" with multiple tags
	// ## GIVEN Input is
	var input = []string{
		// - "package somePackage"
		"package somePackage",
		// - "func TestSomething(t *testing.T) {"
		"func TestSomething(t *testing.T) {",
		// - "// > Tag1, Tag2"
		OLC + " > Tag1, Tag2",
		// - "}" // the func end as not indented
		"}",
	}

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output is:
	require.Nil(t, err, "must be no error")
	// - "packageName" = 'somePackage', ("title", "TOC") = '<empty>'
	require.Equal(t, "somePackage", testData.packageName)
	require.Empty(t, testData.title, "title must be empty")
	require.Empty(t, testData.toc, "TOC must be empty")
	// - "Methods" contains 1 element:
	require.Equal(t, 1, len(testData.methods))
	// -- "name" = 'TestSomething', "scenario" = '<empty>', other fields are empty
	require.Equal(t, "TestSomething", testData.methods[0].name)
	require.Empty(t, testData.methods[0].scenario, "scenario must be empty")
	require.Nil(t, testData.methods[0].steps, "steps must be nil")
	// -- 2 "tags" = 'Tag1', 'Tag2'
	require.Equal(t, 2, len(testData.methods[0].tags))
	require.Equal(t, "Tag1", testData.methods[0].tags[0])
	require.Equal(t, "Tag2", testData.methods[0].tags[1])
}

func TestGoFuncNameGWT(t *testing.T) {
	// > Comments, Go
	// # Parse() returns data with an element in "Methods" with 3 steps - 'GWT' comments
	// ## GIVEN Input is
	var input = []string{
		// - "func TestSomething(t *testing.T) {"
		"func TestSomething(t *testing.T) {",
		// - "// ## GIVEN set"
		OLC + " ## GIVEN set",
		// - "// ## WHEN act"
		OLC + " ## WHEN act",
		// - "// ## THEN check"
		OLC + " ## THEN check",
		// - "}" // the func end as not indented
		"}",
	}

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output is:
	require.Nil(t, err, "must be no error")
	// - "Methods" contains 1 element:
	require.Equal(t, 1, len(testData.methods))
	// -- 3 "steps" :
	require.Equal(t, 3, len(testData.methods[0].steps))
	// --- {0, 'GIVEN set'}
	require.Equal(t, GWT, testData.methods[0].steps[0].kind)
	require.Equal(t, "GIVEN set", testData.methods[0].steps[0].comment)
	// --- {0, 'WHEN act'}
	require.Equal(t, GWT, testData.methods[0].steps[1].kind)
	require.Equal(t, "WHEN act", testData.methods[0].steps[1].comment)
	// --- {0, 'THEN check'}
	require.Equal(t, GWT, testData.methods[0].steps[2].kind)
	require.Equal(t, "THEN check", testData.methods[0].steps[2].comment)
}

func TestGoFuncNameSteps123(t *testing.T) {
	// > Comments, Go
	// # Parse() returns data with an element in "Methods" with 3 steps - comments, common and indented
	// ## GIVEN Input is
	var input = []string{
		// - "func TestSomething(t *testing.T) {"
		"func TestSomething(t *testing.T) {",
		// - "// - common comment"
		OLC + " - common comment",
		// - "// -- indented comment"
		OLC + " -- indented comment",
		// - "// --- indented twice comment"
		OLC + " --- indented twice comment",
		// - "}" // the func end as not indented
		"}",
	}

	// ## WHEN Parse()
	testData, err := Parse(input)
	// ## THEN output is:
	require.Nil(t, err, "must be no error")
	// - "Methods" contains 1 element:
	require.Equal(t, 1, len(testData.methods))
	// -- 3 "steps" :
	require.Equal(t, 3, len(testData.methods[0].steps))
	// --- {1, 'common comment'}
	require.Equal(t, common, testData.methods[0].steps[0].kind)
	require.Equal(t, "common comment", testData.methods[0].steps[0].comment)
	// --- {2, 'indented comment'}
	require.Equal(t, indented, testData.methods[0].steps[1].kind)
	require.Equal(t, "indented comment", testData.methods[0].steps[1].comment)
	// --- {3, 'indented2 comment'}
	require.Equal(t, indented2, testData.methods[0].steps[2].kind)
	require.Equal(t, "indented twice comment", testData.methods[0].steps[2].comment)
}
