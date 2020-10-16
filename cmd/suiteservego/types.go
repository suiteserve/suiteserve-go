package suiteservego

import "time"

type TestEventAction string

const (
	TestEventActionRun    TestEventAction = "run"
	TestEventActionPause  TestEventAction = "pause"
	TestEventActionCont   TestEventAction = "cont"
	TestEventActionPass   TestEventAction = "pass"
	TestEventActionBench  TestEventAction = "bench"
	TestEventActionFail   TestEventAction = "fail"
	TestEventActionOutput TestEventAction = "output"
	TestEventActionSkip   TestEventAction = "skip"
)

type TestEvent struct {
	Time    time.Time
	Action  TestEventAction
	Package string
	Test    string
	Elapsed float64
	Output  string
}
