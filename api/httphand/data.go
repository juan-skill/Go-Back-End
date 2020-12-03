package httphand

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
)

// InfoDomainPage contain the information about the domain
type InfoDomainPage struct {
	Title string
	Logo  string
}

var (
	// ErrEmptyDomainName when check the status server
	ErrEmptyDomainName = errors.New("cannot be empty domain name")
)

// ProcessData to build the domain object
func ProcessData(domainName string) (*models.Domain, error) {
	isDown, err := GetStatusServer(domainName)
	if err != nil {
		logs.Log().Errorf("Error isDown %s", err.Error())
		return nil, err
	}

	infoPage, err := GetInfoDomainPage(domainName)
	if err != nil {
		logs.Log().Errorf("Error infoPage %s", err.Error())
		return nil, err
	}

	domain, err := models.NewDomain(false, isDown, domainName, "", "", infoPage.Logo, infoPage.Title)
	if err != nil {
		logs.Log().Errorf("cannot create the domain %s", err.Error())
		return nil, err
	}

	return domain, nil
}

// GetStatusServer check server status
func GetStatusServer(domainName string) (bool, error) {
	if domainName == "" {
		return false, ErrEmptyDomainName
	}

	url := fmt.Sprintf("https://%s", domainName)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Log().Errorf("Error request wraps %s ", err.Error())
		return true, err
	}

	resp, err := client.Do(request)
	if err != nil {
		logs.Log().Errorf("Error check status server %s ", err.Error())
		return true, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logs.Log().Errorf("Error response body close %s ", err.Error())
		}
	}()

	statusRequest := fmt.Sprintf("%d OK", http.StatusOK)
	if resp.Status != statusRequest {
		logs.Log().Errorf("the server does not work statuscode %d\n", resp.StatusCode)
		return false, nil
	}

	return true, nil
}

// GetInfoDomainPage ...
func GetInfoDomainPage(domainName string) (*InfoDomainPage, error) {
	if domainName == "" {
		return nil, ErrEmptyDomainName
	}

	url := fmt.Sprintf("https://%s", domainName)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Log().Errorf("Error request wraps %s ", err.Error())
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		logs.Log().Errorf("Error check status server %s ", err.Error())
		return nil, err
	}

	defer func() {
		erro := resp.Body.Close()
		if erro != nil {
			logs.Log().Errorf("Error response body close %s ", erro.Error())
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Log().Errorf("Error read document HTML %s ", err.Error())
		return nil, err
	}

	title := ""

	doc.Find("title").Each(func(i int, s *goquery.Selection) { title = s.Text() })

	iconPath := ""

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		rel, err := s.Attr("rel")
		if !err {
			logs.Log().Errorf("Not found rel attribute HTML %s ", err)
			return
		}

		if rel == "shortcut icon" {
			iconPath, err = s.Attr("href")
			if !err {
				logs.Log().Errorf("Not found href attribute HTML %s ", err)
				return
			}
		}

		if rel == "icon" {
			iconPath, err = s.Attr("href")
			if !err {
				logs.Log().Errorf("Not found href attribute HTML %s ", err)
				return
			}

			re, err := regexp.Compile("^/")
			if err != nil {
				logs.Log().Errorf("Not found href attribute HTML %s ", err.Error())
				return
			}

			found := re.MatchString(iconPath)
			if found {
				iconPath = resp.Request.URL.String() + iconPath
			}
			return
		}
	})

	return &InfoDomainPage{
		Title: title,
		Logo:  iconPath,
	}, nil
}
