package sql

// Resulter -
type Resulter interface {
	Unmarshal(v interface{}) error

	SaveSQL(sql string) error
	LastSQL() string
}

// DefaultResult -
type DefaultResult struct {
	sqls []string
	rows []map[string]interface{}
}

// Unmarshal - 将结果解码到变量v中
func (r *DefaultResult) Unmarshal(v interface{}) error {

	return nil
}

// SaveSQL -
func (r *DefaultResult) SaveSQL(sql string) error {
	r.sqls = append(r.sqls, sql)
	return nil
}

// LastSQL - 查询出当前结果的sql集合中的最后一条sql
func (r *DefaultResult) LastSQL() string {
	sqlCnt := len(r.sqls)
	if sqlCnt == 0 {
		return ""
	}
	return r.sqls[sqlCnt-1]
}
