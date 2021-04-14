package tinyredis

func processCommand( /* client *c */ ) {
	// addReply()
}

func initServer() {
	// 1. aeCreateEventLoop

	// 2. start listen

	// 3.
	// aeCreateFileEvent(server.el, server.ipfd[j], AE_READABLE, acceptTcpHandler,NULL) == AE_ERR)
	// aeCreateFileEvent(server.el, server.sofd, AE_READABLE, acceptUnixHandler,NULL) == AE_ERR)
	// aeCreateFileEvent(server.el, server.module_blocked_pipe[0], AE_READABLE, moduleBlockedClientPipeReadable,NULL) == AE_ERR)
}

func main() {
	// 1. if (background) daemonize();

	// 2. initServer();

	// 3. main loop
	// aeSetBeforeSleepProc(server.el, beforeSleep)
	// aeSetAfterSleepProc(server.el, afterSleep)
	// aeMain(server.el)
	// aeDeleteEventLoop(server.el)
}
