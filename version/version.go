package version

import "fmt"

const ServiceName = "cmdb-api"

var (
	GIT_TAG    string
	GIT_COMMIT string
	GIT_BRANCH string
	BUILD_TIME string
	GO_VERSION string
)

func FullVersion() string {
	return fmt.Sprintf("Version   : %s\nBuild Time: %s\nGit Branch: %s\nGit Commit: %s\nGo Version: %s\n", GIT_TAG, BUILD_TIME, GIT_BRANCH, GIT_COMMIT, GO_VERSION)
}

func ShortVersion() string {
	return fmt.Sprintf("%s[%s %s]", GIT_TAG, BUILD_TIME, GIT_COMMIT)
}
