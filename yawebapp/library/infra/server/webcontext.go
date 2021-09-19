package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"yawebapp/library/infra/helper"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var (
	// webcontextHookKeys 内置于webcontext，可用于链路追踪的keys
	webcontextHookKeys = []string{"trace_id", "Authorization"}
)

// WebContext 借由gin context提供能力
type WebContext struct {
	*gin.Context
}

func webcontextMaker(ctx *gin.Context, hookSetter func(*WebContext)) *WebContext {
	if ctx == nil {
		ctx = &gin.Context{}
	}
	wctx := &WebContext{Context: ctx}

	hookSetter(wctx)

	if wctx.GetString("trace_id") == "" {
		wctx.Set("trace_id", helper.NextRequestID())
	}

	return wctx
}

// NewWebContextViaHTTP 通过gin http请求创建webcontext
func NewWebContextViaHTTP(ctx *gin.Context) *WebContext {
	return webcontextMaker(ctx, func(wctx *WebContext) {
		for _, key := range webcontextHookKeys {
			if val := ctx.GetHeader(key); val != "" {
				wctx.Set(key, val)
			}
		}
	})
}

// NewWebContextViaContext 通过go context创建webcontext
func NewWebContextViaContext(ctx context.Context) *WebContext {
	return webcontextMaker(nil, func(wctx *WebContext) {
		for _, key := range webcontextHookKeys {
			if val := ctx.Value(key); val != nil {
				wctx.Set(key, val)
			}
		}
	})
}

// NewWebContextViaGRPC 通过grpc context创建webcontext
func NewWebContextViaGRPC(ctx context.Context) *WebContext {
	return webcontextMaker(nil, func(wctx *WebContext) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for key, val := range md {
				if helper.InArray(key, webcontextHookKeys) {
					wctx.Set(key, val[0])
				}
			}
		}
	})
}

func (wctx *WebContext) Token() string   { return wctx.GetString("Authorization") }
func (wctx *WebContext) TraceID() string { return wctx.GetString("trace_id") }

// ShouldBind gin http请求参数绑定
func (wctx *WebContext) ShouldBind(obj interface{}) error {
	err := wctx.Context.ShouldBind(obj)
	if err != nil {
		return fmt.Errorf("参数错误：%v", err)
	}

	return nil
}

// ShouldBindGRPC grpc请求参数绑定
func (wctx *WebContext) ShouldBindGRPC(data interface{}, obj interface{}) error {
	objT := reflect.TypeOf(obj)
	if objT.Kind() != reflect.Ptr {
		return fmt.Errorf("参数错误：类型非法")
	}

	dataJson, _ := json.Marshal(data)
	err := json.Unmarshal(dataJson, obj)
	if err != nil {
		return fmt.Errorf("参数错误：%v", err)
	}

	// 必传参数检查
	objET := objT.Elem()
	objEV := reflect.ValueOf(obj).Elem()
	for i := 0; i < objET.NumField(); i++ {
		fieldT := objET.Field(i)
		fieldV := objEV.Field(i)
		tags := strings.Split(fieldT.Tag.Get("binding"), ",")
		for _, t := range tags {
			if t == "required" {
				switch fieldV.Kind() {
				case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
					if fieldV.IsNil() {
						return fmt.Errorf("参数错误：必传参数 %v 缺失或零值", fieldT.Name)
					}
				default:
					if !fieldV.IsValid() || fieldV.Interface() == reflect.Zero(fieldV.Type()).Interface() {
						return fmt.Errorf("参数错误：必传参数 %v 缺失或零值", fieldT.Name)
					}
				}
				break
			}
		}
	}

	return nil
}
