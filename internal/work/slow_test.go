package work

import (
	"strings"
	"sync"
	"testing"
)

func TestFibFastMatchesFib(t *testing.T) {
	for n := 0; n <= 30; n++ {
		if got, want := FibFast(n), Fib(n); got != want {
			t.Fatalf("FibFast(%d) = %d, want %d", n, got, want)
		}
	}
	if got := FibFast(38); got != 39088169 {
		t.Fatalf("FibFast(38) = %d", got)
	}
}

func TestBuildReportFastMatches(t *testing.T) {
	for _, n := range []int{0, 1, 10, 1000} {
		if BuildReport(n) != BuildReportFast(n) {
			t.Fatalf("BuildReport(%d) mismatch", n)
		}
	}
	if !strings.HasPrefix(BuildReportFast(3), "note-0;note-1;note-2;") {
		t.Fatal("unexpected format")
	}
}

func TestStoreCounts(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); s.HitSlow("a") }()
		go func() { defer wg.Done(); s.HitFast("a") }()
	}
	wg.Wait()
	if s.Count("a") != 40 {
		t.Fatalf("count = %d, want 40", s.Count("a"))
	}
}

// --- CPU: Фибоначчи ---

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Fib(30) // поменьше, чтобы бенч шёл быстрее
	}
}

func BenchmarkFibFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibFast(30)
	}
}

// --- Память: построение отчёта ---

func BenchmarkBuildReport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BuildReport(2000)
	}
}

func BenchmarkBuildReportFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BuildReportFast(2000)
	}
}

// --- Блокировки: 8 горутин одновременно ---

func benchStore(b *testing.B, hit func(*Store, string)) {
	s := NewStore()
	b.SetParallelism(4) // 4 × GOMAXPROCS горутин
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hit(s, "k")
		}
	})
}

func BenchmarkStoreHitSlow(b *testing.B) { benchStore(b, (*Store).HitSlow) }
func BenchmarkStoreHitFast(b *testing.B) { benchStore(b, (*Store).HitFast) }
