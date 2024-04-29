package errs

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moonlit0114/rest/response"
)

var PgxErrorMap = map[string]*response.RestResponse{
	pgerrcode.UniqueViolation: &response.ErrExistentResponse,
}

func getPgErrResponse(err error) *response.RestResponse {
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if response, ok := PgxErrorMap[pge.Code]; ok {
			return response
		}
	}
	return nil
}
