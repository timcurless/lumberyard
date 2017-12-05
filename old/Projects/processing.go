package Projects

import (
  "net/http"
)

func FormToProject(r *http.Request) (Project, []string) {
	var project Project
	var errStr string
	var errs []string

	project.Name, errStr = processFormField(r, "name")
	errs = appendError(errs, errStr)
	project.Email, errStr = processFormField(r, "email")
	errs = appendError(errs, errStr)

	return project, errs
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