package context

import "sync"

// waitGroup provides hierarchical structure of WaitGroups,
// where children propagate Adds and Dones to the parents.
// A child Waits only for its own Adds, but a parent Waits
// for both its own and its childrens' Adds.
type waitGroup struct {
	parent *waitGroup
	wg     sync.WaitGroup
}

func newWaitGroup(parent *waitGroup) *waitGroup {
	return &waitGroup{parent: parent}
}

func (wg *waitGroup) Add(delta int) {
	if wg == nil {
		return
	}
	wg.wg.Add(delta)
	wg.parent.Add(delta)
}

func (wg *waitGroup) Done() {
	if wg == nil {
		return
	}
	wg.wg.Done()
	wg.parent.Done()
}

func (wg *waitGroup) Wait() {
	wg.wg.Wait()
}
