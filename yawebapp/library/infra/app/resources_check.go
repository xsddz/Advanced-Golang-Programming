package app

import (
	"context"
	"fmt"
)

// resourcesCheckerTable 资源检查注册表，使用slice保证检查的先后
var resourcesCheckerTable = []resourcesChecker{
	{"db", checkDB},
	// {"cache", checkCache},
}

type resourcesChecker struct {
	name    string
	checker func()
}

func processResourcesChecker(c resourcesChecker) {
	fmt.Printf("check resources: [%v] \t......", c.name)
	c.checker()
	fmt.Printf("... [OK]\n")
}

func resourcesCheck() {
	for _, c := range resourcesCheckerTable {
		// TODO::通过appconf判断需要检查的资源
		// ...

		processResourcesChecker(c)
	}
}

//------------------------------------------------------------------------------

func checkDB() {
	for clusterName := range appDBConf {
		if err := DB(context.TODO(), clusterName).Ping(); err != nil {
			panic(fmt.Sprintf("ping cluster [%v] db error: %v", clusterName, err))
		}
	}
}

// func checkCache() {
// 	for clusterName := range appCacheConf {
// 		if _, err := Cache(clusterName).Ping().Result(); err != nil {
// 			panic(fmt.Sprintf("ping cluster [%v] cache error: %v", clusterName, err))
// 		}
// 	}
// }
