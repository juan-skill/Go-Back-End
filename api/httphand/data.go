package httphand

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
)

// InfoLabSSLEndpoints contain the information of a domain
type InfoLabSSLEndpoints struct {
	IPAddress         string
	ServerName        string
	StatusMessage     string
	Grade             string
	GradeTrustIgnored string
	HasWarmings       string
	IsExceptional     bool
	Progress          int64
	Duration          int64
	Delegation        int64
}

// InfoLabSSL contain the info SSL LAB API
type InfoLabSSL struct {
	Host            string
	Port            int64
	Protocol        string
	IsPublic        bool
	Status          string
	StartTime       int64
	TestTime        int64
	EngineVersion   string
	CriteriaVersion string
	Endpoints       []*InfoLabSSLEndpoints
}

// ParseServerJSON model structure for parse server
type ParseServerJSON struct {
	Address  string `json:"address"`
	SSLGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

//InfoWHOISCommand model struct owner and country
type InfoWHOISCommand struct {
	country string
	owner   string
}

// ParseDomainJSON model structure for parse domain
type ParseDomainJSON struct {
	Servers          []*ParseServerJSON `json:"servers"`
	ServerChanged    bool               `json:"servers_changed"`
	SSLGrade         string             `json:"ssl_grade"`
	PreviousSSLGrade string             `json:"previous_ssl_grade"`
	Logo             string             `json:"logo"`
	Title            string             `json:"title"`
	IsDown           bool               `json:"is_down"`
}

// InfoDomainPage contain the information about the domain
type InfoDomainPage struct {
	Title string
	Logo  string
}

const (
	// Timeout time to perform the request to the API
	Timeout = 15 * time.Second
)

var (
	// ErrEmptyDomainName when check the status server
	ErrEmptyDomainName = errors.New("cannot be empty domain name")
	// ErrInvalidServers when search info servers
	ErrInvalidServers = errors.New("cannot extract data about the servers")
	// ErrWithoutAnwserSSLLabs when search
	ErrWithoutAnwserSSLLabs = errors.New("cannot obtain  answser SSL labs info")
	// ErrDomainConsulted when search the domain
	ErrDomainConsulted = errors.New("cannot obtain answer the domain")
)

// ProcessData to build the domain object
func ProcessData(ctx context.Context, domainName string) (*models.Domain, error) {
	isDown, err := GetStatusServer(domainName)
	if err != nil {
		//logs.Log().Errorf("Error isDown %s", err.Error())
		return nil, err
	}

	infoPage, err := GetInfoDomainPage(domainName)
	if err != nil {
		//logs.Log().Errorf("Error infoPage %s", err.Error())
		return nil, err
	}

	domain, err := models.NewDomain(false, isDown, domainName, "", "", infoPage.Logo, infoPage.Title)
	if err != nil {
		//logs.Log().Errorf("cannot create the domain %s", err.Error())
		return nil, err
	}

	infoDomainSSL, err := InfoServers(domainName)
	if err != nil {
		return nil, err
	}

	servers := infoDomainSSL.Endpoints

	serversNumber := len(servers)

	for i := 0; i < serversNumber; i++ {
		serverSSL := servers[i]

		infoWhois, err := getInfoWhois(serverSSL.IPAddress)
		if err != nil {
			//logs.Log().Errorf("cannot extract Country whois command: %s", err.Error())
			return nil, err
		}

		server, err := models.NewServer(serverSSL.IPAddress, serverSSL.Grade, infoWhois.country, infoWhois.owner, domain)
		if err != nil {
			//logs.Log().Errorf("cannot create the server of the domain %s", err.Error())
			return nil, err
		}

		domain.Servers = append(domain.Servers, server)
	}

	return domain, nil
}

// GetStatusServer check server status
func GetStatusServer(domainName string) (bool, error) {
	if domainName == "" {
		return false, ErrEmptyDomainName
	}

	url := fmt.Sprintf("https://%s", domainName)
	timeout := time.Duration(Timeout)
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
		return false, err
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
		logs.Log().Errorf("missing domain name %s ", ErrEmptyDomainName)
		return nil, ErrEmptyDomainName
	}

	url := fmt.Sprintf("https://%s", domainName)
	timeout := time.Duration(Timeout)
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

	statusRequest := fmt.Sprintf("%d OK", http.StatusOK)
	if resp.Status != statusRequest {
		logs.Log().Errorf("the dominio %s does not work: statuscode %d\n", resp.Request.URL, resp.StatusCode)
		return nil, ErrDomainConsulted
	}

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
			logs.Log().Errorf("Not found rel attribute HTML %v ", err)
			return
		}

		if rel == "shortcut icon" {
			iconPath, err = s.Attr("href")
			if !err {
				logs.Log().Errorf("Not found href attribute HTML %v ", err)
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

// InfoServers ...
func InfoServers(domain string) (*InfoLabSSL, error) {
	if domain == "" {
		return nil, ErrEmptyDomainName
	}

	url := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", domain)

	timeout := time.Duration(Timeout)

	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Log().Errorf("Error wraps request %s", err.Error())
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		logs.Log().Errorf("Error wraps request %s", err.Error())
		return nil, err
	}

	defer func() {
		erro := resp.Body.Close()
		if erro != nil {
			logs.Log().Errorf("Error response body close %s ", erro.Error())
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Log().Errorf("Error response body close %s ", err.Error())
		return nil, err
	}

	//var result map[string]map[string]string
	var resultSSL map[string]interface{}

	err = json.Unmarshal(body, &resultSSL)
	if err != nil {
		logs.Log().Errorf("Error unmarshal infoDomainSSL %s ", err.Error())
		return nil, err
	}

	if resultSSL["errors"] != nil {
		logs.Log().Errorf("Error unmarshal infoDomainSSL %s ", resultSSL["errors"])
		return nil, ErrWithoutAnwserSSLLabs
	}

	var infoDomainSSL InfoLabSSL

	err = json.Unmarshal(body, &infoDomainSSL)
	if err != nil {
		logs.Log().Errorf("Error unmarshal infoDomainSSL %s ", err.Error())
		return nil, err
	}

	fmt.Println("struct 1 info ssl-labs: ", infoDomainSSL)

	if infoDomainSSL.Endpoints == nil {
		logs.Log().Errorf("cannot found info servers %s", ErrInvalidServers.Error())
		return nil, ErrInvalidServers
	}

	return &infoDomainSSL, nil
}

// RunWHOIS get info about domin
func RunWHOIS(cmd string, args ...string) (string, error) {
	value, err := exec.Command("bash", args...).Output() //nolint:gosec
	if err != nil {
		fmt.Printf("Failed to execute command: %s", err.Error())
		return "", err
	}

	return string(value), nil
}

// getInfoWhois ....
func getInfoWhois(ipAddress string) (*InfoWHOISCommand, error) {
	infoWhois := new(InfoWHOISCommand)

	command := fmt.Sprintf(`whois %s | grep -i %s | cut -f 2 -d ":" | sed 's/^ *//;s/ *$//'`, ipAddress, "country")

	value, err := RunWHOIS("bash", "-c", command)
	if err != nil {
		fmt.Printf("cannot extract country whois command %s", err.Error())
		return nil, err
	}

	infoWhois.country = string(value)

	command = fmt.Sprintf(`whois %s | grep -i %s | cut -f 2 -d ":" | sed 's/^ *//;s/ *$//'`, ipAddress, "name")

	value, err = RunWHOIS("bash", "-c", command)
	if err != nil {
		fmt.Printf("cannot extract name whois command %s", err.Error())
		return nil, err
	}

	infoWhois.owner = string(value)

	return infoWhois, nil
}

// parseJSON parse the data to return to the API
func parseJSON(domain *models.Domain) *ParseDomainJSON {
	parseDomain := new(ParseDomainJSON)

	serversNumber := len(domain.Servers)

	if len(domain.Servers) == 0 {
		return nil
	}

	for i := 0; i < serversNumber; i++ {
		parseServer := new(ParseServerJSON)
		parseServer.Address = domain.Servers[i].Address
		parseServer.SSLGrade = domain.Servers[i].SSLGrade
		parseServer.Country = domain.Servers[i].Country
		parseServer.Owner = domain.Servers[i].Owner
		parseDomain.Servers = append(parseDomain.Servers, parseServer)
	}

	parseDomain.ServerChanged = domain.ServerChanged
	parseDomain.SSLGrade = domain.SSLGrade
	parseDomain.PreviousSSLGrade = domain.PreviousSSLGrade
	parseDomain.Logo = domain.Logo
	parseDomain.Title = domain.Title
	parseDomain.IsDown = domain.IsDown

	return parseDomain
}

// parseListJSON parse the data to return to the API
func parseListJSON(domains []*models.Domain) []*ParseDomainJSON {
	objects := []*ParseDomainJSON{}

	for i := 0; i < len(domains); i++ {
		parseDomain := parseJSON(domains[i])
		objects = append(objects, parseDomain)
	}

	return objects
}
