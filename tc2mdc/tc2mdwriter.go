package tc2mdc

func Write(data *TestData) []string {
	if data == nil {
		return nil
	}
	var mdText []string

	if data.packageName != "" {
		mdText = append(mdText, "## `"+data.packageName+"`")
	}

	for _, method := range data.methods {
		appendFunc(method.name, &mdText)
		appendTags(method.tags, &mdText)
		appendScenario(method.scenario, &mdText)
		appendSteps(method.steps, &mdText)
		appendFuncEnd(data.packageName, &mdText)
	}

	return mdText
}

func appendSteps(steps []TestStep, mdText *[]string) {
	for _, step := range steps {
		*mdText = append(*mdText, getStepPrefix(step.kind)+step.comment)
	}
}

func getStepPrefix(kind int) string {
	switch kind {
	case GWT:
		{
			return "#### "
		}
	case common:
		{
			return "- "
		}
	case indented:
		{
			return "  - "
		}
	case indented2:
		{
			return "    - "
		}
	}
	return ""
}

func appendScenario(scenario string, mdText *[]string) {
	if scenario != "" {
		*mdText = append(*mdText, "### "+scenario)
	}
}

func appendFunc(name string, mdText *[]string) {
	*mdText = append(*mdText, "---")
	*mdText = append(*mdText, "#### `"+name+"`")
}

func appendFuncEnd(packageName string, mdText *[]string) {
	*mdText = append(*mdText, "")
	*mdText = append(*mdText, getLinkToTop(packageName))
}

func appendTags(tags []string, mdText *[]string) {
	if tags == nil {
		return
	}
	var sep string
	tagsLine := "> "
	for _, tag := range tags {
		tagsLine += sep + tag
		sep = ", "
	}
	*mdText = append(*mdText, tagsLine)
}

func getLinkToTop(name string) string {
	if name == "" {
		name = "top"
	}
	return "[top](#" + name + ")"
}
