// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "upload.bot",
  "name": "Upload Bot",
  "description": "upload bot",
  "version": "0.1.0",
  "min_server_version": "5.21.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64"
    },
    "executable": ""
  },
  "settings_schema": {
    "header": "",
    "footer": "",
    "settings": [
      {
        "key": "BotName",
        "display_name": "Bot Username",
        "type": "text",
        "help_text": "Username that is going be used for the bot",
        "placeholder": "",
        "default": "upload_bot"
      },
      {
        "key": "PublicChannels",
        "display_name": "Channels List",
        "type": "text",
        "help_text": "Public channels that should be monitored separated by ';' ",
        "placeholder": "",
        "default": "test;test1;test2"
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
