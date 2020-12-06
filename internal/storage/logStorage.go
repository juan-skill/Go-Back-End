package storage

import (
	"context"
	"fmt"

	"github.com/other_project/crockroach/models"
)

var (
	// Objects that contain the last record
	Objects map[string]*models.LogDomainStatus
)

// GetRecordByName return a list of records saved an hour or less ago
func (q *Queries) GetRecordByName(domain *models.Domain) (records []*models.Domain, err error) {
	if domain == nil {
		return nil, models.ErrEmptyDomain
	}

	// sobre esta lista procesaremos los dominios visitados
	recordByIds := []*models.Domain{}

	for _, value := range Objects {
		if domain.DomainName == value.Domain.DomainName {
			recordByIds = append(recordByIds, value.Domain)
		}
	}

	fmt.Println()
	fmt.Println(Objects)
	fmt.Println()

	for i := 0; i < len(recordByIds); i++ {
		fmt.Println("getrecordbyname --> ", recordByIds[i])
	}

	return recordByIds, nil
}

// GetLastDomain list the last domains consulted
func (q *Queries) GetLastDomain() map[string]*models.Domain {
	fmt.Println()
	fmt.Println(Objects)
	fmt.Println()

	myObjects := make(map[string]*models.Domain)

	for _, value := range Objects {
		myObjects[value.DomainName] = value.Domain
		fmt.Println("value ---> ", value.Domain)
	}

	fmt.Println()
	fmt.Println("objects --> ", myObjects)

	fmt.Println()

	return myObjects
}

// NewRecord creates a new record about of last record/changes
func (q *Queries) NewRecord(domain *models.Domain) (*models.LogDomainStatus, error) {
	if domain == nil {
		return nil, models.ErrEmptyDomain
	}

	logDomain, err := models.NewLogDomainStatus(domain.DomainName, domain.SSLGrade, domain)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	Objects[fmt.Sprintf("%d.%s", logDomain.LogDomainStatusID, logDomain.DomainName)] = logDomain

	return logDomain, nil
}

// ReloadRecord the records saved an hour or less ago
func (q *Queries) ReloadRecord(ctx context.Context) (myObjects map[string]*models.LogDomainStatus, err error) {
	domains, err := GetDomains(ctx, "1 hours")
	if err != nil {
		return nil, err
	}

	fmt.Println("reload ", len(domains))

	for domain := range domains {
		_, err := NewRecord(&domains[domain])
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(Objects)

	return Objects, nil
}

func init() {
	Objects = make(map[string]*models.LogDomainStatus)
}
