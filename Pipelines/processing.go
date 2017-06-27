package Pipelines

import "net/http"

// FormToPipeline -- fills a Pipeline struct with submitted form data
// params:
// r - request reader to fetch form data or url params (unused here)
// returns:
// Pipeline struct if successful
// array of strings of errors if any occur during processing
func FormToPipeline(r *http.Request) (Pipeline, []string) {
	var pipeline Pipeline
	var errStr string
	var errs []string

	pipeline.Name, errStr = processFormField(r, "name")
	errs = appendError(errs, errStr)
	pipeline.Description, errStr = processFormField(r, "description")
	errs = appendError(errs, errStr)

	return pipeline, errs
}

func appendError(errs []string, errStr string) []string {
	if len(errStr) > 0 {
		errs = append(errs, errStr)
	}
	return errs
}

func processFormField(r *http.Request, field string) (string, string) {
	fieldData := r.PostFormValue(field)
	if len(fieldData) == 0 {
		return "", "Missing '" + field + "' parameter, cannot continue"
	}
	return fieldData, ""
}
