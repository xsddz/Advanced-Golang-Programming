package tinyredis

type redisServer struct {
	supervised bool // int supervised; /* 1 if supervised, 0 otherwise. */
	daemonize  bool // int daemonize; /* True if running as a daemon */
}

var server redisServer

func processCommand( /* client *c */ ) {
	// addReply()
}

// daemonize -
// should:
// 1. set umask
// 2. fork() != 0 and exit()
// 3. setsid()
// 4. current dir
// 5. close fds
// 6. redirect stdin/stdout/stderr to /dev/null
// 7. syslog or current log system
// 8. single instanse lock file
// 9. conf && signal && audo start set
func daemonize() {}

// check file /proc/sys/net/core/somaxconn value with server.tcp_backlog
func checkTcpBacklogSettings() {}

func initServer() {
	// 1.
	// server.el = aeCreateEventLoop(server.maxclients+CONFIG_FDSET_INCR);

	// 2. start listen
	/* Open the TCP listening socket for the user commands. */
	/* Open the listening Unix domain socket. */

	/* Create the Redis databases, and initialize other internal state. */
	/* Initialize the LRU keys pool. */

	// 3.
	// aeCreateTimeEvent(server.el, 1, serverCron, NULL, NULL)
	// aeCreateFileEvent(server.el, server.ipfd[j], AE_READABLE, acceptTcpHandler,NULL) == AE_ERR)
	// aeCreateFileEvent(server.el, server.sofd, AE_READABLE, acceptUnixHandler,NULL) == AE_ERR)
	// aeCreateFileEvent(server.el, server.module_blocked_pipe[0], AE_READABLE, moduleBlockedClientPipeReadable,NULL) == AE_ERR)

	// 4.
	/* Open the AOF file if needed. */

	// 5.
	/* 32 bit instances are limited to 4GB of address space, so if there is
	 * no explicit limit in the user provided configuration we set a limit
	 * at 3 GB using maxmemory with 'noeviction' policy'. This avoids
	 * useless crashes of the Redis instance for out of memory. */
}

func loadDataFromDisk() {
	// rdbLoad(server.rdb_filename,&rsi) == C_OK
}

func main() {
	background := server.daemonize && !server.supervised
	if background {
		daemonize()
	}

	initServer()
	checkTcpBacklogSettings()

	loadDataFromDisk()

	// 3. main loop
	// aeSetBeforeSleepProc(server.el, beforeSleep)
	// aeSetAfterSleepProc(server.el, afterSleep)
	// aeMain(server.el)
	// aeDeleteEventLoop(server.el)
}
