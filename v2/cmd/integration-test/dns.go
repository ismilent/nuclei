package main

import (
	"github.com/ismilent/nuclei/v2/pkg/testutils"
)

var dnsTestCases = map[string]testutils.TestCase{
	"dns/basic.yaml":     &dnsBasic{},
	"dns/ptr.yaml":       &dnsPtr{},
	"dns/caa.yaml":       &dnsCAA{},
	"dns/variables.yaml": &dnsVariables{},
}

type dnsBasic struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsBasic) Execute(filePath string) error {
	var routerErr error

	results, err := testutils.RunNucleiTemplateAndGetResults(filePath, "one.one.one.one", debug)
	if err != nil {
		return err
	}
	if routerErr != nil {
		return routerErr
	}
	return expectResultsCount(results, 1)
}

type dnsPtr struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsPtr) Execute(filePath string) error {
	var routerErr error

	results, err := testutils.RunNucleiTemplateAndGetResults(filePath, "1.1.1.1", debug)
	if err != nil {
		return err
	}
	if routerErr != nil {
		return routerErr
	}
	return expectResultsCount(results, 1)
}

type dnsCAA struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsCAA) Execute(filePath string) error {
	var routerErr error

	results, err := testutils.RunNucleiTemplateAndGetResults(filePath, "google.com", debug)
	if err != nil {
		return err
	}
	if routerErr != nil {
		return routerErr
	}
	return expectResultsCount(results, 1)
}

type dnsVariables struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsVariables) Execute(filePath string) error {
	var routerErr error

	results, err := testutils.RunNucleiTemplateAndGetResults(filePath, "one.one.one.one", debug)
	if err != nil {
		return err
	}
	if routerErr != nil {
		return routerErr
	}
	return expectResultsCount(results, 1)
}
