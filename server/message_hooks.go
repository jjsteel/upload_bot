package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"strings"
)

// MessageWillBePosted comment complience
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	configuration := p.getConfiguration()
	channel, _ := p.API.GetChannel(post.ChannelId)

	if configuration.disabled || channel.Type != "O" {
		return post, ""
	}

	for _, c := range configuration.channels {
		if c == post.ChannelId {
			if post.FileIds == nil {
				return post, ""
			}
			for _, f := range post.FileIds {
				file, _ := p.API.GetFileInfo(f)
				exist, err := p.API.FileExists(file.Path)
				if err != nil {
					return nil, err.DetailedError
				}
				if exist {
					p.API.SendEphemeralPost(post.UserId, &model.Post{
						UserId:    configuration.botID,
						ChannelId: post.ChannelId,
						Message:   "File: " + file.Path,
					})
				}
				if err := p.API.RemoveFile(file.Path); err != nil {
					return nil, err.DetailedError
				}
				sli := strings.LastIndex(file.Path, "/")
				if err := p.API.RemoveDirectory(file.Path[:sli]); err != nil {
					return nil, err.DetailedError
				}
			}
			p.API.SendEphemeralPost(post.UserId, &model.Post{
				UserId:    configuration.botID,
				ChannelId: post.ChannelId,
				Message:   "Channel Reserved For Discussion!",
			})
			return nil, plugin.DismissPostError
		}
	}

	if post.FileIds == nil {
		p.API.SendEphemeralPost(post.UserId, &model.Post{
			UserId:    configuration.botID,
			ChannelId: post.ChannelId,
			Message:   "Channel Reserved For Upload!",
		})
		return nil, plugin.DismissPostError
	}

	chown, err := p.API.CopyFileInfos(configuration.botID, post.FileIds)
	if err != nil {
		p.API.SendEphemeralPost(post.UserId, &model.Post{
			UserId:    configuration.botID,
			ChannelId: post.ChannelId,
			Message:   "Something Is Wrong Try Again Later",
		})
		return nil, plugin.DismissPostError
	}

	post.UserId = configuration.botID
	post.FileIds = chown

	return post, ""
}
