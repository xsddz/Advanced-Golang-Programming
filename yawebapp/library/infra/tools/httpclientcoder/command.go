package httpclientcoder

type CommandType int

const (
	CMD_ERROR           CommandType = iota
	CMD_NONE                        //
	CMD_COMMENT_INLINE              // // xxxxxx or # xxxxxx
	CMD_COMMENT                     // /** xxxxxxx
	CMD_COMMENT_CONTENT             // * xxxxx
	CMD_COMMENT_END                 // */
	CMD_CONF                        // conf xxxxx
	CMD_CLIENT                      // client xxxx {
	CMD_CLIENT_ENDPOINT             // endpoint xxxx
	CMD_CLIENT_METHOD               // {get|post|...} "path" reqEntity respEntity `headers:"xxxx:xxxx"`
	CMD_ENTITY                      // entity xxxx {
	CMD_ENTITY_PROP                 // xxxx {string|int|...}
	CMD_END_BLOCK                   // }
)

type Cmder interface{}

type HeadCmd struct {
	CoderVer  string
	ProtoFile string
}
type ClientCmd struct {
	Symbol string
}
type ClientMethodCmd struct {
	BeginBlockCmd Cmder
	Symbol        string
	Path          string
	ReqSymbol     string
	RespSymbol    string
	Tag           map[string]string
}
type EntityCmd struct {
	Symbol string
	Type   string
}
type EntityPropCmd struct {
	BeginBlockCmd Cmder
	Symbol        string
	Type          string
}
type EndBlockCmd struct {
	BeginBlockCmd Cmder
}
