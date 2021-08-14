package ae

// timeval ref <sys/time.h>
type timeval struct{}

type aeAPI interface {
	// for example, use epoll_create()
	aeApiCreate(eventLoop *aeEventLoop) int

	aeApiResize(eventLoop *aeEventLoop, setsize int) int
	aeApiFree(eventLoop *aeEventLoop)

	// for example, use epoll_ctl()
	aeApiAddEvent(eventLoop *aeEventLoop, fd int, mask int) int
	aeApiDelEvent(eventLoop *aeEventLoop, fd int, delmask int)

	// for example, use epoll_wait()
	aeApiPoll(eventLoop *aeEventLoop, tvp *timeval)

	aeApiName() string
}
