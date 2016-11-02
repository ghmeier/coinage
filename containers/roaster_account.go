package containers

import(
	"database/sql"

	"github.com/pborman/uuid"
)

type RoasterAccount struct {
	Id uuid.UUID `json:"id"`
}

func FromSql(rows *sql.Rows) ([]*RoasterAccount, error) {
	roasterAccount := make([]*RoasterAccount,0)

	for rows.Next() {
		c := &RoasterAccount{}
		rows.Scan(&c.Id)
		roasterAccount = append(roasterAccount, c)
	}

	return roasterAccount, nil
}
