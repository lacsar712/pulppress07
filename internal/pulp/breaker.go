package pulp

import "sync"

type TripBreaker struct {
	mu       sync.Mutex
	failures int
	open     bool
	limit    int
}

func NewTripBreaker(limit int) *TripBreaker {
	if limit <= 0 {
		limit = 2
	}
	return &TripBreaker{limit: limit}
}

func (b *TripBreaker) Fail() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.limit {
		b.open = true
	}
}

func (b *TripBreaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.open = false
}

func (b *TripBreaker) Open() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.open
}

// RecordOutcome 把一次出包结果记入压区跳闸熔断。
//
// 断点标记（工单 #P-44）：此处曾对任何非空 err 一律 Fail()，把
// 「可再试」回包与 terminal 失败同等计数，推动压区被掐死。改为
// 只对 terminal（不可恢复）失败计数；ok 与 retry 都不推动熔断，
// 让重试语义与分类入口一致放行。
func RecordOutcome(b *TripBreaker, err error) {
	switch ClassifyOutcome(err) {
	case "terminal":
		b.Fail()
	case "retry", "ok":
		b.Success()
	}
}
