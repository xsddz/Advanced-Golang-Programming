package httpclientcoder

import (
	"strings"
)

const (
	TPL_FILE = `// Code generated by http-client-coder. DO NOT EDIT.
// versions:
// 	http-client-coder: {{VERSION}}
// source: {{PROTO_FILE}}

package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/config"
	"yawebapp/library/infra/tools/httpclientcoder/base"
)`

	TPL_CLIENT = `
type {{SYMBOL}}Client struct {
	req base.HTTPRequest
}

func New{{SYMBOL}}Client() (*{{SYMBOL}}Client, error) {
	return New{{SYMBOL}}ClientWithAuth(base.DoNothingAuthHandler)
}

func New{{SYMBOL}}ClientWithAuth(authFunc base.AuthHandlerFunc) (*{{SYMBOL}}Client, error) {
	endpointname := "{{SYMBOL}}"

	var confmap map[string]base.ClientConf
	err := config.LoadConf(app.ConfPath()+"/httpclient.toml", &confmap)
	if err != nil {
		return nil, err
	}
	if _, ok := confmap[endpointname]; !ok {
		return nil, fmt.Errorf("no conf for %s", endpointname)
	}

	req := base.HTTPRequest{
		AuthFunc: authFunc,
		Method:   "get",
		Endpoint: confmap[endpointname].Endpoint,
		Path:     "/",
		Headers:  make(map[string]string),
	}

	return &{{SYMBOL}}Client{req}, nil
}`

	TPL_CLIENT_METHOD = `
func (c *{{CLIENT_SYMBOL}}Client) {{SYMBOL_METHOD}}(ctx context.Context, input Request{{REQ_SYMBOL}}) (output Response{{RESP_SYMBOL}}, err error) {
	request := c.req
	request.Ctx = ctx
	request.Method = "{{METHOD}}"
	request.Path = "{{SYMBOL}}"
	request.Headers = {{HEADER_MAP}}
	request.Body = input

	err = request.BindResponse(&output)
	return
}
`
	TPL_CLIENT_METHOD_PARAM      = `{{PARAM_PROP}} interface{} 'json:"{{PARAM}}"'`
	TPL_CLIENT_METHOD_WITH_PARAM = `
type Param{{PARAM_SYMBOL}} struct {
	{{PARAM_SYMBOL_PROPS}}
}
func (p Param{{PARAM_SYMBOL}}) Map() (m map[string]interface{}) {
	j, _ := json.Marshal(p)
	json.Unmarshal(j, &m)
	return
}

func (c *{{CLIENT_SYMBOL}}Client) {{SYMBOL_METHOD}}(ctx context.Context, param Param{{PARAM_SYMBOL}}, input Request{{REQ_SYMBOL}}) (output Response{{RESP_SYMBOL}}, err error) {
	request := c.req
	request.Ctx = ctx
	request.Method = "{{METHOD}}"
	request.Path = "{{SYMBOL}}"
	request.PathParam = param.Map()
	request.Headers = {{HEADER_MAP}}
	request.Body = input

	err = request.BindResponse(&output)
	return
}
`

	TPL_ENTITY_REQ      = `type Request{{SYMBOL}} struct {`
	TPL_ENTITY_REQ_PROP = `{{SYMBOL_PROP}} {{TYPE}} 'json:"{{SYMBOL}},omitempty" url:"{{SYMBOL}},omitempty"'`
	TPL_ENTITY_REQ_END  = `}

func (r Request{{SYMBOL}}) Json() string {
	j, _ := json.Marshal(r)
	return string(j)
}
func (r Request{{SYMBOL}}) URLEncoded() string {
	v, _ := query.Values(r)
	return v.Encode()
}`
	TPL_ENTITY_RESP      = `type Response{{SYMBOL}} struct {`
	TPL_ENTITY_RESP_PROP = `{{SYMBOL_PROP}} {{TYPE}} 'json:"{{SYMBOL}}"'`
	TPL_ENTITY_RESP_END  = `}`
)

type Coder struct{}

func NewCoder() (*Coder, error) {
	return &Coder{}, nil
}

func (c *Coder) Head(cmd Cmder) string {
	hc := cmd.(*HeadCmd)
	content := strings.ReplaceAll(TPL_FILE, "{{PROTO_FILE}}", hc.ProtoFile)
	content = strings.ReplaceAll(content, "{{VERSION}}", hc.CoderVer)
	return content + "\n"
}

func (c *Coder) Client(cmd Cmder) string {
	cc := cmd.(*ClientCmd)
	content := strings.ReplaceAll(TPL_CLIENT, "{{SYMBOL}}", cc.Symbol)
	return content + "\n"
}

