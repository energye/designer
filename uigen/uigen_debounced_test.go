package uigen

import (
	"testing"
	"time"

	projBean "github.com/energye/designer/options/bean"
)

func TestFlushPendingSkipsStaleEpoch(t *testing.T) {
	debounceMutex.Lock()
	origEpoch := generationEpoch
	generationEpoch = 5
	timer := time.NewTimer(time.Hour)
	pendingGenerations[1] = &pendingGeneration{epoch: 4, timer: timer}
	debounceMutex.Unlock()
	t.Cleanup(func() {
		timer.Stop()
		debounceMutex.Lock()
		delete(pendingGenerations, 1)
		generationEpoch = origEpoch
		debounceMutex.Unlock()
	})

	flushPending(1, 4)

	debounceMutex.Lock()
	_, exists := pendingGenerations[1]
	debounceMutex.Unlock()
	if !exists {
		t.Fatal("过期定时器不应消费新的 pending 任务")
	}
}

func TestFlushAllInvalidatesPending(t *testing.T) {
	debounceMutex.Lock()
	origEpoch := generationEpoch
	timer := time.NewTimer(time.Hour)
	pendingGenerations[1] = &pendingGeneration{epoch: origEpoch, timer: timer}
	debounceMutex.Unlock()
	t.Cleanup(func() {
		timer.Stop()
		debounceMutex.Lock()
		delete(pendingGenerations, 1)
		generationEpoch = origEpoch
		debounceMutex.Unlock()
	})

	origProject := projBean.GProject
	projBean.GProject = &projBean.TProject{}
	t.Cleanup(func() { projBean.GProject = origProject })

	FlushAll()

	debounceMutex.Lock()
	count := len(pendingGenerations)
	epoch := generationEpoch
	debounceMutex.Unlock()
	if count != 0 {
		t.Fatalf("FlushAll 后不应残留 pending, got %d", count)
	}
	if epoch == origEpoch {
		t.Fatal("FlushAll 应推进 generationEpoch")
	}
}
