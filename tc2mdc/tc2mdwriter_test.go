package tc2mdc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteNil(t *testing.T) {
	// > Write to MD
	// # Write() returns 'nil' data as 'nil'
	// ## WHEN Write(nil)
	mdText := Write(nil)
	// ## THEN - MD text is 'nil'
	require.Nil(t, mdText, "MD text must be empty")
}

func TestWritePackageName(t *testing.T) {
	// > Write to MD
	// # Write() returns "packageName" as header(2)
	// ## GIVEN - testData: "packageName" = 'SomePackage'
	var testData = new(TestData)
	testData.packageName = "SomePackage"

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes 1 line:
	require.Equal(t, []string{
		// - "## `SomePackage`"
		"## `" + testData.packageName + "`",
	}, mdText)
}

func TestWriteOneMethod(t *testing.T) {
	// > Write to MD
	// # Write() returns test method name as header(4) with separator before name and a link to the top after
	// ## GIVEN - testData: "packageName" = ''
	var testData = new(TestData)
	// - 1 element in "methods": "name" = 'TestSomething'
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething"})

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes 4 lines:
	require.Equal(t, []string{
		// - "---"
		"---",
		// - "#### `TestSomething`"
		"#### `TestSomething`",
		// - "" // to separate the link and the last step
		"",
		// - "[top]#top" - link to the top
		"[top](#top)",
	}, mdText)
}

func TestWritePackageTwoMethods(t *testing.T) {
	// > Write to MD
	// # Write() returns test method name as header(4) with separator before name and a link to the top after
	// ## GIVEN - testData: "packageName" = 'somePackage'
	var testData = new(TestData)
	testData.packageName = "somePackage"
	// - 2 elements in "methods": "name" = 'TestSomething1' and 'TestSomething2'
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething1"})
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething2"})

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes 11 lines:
	require.Equal(t, []string{
		// - "## `somePackage`"
		"## `somePackage`",
		// - "---"
		"---",
		// - "#### `TestSomething1`"
		"#### `TestSomething1`",
		// - "", "[top]#somePackage" - link to the package line
		"",
		"[top](#somePackage)",
		// - "---"
		"---",
		// - "#### `TestSomething2`"
		"#### `TestSomething2`",
		// - "", "[top]#somePackage" - link to the package line
		"",
		"[top](#somePackage)",
	}, mdText)
}

func TestWriteMethodTags(t *testing.T) {
	// > Write to MD
	// # Write() returns "tags" as one line - comma separated list prefixed with '>'
	// ## GIVEN - testData: "packageName" = ''
	var testData = new(TestData)
	// - 1 element in "methods": "name" = 'TestSomething'
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething"})
	// - two "tags": "Tag", "Complex Tag"
	testData.methods[0].tags = []string{"Tag", "Complex Tag"}

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes one tags line after method name line:
	require.Equal(t, []string{
		"---",
		// - "#### `TestSomething`"
		"#### `TestSomething`",
		// - "> Tag, Complex Tag"
		"> Tag, Complex Tag",
		"",
		"[top](#top)",
	}, mdText)
}

func TestWriteMethodScenario(t *testing.T) {
	// > Write to MD
	// # Write() returns "scenario" as header(3) line
	// ## GIVEN - testData: "packageName" = ''
	var testData = new(TestData)
	// - 1 element in "methods": "name" = 'TestSomething'
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething"})
	// - "scenario" = 'Something happens'
	testData.methods[0].scenario = "Something happens"

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes a scenario line after method name line:
	require.Equal(t, []string{
		"---",
		// - "#### `TestSomething`"
		"#### `TestSomething`",
		// - "### Something happens"
		"### Something happens",
		"",
		"[top](#top)",
	}, mdText)
}

func TestWriteMethodGWT(t *testing.T) {
	// > Write to MD
	// # Write() returns 'GWT' comments as header(4) lines
	// ## GIVEN - testData: "packageName" = ''
	var testData = new(TestData)
	// - 1 element in "methods": "name" = 'TestSomething'
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething"})
	// - 3 steps of 'GWT' kind: "GIVEN set", "WHEN act", "THEN check"
	testData.methods[0].steps = append(testData.methods[0].steps, TestStep{kind: GWT, comment: "GIVEN set"})
	testData.methods[0].steps = append(testData.methods[0].steps, TestStep{kind: GWT, comment: "WHEN act"})
	testData.methods[0].steps = append(testData.methods[0].steps, TestStep{kind: GWT, comment: "THEN check"})

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes 3 'GWT' step lines as headers(4):
	require.Equal(t, []string{
		"---",
		// - "#### `TestSomething`"
		"#### `TestSomething`",
		// - "#### GIVEN set"
		"#### GIVEN set",
		// - "#### WHEN act"
		"#### WHEN act",
		// - "#### THEN check"
		"#### THEN check",
		"",
		"[top](#top)",
	}, mdText)
}

func TestWriteMethodIndentedSteps(t *testing.T) {
	// > Write to MD
	// # Write() returns indented comments as indented lines prefixed with '-'
	// ## GIVEN - testData: "packageName" = ''
	var testData = new(TestData)
	// - 1 element in "methods": "name" = 'TestSomething'
	testData.methods = append(testData.methods, TestMethod{name: "TestSomething"})
	// - 3 steps: common 'Step1', indented 'Step2', indented twice 'Step3'
	testData.methods[0].steps = append(testData.methods[0].steps, TestStep{kind: common, comment: "Step1"})
	testData.methods[0].steps = append(testData.methods[0].steps, TestStep{kind: indented, comment: "Step2"})
	testData.methods[0].steps = append(testData.methods[0].steps, TestStep{kind: indented2, comment: "Step3"})

	// ## WHEN Write()
	mdText := Write(testData)

	// ## THEN - MD text includes 3 indented step lines of 3 levels:
	require.Equal(t, []string{
		"---",
		// - "#### `TestSomething`"
		"#### `TestSomething`",
		// - "- Step1"
		"- Step1",
		// - "__- Step2"
		"  - Step2",
		// - "____- Step3"
		"    - Step3",
		"",
		"[top](#top)",
	}, mdText)
}
