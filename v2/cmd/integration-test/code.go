package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"

	"github.com/ismilent/nuclei/v2/pkg/catalog"
	"github.com/ismilent/nuclei/v2/pkg/catalog/config"
	"github.com/ismilent/nuclei/v2/pkg/catalog/loader"
	"github.com/ismilent/nuclei/v2/pkg/core"
	"github.com/ismilent/nuclei/v2/pkg/core/inputs"
	"github.com/ismilent/nuclei/v2/pkg/output"
	"github.com/ismilent/nuclei/v2/pkg/parsers"
	"github.com/ismilent/nuclei/v2/pkg/protocols"
	"github.com/ismilent/nuclei/v2/pkg/protocols/common/hosterrorscache"
	"github.com/ismilent/nuclei/v2/pkg/protocols/common/interactsh"
	"github.com/ismilent/nuclei/v2/pkg/protocols/common/protocolinit"
	"github.com/ismilent/nuclei/v2/pkg/protocols/common/protocolstate"
	"github.com/ismilent/nuclei/v2/pkg/reporting"
	"github.com/ismilent/nuclei/v2/pkg/testutils"
	"github.com/ismilent/nuclei/v2/pkg/types"
	"github.com/julienschmidt/httprouter"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/goflags"
	"go.uber.org/ratelimit"
)

var codeTestcases = map[string]testutils.TestCase{
	"code/test.yaml": &goIntegrationTest{},
}

type goIntegrationTest struct{}

// Execute executes a test case and returns an error if occurred
//
// Execute the docs at ../DESIGN.md if the code stops working for integration.
func (h *goIntegrationTest) Execute(templatePath string) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "This is test matcher text")
		if strings.EqualFold(r.Header.Get("test"), "nuclei") {
			fmt.Fprintf(w, "This is test headers matcher text")
		}
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	results, err := executeNucleiAsCode(templatePath, ts.URL)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}

// executeNucleiAsCode contains an example
func executeNucleiAsCode(templatePath, templateURL string) ([]string, error) {
	cache := hosterrorscache.New(30, hosterrorscache.DefaultMaxHostsCount)
	defer cache.Close()

	mockProgress := &testutils.MockProgressClient{}
	reportingClient, _ := reporting.New(&reporting.Options{}, "")
	defer reportingClient.Close()

	outputWriter := testutils.NewMockOutputWriter()
	var results []string
	outputWriter.WriteCallback = func(event *output.ResultEvent) {
		results = append(results, fmt.Sprintf("%v\n", event))
	}

	defaultOpts := types.DefaultOptions()
	_ = protocolstate.Init(defaultOpts)
	_ = protocolinit.Init(defaultOpts)

	defaultOpts.Templates = goflags.FileOriginalNormalizedStringSlice{templatePath}
	defaultOpts.ExcludeTags = config.ReadIgnoreFile().Tags

	interactOpts := interactsh.NewDefaultOptions(outputWriter, reportingClient, mockProgress)
	interactClient, err := interactsh.New(interactOpts)
	if err != nil {
		return nil, errors.Wrap(err, "could not create interact client")
	}
	defer interactClient.Close()

	home, _ := os.UserHomeDir()
	catalog := catalog.New(path.Join(home, "nuclei-templates"))
	executerOpts := protocols.ExecuterOptions{
		Output:          outputWriter,
		Options:         defaultOpts,
		Progress:        mockProgress,
		Catalog:         catalog,
		IssuesClient:    reportingClient,
		RateLimiter:     ratelimit.New(150),
		Interactsh:      interactClient,
		HostErrorsCache: cache,
		Colorizer:       aurora.NewAurora(true),
		ResumeCfg:       types.NewResumeCfg(),
	}
	engine := core.New(defaultOpts)
	engine.SetExecuterOptions(executerOpts)

	workflowLoader, err := parsers.NewLoader(&executerOpts)
	if err != nil {
		log.Fatalf("Could not create workflow loader: %s\n", err)
	}
	executerOpts.WorkflowLoader = workflowLoader

	configObject, err := config.ReadConfiguration()
	if err != nil {
		return nil, errors.Wrap(err, "could not read configuration file")
	}
	store, err := loader.New(loader.NewConfig(defaultOpts, configObject, catalog, executerOpts))
	if err != nil {
		return nil, errors.Wrap(err, "could not create loader")
	}
	store.Load()

	input := &inputs.SimpleInputProvider{Inputs: []string{templateURL}}
	_ = engine.Execute(store.Templates(), input)
	engine.WorkPool().Wait() // Wait for the scan to finish

	return results, nil
}
