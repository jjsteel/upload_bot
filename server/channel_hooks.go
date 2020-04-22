package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// ChannelHasBeenCreated comment complience
func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	if channel.Type == "O" {
		if _, err := p.API.CreateTeamMember(channel.TeamId, configuration.botID); err != nil {
			fmt.Println(err)
		}
		if _, err := p.API.AddChannelMember(channel.Id, configuration.botID); err != nil {
			fmt.Println(err)
		}
	}

}
