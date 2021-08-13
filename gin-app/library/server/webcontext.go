package server

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebContext struct {
	*gin.Context
}

func NewWebContext(ctx *gin.Context) *WebContext {
	return &WebContext{Context: ctx}
}

func (wctx *WebContext) ShouldBind(obj interface{}) error {
	err := wctx.Context.ShouldBind(obj)
	if err != nil {
		return fmt.Errorf("参数错误：%v", err)
	}

	return nil
}

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
