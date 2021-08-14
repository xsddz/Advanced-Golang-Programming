package sql

// Field -
type Field struct {
	namePrefix string
	name       string
	nameAlias  string
	value      interface{}
	valueType  string
	operate    interface{} // 将name和value作为参数的过程
}

// Conditions -
type Conditions []Conditions

// AndConditions -
type AndConditions []Field

// OrConditions -
type OrConditions []Field
