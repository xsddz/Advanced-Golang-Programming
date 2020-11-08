package sql

// CURDer -
type CURDer interface {
	Insert(keyValues map[string]interface{}) (res Resulter, err error)
	Delete(conds [][]interface{}) (res Resulter, err error)
	Update(keyValues map[string]interface{}, conds [][]interface{}) (res Resulter, err error)
	Select(fields []string, conds [][]interface{}, appends []string) (res Resulter, err error)
}