func (c *Coder) Method(cmd Cmder) string {
	mc := cmd.(*ClientMethodCmd)
	cc := mc.BeginBlockCmd.(*ClientCmd)

	symbolMethod, paramProps := "", ""
	pathArr := strings.Split(strings.Trim(mc.Path, "/"), "/")
	for _, v := range pathArr {
		vLen := len(v)
		if v[0:1] == "{" && v[vLen-1:] == "}" {
			v = v[1 : vLen-1]
			prop := strings.ReplaceAll(TPL_CLIENT_METHOD_PARAM, "{{PARAM}}", v)
			prop = strings.ReplaceAll(prop, "{{PARAM_PROP}}", strings.Title(v))
			prop = strings.ReplaceAll(prop, "'", "`")
			paramProps += prop + "\n"
		}
		symbolMethod += strings.Title(v)
	}

	content := ""
	if len(paramProps) > 0 {
		content = strings.ReplaceAll(TPL_CLIENT_METHOD_WITH_PARAM, "{{CLIENT_SYMBOL}}", cc.Symbol)
	} else {
		content = strings.ReplaceAll(TPL_CLIENT_METHOD, "{{CLIENT_SYMBOL}}", cc.Symbol)
	}

	content = strings.ReplaceAll(content, "{{SYMBOL_METHOD}}", symbolMethod)
	content = strings.ReplaceAll(content, "{{PARAM_SYMBOL}}", cc.Symbol+symbolMethod)
	content = strings.ReplaceAll(content, "{{PARAM_SYMBOL_PROPS}}", paramProps)
	content = strings.ReplaceAll(content, "{{REQ_SYMBOL}}", mc.ReqSymbol)
	content = strings.ReplaceAll(content, "{{RESP_SYMBOL}}", mc.RespSymbol)

	content = strings.ReplaceAll(content, "{{METHOD}}", mc.Symbol)
	content = strings.ReplaceAll(content, "{{SYMBOL}}", mc.Path)

	headerMap := "map[string]string{"
	if headers, ok := mc.Tag["headers"]; ok {
		for _, item := range strings.Split(headers, ",") {
			if strings.HasPrefix(item, "content-type:") {
				headerMap += `"content-type":"` + strings.ReplaceAll(item, "content-type:", "") + `"`
			}
		}
	}
	headerMap += "}"
	content = strings.ReplaceAll(content, "{{HEADER_MAP}}", headerMap)

	return content + "\n"
}

func (c *Coder) Entity(cmd Cmder) string {
	ec := cmd.(*EntityCmd)
	switch ec.Type {
	case "req":
		content := strings.ReplaceAll(TPL_ENTITY_REQ, "{{SYMBOL}}", ec.Symbol)
		return content + "\n"
	case "resp":
		content := strings.ReplaceAll(TPL_ENTITY_RESP, "{{SYMBOL}}", ec.Symbol)
		return content + "\n"
	default:
		return ""
	}
}

func (c *Coder) EntityProp(cmd Cmder) string {
	epc := cmd.(*EntityPropCmd)
	ec := epc.BeginBlockCmd.(*EntityCmd)
	switch ec.Type {
	case "req":
		content := strings.ReplaceAll(TPL_ENTITY_REQ_PROP, "{{SYMBOL}}", epc.Symbol)
		content = strings.ReplaceAll(content, "{{TYPE}}", epc.Type)

		symbolProp := strings.ReplaceAll(epc.Symbol, "_", " ")
		symbolProp = strings.Title(symbolProp)
		symbolProp = strings.ReplaceAll(symbolProp, " ", "")
		content = strings.ReplaceAll(content, "{{SYMBOL_PROP}}", symbolProp)

		content = strings.ReplaceAll(content, "'", "`")

		return content + "\n"
	case "resp":
		content := strings.ReplaceAll(TPL_ENTITY_RESP_PROP, "{{SYMBOL}}", epc.Symbol)
		content = strings.ReplaceAll(content, "{{TYPE}}", epc.Type)

		symbolProp := strings.ReplaceAll(epc.Symbol, "_", " ")
		symbolProp = strings.Title(symbolProp)
		symbolProp = strings.ReplaceAll(symbolProp, " ", "")
		content = strings.ReplaceAll(content, "{{SYMBOL_PROP}}", symbolProp)

		content = strings.ReplaceAll(content, "'", "`")

		return content + "\n"
	default:
		return ""
	}
}

func (c *Coder) EndBlock(cmd Cmder) string {
	ebc := cmd.(*EndBlockCmd)
	switch bcmd := ebc.BeginBlockCmd.(type) {
	case *EntityCmd:
		switch bcmd.Type {
		case "req":
			content := strings.ReplaceAll(TPL_ENTITY_REQ_END, "{{SYMBOL}}", bcmd.Symbol)
			return content + "\n"
		case "resp":
			return TPL_ENTITY_RESP_END + "\n"
		default:
			return ""
		}
	default:
		return ""
	}
}