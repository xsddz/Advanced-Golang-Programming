package sql

// ORM -
type ORM interface {
	Insert(keyValues map[string]interface{}) (Resulter, error)
	BatchInsert(keyValues []map[string]interface{}) ([]Resulter, error)
	Delete(conds [][]interface{}) (Resulter, error)
	Update(keyValues map[string]interface{}, conds [][]interface{}) (Resulter, error)
	Select(fields []string, conds [][]interface{}, appends []string) (Resulter, error)
}
