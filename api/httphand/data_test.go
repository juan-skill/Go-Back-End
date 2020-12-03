package httphand

import (
	"testing"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/stretchr/testify/require"
)

func TestProcessData(t *testing.T) {
	c := require.New(t)

	err := logs.InitLogger()
	c.NoError(err)

	_, err = ProcessData("netflix.com")
	c.NoError(err)
}

func TestGetStatusServer(t *testing.T) {
	c := require.New(t)

	err := logs.InitLogger()
	c.NoError(err)

	isdown, err := GetStatusServer("google.com")
	c.NoError(err)
	c.Equal(true, isdown)

	isdown, err = GetStatusServer("")
	c.EqualError(ErrEmptyDomainName, err.Error())
	c.Equal(false, isdown)

	isdown, err = GetStatusServer("googlee.com")
	c.Error(err)
	c.Equal(true, isdown)
}

func TestGetInfoDomainPage(t *testing.T) {
	c := require.New(t)

	err := logs.InitLogger()
	c.NoError(err)

	infoPage, err := GetInfoDomainPage("gitlab.com")
	c.NoError(err)
	c.NotEmpty(infoPage)
	c.Equal("https://about.gitlab.com//ico/favicon-32x32.png", infoPage.Logo)
	c.Equal("\nDevOps Platform Delivered as a Single Application\n|\nGitLab\n", infoPage.Title)

	infoPage, err = GetInfoDomainPage("")
	c.Error(err)
	c.Empty(infoPage)
	c.EqualError(ErrEmptyDomainName, err.Error())

	_, err = GetInfoDomainPage("gitlabb.com")
	c.Error(err)
}
