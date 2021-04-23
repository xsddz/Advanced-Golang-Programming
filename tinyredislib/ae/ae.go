package ae

const (
	AE_OK  = 0
	AE_ERR = -1

	AE_NONE     = 0 /* No events registered. */
	AE_READABLE = 1 /* Fire when descriptor is readable. */
	AE_WRITABLE = 2 /* Fire when descriptor is writable. */
	AE_BARRIER  = 4 /* With WRITABLE, never fire the event if the
	   READABLE event already fired in the same event
	   loop iteration. Useful when you want to persist
	   things to disk before sending replies, and want
	   to do that in a group fashion. */

	AE_FILE_EVENTS      = 1
	AE_TIME_EVENTS      = 2
	AE_ALL_EVENTS       = (AE_FILE_EVENTS | AE_TIME_EVENTS)
	AE_DONT_WAIT        = 4
	AE_CALL_AFTER_SLEEP = 8

	AE_NOMORE           = -1
	AE_DELETED_EVENT_ID = -1
)

type time_t int64 // "ref: <sys/types.h>"

type aeFileProc func(eventLoop *aeEventLoop, fd int, clientData interface{}, mask int)
type aeTimeProc func(eventLoop *aeEventLoop, id int64, clientData interface{})
type aeEventFinalizerProc func(eventLoop *aeEventLoop, clientData interface{})
type aeBeforeSleepProc func(eventLoop *aeEventLoop)

type aeFileEvent struct{}
type aeTimeEvent struct{}
type aeFiredEvent struct{}

type aeEventLoop struct {
	maxfd           int /* highest file descriptor currently registered */
	setsize         int /* max number of file descriptors tracked */
	timeEventNextId int64
	lastTime        time_t        /* Used to detect system clock skew */
	events          *aeFileEvent  /* Registered events */
	fired           *aeFiredEvent /* Fired events */
	timeEventHead   *aeTimeEvent
	stop            int
	apidata         interface{} /* This is used for polling API specific data */
	beforesleep     aeBeforeSleepProc
	aftersleep      aeBeforeSleepProc
}

func aeCreateEventLoop(setsize int) *aeEventLoop { return &aeEventLoop{} }

func (e *aeEventLoop) aeDeleteEventLoop() {}
func (e *aeEventLoop) aeStop()            {}

func (e *aeEventLoop) aeCreateFileEvent(fd int, mask int, proc aeFileProc, clientData interface{}) int {
	return 0
}
func (e *aeEventLoop) aeDeleteFileEvent(fd int, mask int) {}
func (e *aeEventLoop) aeGetFileEvents(fd int) int         { return 0 }

func (e *aeEventLoop) aeCreateTimeEvent(milliseconds int64, proc aeTimeProc, clientData interface{}, finalizerProc aeEventFinalizerProc) int64 {
	return 0
}
func (e *aeEventLoop) aeDeleteTimeEvent(id int64) int { return 0 }

/* Process every pending time event, then every pending file event
 * (that may be registered by time event callbacks just processed).
 * Without special flags the function sleeps until some file event
 * fires, or when the next time event occurs (if any).
 *
 * If flags is 0, the function does nothing and returns.
 * if flags has AE_ALL_EVENTS set, all the kind of events are processed.
 * if flags has AE_FILE_EVENTS set, file events are processed.
 * if flags has AE_TIME_EVENTS set, time events are processed.
 * if flags has AE_DONT_WAIT set the function returns ASAP until all
 * if flags has AE_CALL_AFTER_SLEEP set, the aftersleep callback is called.
 * the events that's possible to process without to wait are processed.
 *
 * The function returns the number of events processed. */
func (e *aeEventLoop) aeProcessEvents(flags int) int { return 0 }

/* Wait for milliseconds until the given file descriptor becomes
 * writable/readable/exception */
func (e *aeEventLoop) aeWait(fd int, mask int, milliseconds int64) int { return 0 }

func (e *aeEventLoop) aeMain() {
	e.stop = 0
	for e.stop == 0 {
		if e.beforesleep != nil {
			e.beforesleep(e)
		}
		e.aeProcessEvents(AE_ALL_EVENTS | AE_CALL_AFTER_SLEEP)
	}
}

func (e *aeEventLoop) aeGetApiName() string                               { return "" }
func (e *aeEventLoop) aeSetBeforeSleepProc(beforesleep aeBeforeSleepProc) {}
func (e *aeEventLoop) aeSetAfterSleepProc(aftersleep aeBeforeSleepProc)   {}
func (e *aeEventLoop) aeGetSetSize() int                                  { return 0 }
func (e *aeEventLoop) aeResizeSetSize(setsize int) int                    { return 0 }
