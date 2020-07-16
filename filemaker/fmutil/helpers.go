package fmutil

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/api"
	fm "github.com/amanbolat/gofmcon"
	"github.com/pkg/errors"
	"strings"
)

const (
	SCRIPT_AUDIT_LOG    = "api_audit_log"
	SCRIPT_DELIMITER    = "|"
	COMMON_QUERY_LAYOUT = "common_query_layout"
)

func validatePagination(q *fm.FMQuery, rm api.RequestMeta) {
	q = q.Skip(rm.Skip)
	if rm.PerPage >= 1 {
		q = q.Max(rm.PerPage)
	}
}

// withAudit adds script and params for audit log
func WithAudit(q *fm.FMQuery, id, table, field, data, user string) {
	q.WithPostFindScript(SCRIPT_AUDIT_LOG, strings.Join([]string{id, table, field, data, user}, SCRIPT_DELIMITER))
}

type Store interface {
	DBName() string
	FMConn() *fm.FMConnector
}

func GetFileMakerRecordList(store Store, q *fm.FMQuery, reqMeta api.RequestMeta) ([]*fm.Record, api.ResponseMeta, error) {
	var resMeta api.ResponseMeta
	if q == nil {
		return nil, resMeta, errors.New("filemaker query is empty")
	}

	for _, sf := range reqMeta.SortFields {
		if strings.TrimSpace(sf.Name) == "" {
			continue
		}
		order := fm.Ascending
		if sf.Descending {
			order = fm.Descending
		}

		q.WithSortFields(fm.FMSortField{
			Name:  sf.Name,
			Order: order,
		})
	}

	reqMeta.Check()
	validatePagination(q, reqMeta)

	resMeta = api.ResponseMeta{
		Page:  reqMeta.Page,
		Count: 0,
		Total: 0,
	}

	fmset, err := store.FMConn().Query(q)
	if err != nil {
		// if error is "No records match the request"
		// don't return error
		if err.Error() == fmt.Sprintf("filemaker_error: %s", fm.FileMakerErrorCodes[401]) {
			return []*fm.Record{}, resMeta, nil
		}

		return nil, resMeta, err
	}

	resMeta.Count = len(fmset.Resultset.Records)
	resMeta.Total = fmset.Resultset.Count

	if len(fmset.Resultset.Records) < 1 {
		return []*fm.Record{}, resMeta, nil
	}

	return fmset.Resultset.Records, resMeta, nil
}

func GetFileMakerRecordSingle(store Store, q *fm.FMQuery) (*fm.Record, error) {
	if q == nil {
		return nil, errors.New("filemaker query is empty")
	}

	q.Max(1)

	fmSet, err := store.FMConn().Query(q)
	if err != nil {
		return nil, errors.WithMessage(err, "database_error")
	}

	if len(fmSet.Resultset.Records) == 0 {
		return nil, errors.New("record_not_found")
	}

	return fmSet.Resultset.Records[0], nil
}

func ConvertToBool(i int) bool {
	return i == 1
}

func ConvertBoolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
