package evaluator

type ObjectType string

const (
	NullObject            ObjectType = "NULL"
	ReturnValueObject     ObjectType = "RETURN_VALUE"
	BooleanObject         ObjectType = "BOOLEAN"
	IntegerObject         ObjectType = "INTEGER"
	DecimalObject         ObjectType = "DECIMAL"
	StringObject          ObjectType = "STRING"
	DateTimeObject        ObjectType = "DATETIME"
	FunctionObject        ObjectType = "FUNCTION"
	ArrayObject           ObjectType = "ARRAY"
	HashObject            ObjectType = "HASH"
	FileObject            ObjectType = "FILE"
	BuiltInFunctionObject ObjectType = "BUILTIN_FUNCTION"
)

type Object interface {
	Type() ObjectType
	Debug() string
}
