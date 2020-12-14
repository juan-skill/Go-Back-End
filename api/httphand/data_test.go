package httphand

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/stretchr/testify/require"
)

func TestProcessData(t *testing.T) {
	c := require.New(t)

	err := logs.InitLogger()
	c.NoError(err)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = ProcessData(ctx, "netflix.com")
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

func TestInfoServers(t *testing.T) {
	c := require.New(t)

	info, err := InfoServers("netflix.com")
	c.NoError(err)
	c.NotEmpty(info)

	info, err = InfoServers("rappi.com")
	c.NoError(err)
	c.NotEmpty(info)
	c.NotEmpty(info.Endpoints)

	info, err = InfoServers("eltiempo.com")
	c.NoError(err)
	c.NotEmpty(info)
	c.Equal(len(info.Endpoints), 2)

	_, err = InfoServers("")
	c.Error(err)
}

func TestRunWHOIS(t *testing.T) {
	c := require.New(t)

	// 52.73.161.171 server netflix
	iPAddress := "52.73.161.171"

	command := fmt.Sprintf(`whois %s | grep -i %s | cut -f 2 -d ":" | sed 's/^ *//;s/ *$//'`, iPAddress, "country")
	out, err := RunWHOIS("bash", "-c", command)
	c.NoError(err)
	c.NotEmpty(string(out))
}
