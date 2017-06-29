package Stages

import (
	"net/http"
	"strconv"
)

// FormToStage -- fills a Stage struct with submitted form data
// params:
// r - request reader to fetch form data or url params (unused here)
// returns:
// Stage struct if successful
// array of strings of errors if any occur during processing
func FormToStage(r *http.Request) (Stage, []string) {
	var stage Stage
	var errStr, versionStr string
	var errs []string
	var err error

	stage.Name, errStr = processFormField(r, "name")
	errs = appendError(errs, errStr)
	stage.Description, errStr = processFormField(r, "description")
	errs = appendError(errs, errStr)
	stage.Type, errStr = processFormField(r, "type")
	errs = appendError(errs, errStr)
	stage.Payload, errStr = processFormField(r, "payload")
	errs = appendError(errs, errStr)

	versionStr, errStr = processFormField(r, "version")
	if len(errStr) != 0 {
		errs = append(errs, errStr)
	} else {
		stage.Version, err = strconv.Atoi(versionStr)
		if err != nil {
			errs = append(errs, "Paramenter 'version' not an integer")
		}
	}

	return stage, errs
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
