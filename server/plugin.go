package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin comment complience
type Plugin struct {
	plugin.MattermostPlugin
	configurationLock sync.RWMutex
	configuration     *configuration
}
