package networkclientpool

import (
	"github.com/ismilent/nuclei/v2/pkg/protocols/common/protocolstate"
	"github.com/ismilent/nuclei/v2/pkg/types"
	"github.com/projectdiscovery/fastdialer/fastdialer"
)

var (
	normalClient *fastdialer.Dialer
)

// Init initializes the clientpool implementation
func Init(options *types.Options /*TODO review unused parameter*/) error {
	// Don't create clients if already created in the past.
	if normalClient != nil {
		return nil
	}
	normalClient = protocolstate.Dialer
	return nil
}

// Configuration contains the custom configuration options for a client
type Configuration struct{}

// Hash returns the hash of the configuration to allow client pooling
func (c *Configuration) Hash() string {
	return ""
}

// Get creates or gets a client for the protocol based on custom configuration
func Get(options *types.Options, configuration *Configuration /*TODO review unused parameters*/) (*fastdialer.Dialer, error) {
	return normalClient, nil
}
