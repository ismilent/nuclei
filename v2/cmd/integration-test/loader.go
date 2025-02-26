package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/ismilent/nuclei/v2/pkg/testutils"
)

var loaderTestcases = map[string]testutils.TestCase{
	"loader/template-list.yaml":             &remoteTemplateList{},
	"loader/workflow-list.yaml":             &remoteWorkflowList{},
	"loader/excluded-template.yaml":         &excludedTemplate{},
	"loader/nonexistent-template-list.yaml": &nonExistentTemplateList{},
	"loader/nonexistent-workflow-list.yaml": &nonExistentWorkflowList{},
	"loader/template-list-not-allowed.yaml": &remoteTemplateListNotAllowed{},
}

type remoteTemplateList struct{}

// Execute executes a test case and returns an error if occurred
func (h *remoteTemplateList) Execute(templateList string) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "This is test matcher text")
		if strings.EqualFold(r.Header.Get("test"), "nuclei") {
			fmt.Fprintf(w, "This is test headers matcher text")
		}
	})

	router.GET("/template_list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		file, err := os.ReadFile(templateList)
		if err != nil {
			w.WriteHeader(500)
		}
		_, err = w.Write(file)
		if err != nil {
			w.WriteHeader(500)
		}
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	configFileData := `remote-template-domain: [ "` + ts.Listener.Addr().String() + `" ]`
	err := ioutil.WriteFile("test-config.yaml", []byte(configFileData), os.ModePerm)
	if err != nil {
		return err
	}
	defer os.Remove("test-config.yaml")

	results, err := testutils.RunNucleiBareArgsAndGetResults(debug, "-target", ts.URL, "-tu", ts.URL+"/template_list", "-config", "test-config.yaml")
	if err != nil {
		return err
	}

	return expectResultsCount(results, 2)
}

type excludedTemplate struct{}

// Execute executes a test case and returns an error if occurred
func (h *excludedTemplate) Execute(templateList string) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "This is test matcher text")
		if strings.EqualFold(r.Header.Get("test"), "nuclei") {
			fmt.Fprintf(w, "This is test headers matcher text")
		}
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	results, err := testutils.RunNucleiBareArgsAndGetResults(debug, "-target", ts.URL, "-t", templateList, "-include-templates", templateList)
	if err != nil {
		return err
	}

	return expectResultsCount(results, 1)
}

type remoteTemplateListNotAllowed struct{}

// Execute executes a test case and returns an error if occurred
func (h *remoteTemplateListNotAllowed) Execute(templateList string) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "This is test matcher text")
		if strings.EqualFold(r.Header.Get("test"), "nuclei") {
			fmt.Fprintf(w, "This is test headers matcher text")
		}
	})

	router.GET("/template_list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		file, err := os.ReadFile(templateList)
		if err != nil {
			w.WriteHeader(500)
		}
		_, err = w.Write(file)
		if err != nil {
			w.WriteHeader(500)
		}
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := testutils.RunNucleiBareArgsAndGetResults(debug, "-target", ts.URL, "-tu", ts.URL+"/template_list")
	if err == nil {
		return fmt.Errorf("expected error for not allowed remote template list url")
	}

	return nil

}

type remoteWorkflowList struct{}

// Execute executes a test case and returns an error if occurred
func (h *remoteWorkflowList) Execute(workflowList string) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "This is test matcher text")
		if strings.EqualFold(r.Header.Get("test"), "nuclei") {
			fmt.Fprintf(w, "This is test headers matcher text")
		}
	})

	router.GET("/workflow_list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		file, err := os.ReadFile(workflowList)
		if err != nil {
			w.WriteHeader(500)
		}
		_, err = w.Write(file)
		if err != nil {
			w.WriteHeader(500)
		}
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	configFileData := `remote-template-domain: [ "` + ts.Listener.Addr().String() + `" ]`
	err := ioutil.WriteFile("test-config.yaml", []byte(configFileData), os.ModePerm)
	if err != nil {
		return err
	}
	defer os.Remove("test-config.yaml")

	results, err := testutils.RunNucleiBareArgsAndGetResults(debug, "-target", ts.URL, "-wu", ts.URL+"/workflow_list", "-config", "test-config.yaml")
	if err != nil {
		return err
	}

	return expectResultsCount(results, 3)
}

type nonExistentTemplateList struct{}

// Execute executes a test case and returns an error if occurred
func (h *nonExistentTemplateList) Execute(nonExistingTemplateList string) error {
	router := httprouter.New()
	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := testutils.RunNucleiBareArgsAndGetResults(debug, "-target", ts.URL, "-tu", ts.URL+"/404")
	if err == nil {
		return fmt.Errorf("expected error for nonexisting workflow url")
	}

	return nil
}

type nonExistentWorkflowList struct{}

// Execute executes a test case and returns an error if occurred
func (h *nonExistentWorkflowList) Execute(nonExistingWorkflowList string) error {
	router := httprouter.New()
	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := testutils.RunNucleiBareArgsAndGetResults(debug, "-target", ts.URL, "-wu", ts.URL+"/404")
	if err == nil {
		return fmt.Errorf("expected error for nonexisting workflow url")
	}

	return nil
}
