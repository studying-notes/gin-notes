package main
//
//import (
//	"testing"
//	"time"
//)
//
//type request struct {
//	t   time.Time
//	n   int
//	act time.Time
//	ok  bool
//}
//
//// dFromDuration converts a duration to a multiple of the global constant d
//func dFromDuration(dur time.Duration) int {
//	// Adding a millisecond to be swallowed by the integer division
//	// because we don't care about small inaccuracies
//	return int((dur + time.Millisecond) / d)
//}
//
//// dSince returns multiples of d since t0
//func dSince(t time.Time) int {
//	return dFromDuration(t.Sub(t0))
//}
//
//func runReserve(t *testing.T, lim *Limiter, req request) *Reservation {
//	return runReserveMax(t, lim, req, InfDuration)
//}
//
//func runReserveMax(t *testing.T, lim *Limiter, req request, maxReserve time.Duration) *Reservation {
//	r := lim.reserveN(req.t, req.n, maxReserve)
//	if r.ok && (dSince(r.timeToAct) != dSince(req.act)) || r.ok != req.ok {
//		t.Errorf("lim.reserveN(t%d, %v, %v) = (t%d, %v) want (t%d, %v)",
//			dSince(req.t), req.n, maxReserve, dSince(r.timeToAct), r.ok, dSince(req.act), req.ok)
//	}
//	return &r
//}
//
//func TestSimpleReserve(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	runReserve(t, lim, request{t0, 2, t2, true})
//	runReserve(t, lim, request{t3, 2, t4, true})
//}
//
//func TestMix(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 3, t1, false}) // should return false because n > Burst
//	runReserve(t, lim, request{t0, 2, t0, true})
//	run(t, lim, []allow{{t1, 2, false}}) // not enought tokens - don't allow
//	runReserve(t, lim, request{t1, 2, t2, true})
//	run(t, lim, []allow{{t1, 1, false}}) // negative tokens - don't allow
//	run(t, lim, []allow{{t3, 1, true}})
//}
//
//func TestCancelInvalid(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	r := runReserve(t, lim, request{t0, 3, t3, false})
//	r.CancelAt(t0)                               // should have no effect
//	runReserve(t, lim, request{t0, 2, t2, true}) // did not get extra tokens
//}
//
//func TestCancelLast(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	r := runReserve(t, lim, request{t0, 2, t2, true})
//	r.CancelAt(t1) // got 2 tokens back
//	runReserve(t, lim, request{t1, 2, t2, true})
//}
//
//func TestCancelTooLate(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	r := runReserve(t, lim, request{t0, 2, t2, true})
//	r.CancelAt(t3) // too late to cancel - should have no effect
//	runReserve(t, lim, request{t3, 2, t4, true})
//}
//
//func TestCancel0Tokens(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	r := runReserve(t, lim, request{t0, 1, t1, true})
//	runReserve(t, lim, request{t0, 1, t2, true})
//	r.CancelAt(t0) // got 0 tokens back
//	runReserve(t, lim, request{t0, 1, t3, true})
//}
//
//func TestCancel1Token(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	r := runReserve(t, lim, request{t0, 2, t2, true})
//	runReserve(t, lim, request{t0, 1, t3, true})
//	r.CancelAt(t2) // got 1 token back
//	runReserve(t, lim, request{t2, 2, t4, true})
//}
//
//func TestCancelMulti(t *testing.T) {
//	lim := NewLimiter(10, 4)
//
//	runReserve(t, lim, request{t0, 4, t0, true})
//	rA := runReserve(t, lim, request{t0, 3, t3, true})
//	runReserve(t, lim, request{t0, 1, t4, true})
//	rC := runReserve(t, lim, request{t0, 1, t5, true})
//	rC.CancelAt(t1) // get 1 token back
//	rA.CancelAt(t1) // get 2 tokens back, as if C was never reserved
//	runReserve(t, lim, request{t1, 3, t5, true})
//}
//
//func TestReserveJumpBack(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t1, 2, t1, true}) // start at t1
//	runReserve(t, lim, request{t0, 1, t1, true}) // should violate Limit,Burst
//	runReserve(t, lim, request{t2, 2, t3, true})
//}
//
//func TestReserveJumpBackCancel(t *testing.T) {
//	lim := NewLimiter(10, 2)
//
//	runReserve(t, lim, request{t1, 2, t1, true}) // start at t1
//	r := runReserve(t, lim, request{t1, 2, t3, true})
//	runReserve(t, lim, request{t1, 1, t4, true})
//	r.CancelAt(t0)                               // cancel at t0, get 1 token back
//	runReserve(t, lim, request{t1, 2, t4, true}) // should violate Limit,Burst
//}
//
//func TestReserveSetLimit(t *testing.T) {
//	lim := NewLimiter(5, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	runReserve(t, lim, request{t0, 2, t4, true})
//	lim.SetLimitAt(t2, 10)
//	runReserve(t, lim, request{t2, 1, t4, true}) // violates Limit and Burst
//}
//
//func TestReserveSetBurst(t *testing.T) {
//	lim := NewLimiter(5, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	runReserve(t, lim, request{t0, 2, t4, true})
//	lim.SetBurstAt(t3, 4)
//	runReserve(t, lim, request{t0, 4, t9, true}) // violates Limit and Burst
//}
//
//func TestReserveSetLimitCancel(t *testing.T) {
//	lim := NewLimiter(5, 2)
//
//	runReserve(t, lim, request{t0, 2, t0, true})
//	r := runReserve(t, lim, request{t0, 2, t4, true})
//	lim.SetLimitAt(t2, 10)
//	r.CancelAt(t2) // 2 tokens back
//	runReserve(t, lim, request{t2, 2, t3, true})
//}
//
//func TestReserveMax(t *testing.T) {
//	lim := NewLimiter(10, 2)
//	maxT := d
//
//	runReserveMax(t, lim, request{t0, 2, t0, true}, maxT)
//	runReserveMax(t, lim, request{t0, 1, t1, true}, maxT)  // reserve for close future
//	runReserveMax(t, lim, request{t0, 1, t2, false}, maxT) // time to act too far in the future
//}
//
//type wait struct {
//	name   string
//	ctx    context.Context
//	n      int
//	delay  int // in multiples of d
//	nilErr bool
//}
//
//func runWait(t *testing.T, lim *Limiter, w wait) {
//	start := time.Now()
//	err := lim.WaitN(w.ctx, w.n)
//	delay := time.Now().Sub(start)
//	if (w.nilErr && err != nil) || (!w.nilErr && err == nil) || w.delay != dFromDuration(delay) {
//		errString := "<nil>"
//		if !w.nilErr {
//			errString = "<non-nil error>"
//		}
//		t.Errorf("lim.WaitN(%v, lim, %v) = %v with delay %v ; want %v with delay %v",
//			w.name, w.n, err, delay, errString, d*time.Duration(w.delay))
//	}
//}
//
//func TestWaitSimple(t *testing.T) {
//	lim := NewLimiter(10, 3)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	cancel()
//	runWait(t, lim, wait{"already-cancelled", ctx, 1, 0, false})
//
//	runWait(t, lim, wait{"exceed-burst-error", context.Background(), 4, 0, false})
//
//	runWait(t, lim, wait{"act-now", context.Background(), 2, 0, true})
//	runWait(t, lim, wait{"act-later", context.Background(), 3, 2, true})
//}
//
//func TestWaitCancel(t *testing.T) {
//	lim := NewLimiter(10, 3)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	runWait(t, lim, wait{"act-now", ctx, 2, 0, true}) // after this lim.tokens = 1
//	go func() {
//		time.Sleep(d)
//		cancel()
//	}()
//	runWait(t, lim, wait{"will-cancel", ctx, 3, 1, false})
//	// should get 3 tokens back, and have lim.tokens = 2
//	t.Logf("tokens:%v last:%v lastEvent:%v", lim.tokens, lim.last, lim.lastEvent)
//	runWait(t, lim, wait{"act-now-after-cancel", context.Background(), 2, 0, true})
//}
//
//func TestWaitTimeout(t *testing.T) {
//	lim := NewLimiter(10, 3)
//
//	ctx, cancel := context.WithTimeout(context.Background(), d)
//	defer cancel()
//	runWait(t, lim, wait{"act-now", ctx, 2, 0, true})
//	runWait(t, lim, wait{"w-timeout-err", ctx, 3, 0, false})
//}
//
//func TestWaitInf(t *testing.T) {
//	lim := NewLimiter(Inf, 0)
//
//	runWait(t, lim, wait{"exceed-burst-no-error", context.Background(), 3, 0, true})
//}
//
//func BenchmarkAllowN(b *testing.B) {
//	lim := NewLimiter(Every(1*time.Second), 1)
//	now := time.Now()
//	b.ReportAllocs()
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			lim.AllowN(now, 1)
//		}
//	})
//}
//
//func BenchmarkWaitNNoDelay(b *testing.B) {
//	lim := NewLimiter(Limit(b.N), b.N)
//	ctx := context.Background()
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		lim.WaitN(ctx, 1)
//	}
//}
