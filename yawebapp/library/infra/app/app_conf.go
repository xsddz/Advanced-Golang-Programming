package app

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"yawebapp/library/infra/config"
	"yawebapp/library/infra/storage"
)

// LoadConf 加载toml配置文件内容至obj对应的数据结构中
// 如果使用apollo，会根据配置文件中指示的key，进一步获取apollo的配置
func LoadConf(filename string, obj interface{}) error {
	file := fmt.Sprintf("%v/%v.toml", ConfPath(), filename)
	err := config.LoadConf(file, obj)
	if err != nil {
		return fmt.Errorf("load conf [%v] error: %v", file, err)
	}

	if UseApollo() {
		reloadFromApollo(reflect.ValueOf(obj))
	}

	return nil
}

func reloadFromApollo(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Ptr:
		reloadFromApollo(rv.Elem())
	case reflect.Map:
		reloadMapValueFromApollo(rv)
	case reflect.Struct:
		reloadStructValueFromApollo(rv)
	case reflect.String:
		reloadStringValueFromApollo(rv)
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
	case reflect.Float32, reflect.Float64:
	default:
	}

	return nil
}

func reloadMapValueFromApollo(rv reflect.Value) error {
	for _, k := range rv.MapKeys() {
		v := rv.MapIndex(k)
		switch v.Kind() {
		case reflect.Ptr:
			reloadFromApollo(v)
		case reflect.Map:
			reloadMapValueFromApollo(v)
		case reflect.Struct:
			s := reflect.New(v.Type())
			s.Elem().Set(v)
			reloadStructValueFromApollo(s.Elem())
			rv.SetMapIndex(k, s.Elem())
		case reflect.String:
			reloadStringValueFromApollo(v)
		}
	}

	return nil
}

func reloadStructValueFromApollo(rv reflect.Value) error {
	for i := 0; i < rv.NumField(); i++ {
		v := rv.Field(i)
		switch v.Kind() {
		case reflect.Ptr:
			reloadFromApollo(v)
		case reflect.Map:
			reloadMapValueFromApollo(v)
		case reflect.Struct:
			reloadStructValueFromApollo(v)
		case reflect.String:
			reloadStringValueFromApollo(v)
		}
	}
	return nil
}

func reloadStringValueFromApollo(rv reflect.Value) {
	newVal := Apollo().Get(rv.String())
	if newVal != "" {
		rv.SetString(newVal)
	}
}

////////////////////////////////////////////////////////////////////////////////
// confLoaderTable 配置加载注册表，使用slice保证加载的先后
var confLoaderTable = []confLoader{
	{"apollo", loadApolloConf},
	{"app", loadAppConf},
	{"db", loadDBConf},
	{"cache", loadCacheConf},
}

type confLoader struct {
	name   string
	loader func()
}

func processConfLoader(l confLoader) {
	fmt.Printf("load conf: [%v] \t......", l.name)
	l.loader()
	fmt.Printf("... [OK]\n")
}

func loadConf() {
	for _, l := range confLoaderTable {
		// TODO::通过appconf判断需要加载的配置
		// ...

		processConfLoader(l)
	}
}

//------------------------------------------------------------------------------
var appApolloConf *config.ApolloConf

// UseApollo 是否使用apollo
func UseApollo() bool { return ConfApollo().UseApollo }

// ConfApollo 获取apollo服务的配置
func ConfApollo() config.ApolloConf {
	if appApolloConf != nil {
		return *appApolloConf
	}

	panic("must load apollo conf before use")
}

func loadApolloConf() {
	err := config.LoadConf(ConfPath()+"/apollo.toml", &appApolloConf)
	if err != nil {
		panic(fmt.Sprint("[loadApolloConf] init agollo conf error: ", err))
	}
}

//------------------------------------------------------------------------------
var appConf *AppConf

type AppConf struct {
	AppName     string `toml:"app_name"`
	AppEnv      string `toml:"app_env"`
	OutputColor string `toml:"output_color"`
	Server      map[string]struct {
		Enable string `toml:"enable"`
		Port   string `toml:"port"`
	} `toml:"server"`
}

// OutputColor 颜色输出？仅对终端有效
func OutputColor() bool {
	switch ConfApp().OutputColor {
	case "true", "True", "TRUE", "1":
		return true
	default:
		return false
	}
}

// ConfApp 获取app配置
func ConfApp() AppConf {
	if appConf != nil {
		return *appConf
	}

	panic("must load app conf before use")
}

func loadAppConf() {
	err := LoadConf("app", &appConf)
	if err != nil {
		panic(fmt.Sprint("[loadAppConf] load app conf error: ", err))
	}
}

//------------------------------------------------------------------------------
var appDBConf map[string]*storage.DBConf

// ConfDB 获取db配置
func ConfDB(clusterNames ...string) storage.DBConf {
	clusterName := "Default"
	if len(clusterNames) > 0 {
		clusterName = clusterNames[0]
	}

	if c, ok := appDBConf[clusterName]; ok && c != nil {
		return *c
	}

	panic(fmt.Sprintf("must load db cluster [%v] conf before use", clusterName))
}

func loadDBConf() {
	err := LoadConf("db", &appDBConf)
	if err != nil {
		panic(fmt.Sprint("[loadDBConf] load db conf error: ", err))
	}

	for key := range appDBConf {
		for _, hp := range strings.Split(appDBConf[key].HostPorts, ",") {
			hpArr := strings.Split(hp, ":")
			if len(hpArr) == 2 {
				port, _ := strconv.Atoi(hpArr[1])
				appDBConf[key].Hosts = append(appDBConf[key].Hosts, storage.DBConfHost{IP: hpArr[0], Port: port})
			}
		}
	}
}

//------------------------------------------------------------------------------
var appCacheConf map[string]*storage.RedisConf

// ConfCache 获取缓存配置
func ConfCache(clusterNames ...string) storage.RedisConf {
	clusterName := "Default"
	if len(clusterNames) > 0 {
		clusterName = clusterNames[0]
	}

	if c, ok := appCacheConf[clusterName]; ok && c != nil {
		return *c
	}

	panic(fmt.Sprintf("must load cache cluster [%v] conf before use", clusterName))
}

func loadCacheConf() {
	if err := LoadConf("cache", &appCacheConf); err != nil {
		panic(fmt.Sprint("[loadCacheConf] load cache conf error: ", err))
	}

	for key := range appCacheConf {
		for _, hp := range strings.Split(appCacheConf[key].HostPorts, ",") {
			hpArr := strings.Split(hp, ":")
			if len(hpArr) == 2 {
				port, _ := strconv.Atoi(hpArr[1])
				appCacheConf[key].Hosts = append(appCacheConf[key].Hosts, storage.RedisConfHost{IP: hpArr[0], Port: port})
			}
		}
	}
}
