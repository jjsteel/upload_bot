module github.com/mattermost/upload_bot

go 1.12

require (
	github.com/mattermost/mattermost-server/v5 v5.0.0
	github.com/mholt/archiver/v3 v3.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
)
replace github.com/mattermost/mattermost-server/v5 => ../mattermost-server
