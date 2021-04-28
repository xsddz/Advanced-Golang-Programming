package stub

var ServiceTable = map[string]interface{}{}

func RegisteService(srvname string, f interface{}) {
	ServiceTable[srvname] = f
}

func DeleteService(srvname string) {
	delete(ServiceTable, srvname)
}
