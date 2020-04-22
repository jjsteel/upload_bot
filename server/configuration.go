package main

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type configuration struct {
	BotName        string
	PublicChannels string
	disabled       bool
	botID          string
	channels       []string
}

func (c *configuration) Clone() *configuration {
	return &configuration{
		BotName:        c.BotName,
		PublicChannels: c.PublicChannels,
		disabled:       c.disabled,
		botID:          c.botID,
		channels:       c.channels,
	}
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

// OnConfigurationChange comment complience
func (p *Plugin) OnConfigurationChange() error {
	configuration := p.getConfiguration().Clone()
	var err error

	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}
	// EnsureBot runs only once if no bot is present
	bot, ensureBotError := p.Helpers.EnsureBot(&model.Bot{
		Username:    configuration.BotName,
		DisplayName: "Upload Bot",
		Description: " monitor file upload in public channels",
	}, plugin.IconImagePath("/assets/robot_tv.svg"))
	if ensureBotError != nil {
		return errors.Wrap(ensureBotError, "failed to ensure upload bot.")
	}

	configuration.botID = bot
	configuration.channels, err = p.getPublicChannels(configuration)

	if err != nil {
		return errors.Wrap(err, "failed to get public channels")
	}

	p.setConfiguration(configuration)
	return nil
}

func (p *Plugin) getPublicChannels(configuration *configuration) ([]string, error) {
	teams, err := p.API.GetTeams()
	if err != nil {
		return nil, err
	}

	channels := strings.Split(configuration.PublicChannels, ";")
	publicChannelIds := make([]string, len(teams)*len(channels))
	i := 0
	for _, team := range teams {
		for _, channel := range channels {
			ch, err := p.API.GetChannelByName(team.Id, channel, false)
			if err != nil {
				continue
			}
			if _, err := p.API.CreateTeamMember(team.Id, configuration.botID); err != nil {
				return nil, errors.Wrap(err, "Team Member Cration")
			}
			publicChannelIds[i] = ch.Id
			i++
		}
	}
	return publicChannelIds[:i], nil
}

// setEnabled wraps setConfiguration to configure if the plugin is enabled.
func (p *Plugin) setEnabled(enabled bool) {
	var configuration = p.getConfiguration().Clone()
	configuration.disabled = !enabled

	p.setConfiguration(configuration)
}
