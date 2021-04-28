package server

var RPCServiceTable = map[string]interface{}{}

func RegisteService(srvname string, f interface{}) {
	RPCServiceTable[srvname] = f
}

func DeleteService(srvname string) {
	delete(RPCServiceTable, srvname)
}
