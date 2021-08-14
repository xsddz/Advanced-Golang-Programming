package networking

func acceptTcpHandler( /*aeEventLoop *el, int fd, void *privdata, int mask*/ ) {
	for {
		// 1. anetTcpAccept(server.neterr, fd, cip, sizeof(cip), &cport);
		// 2. acceptCommonHandler(cfd,0,cip);
	}
}

func acceptCommonHandler( /* int fd, int flags, char *ip*/ ) {
	// 1. createClient(fd)
}

func createClient(fd int) {
	// 1. anetNonBlock(NULL,fd);
	// 2. anetEnableTcpNoDelay(NULL,fd);
	// 3. aeCreateFileEvent(server.el, fd, AE_READABLE, readQueryFromClient, c)
}

func readQueryFromClient( /* aeEventLoop *el, int fd, void *privdata, int mask */ ) {
	// 1. read(fd, c->querybuf+qblen, readlen);
	// 2. processInputBufferAndReplicate(c);
}

func processInputBufferAndReplicate( /* client *c */ ) {
	// processInputBuffer(c);
}

func processInputBuffer( /* client *c */ ) {
	// processCommand(c)
}
