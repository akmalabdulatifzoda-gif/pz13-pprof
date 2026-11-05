package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // регистрирует /debug/pprof/*
	"os"
	"runtime"
	"strconv"

	"github.com/akmalabdulatifzoda-gif/pz13-pprof/internal/work"
)

// enableLocks включает сбор Block- и Mutex-профилей.
func enableLocks() {
	runtime.SetBlockProfileRate(1)     // записывать каждое событие блокировки
	runtime.SetMutexProfileFraction(1) // записывать каждое событие конкуренции за мьютекс
}

// fast возвращает true, если нужна оптимизированная версия.
// Режим задаётся параметром ?impl=fast или переменной окружения IMPL=fast.
func fast(r *http.Request) bool {
	if v := r.URL.Query().Get("impl"); v != "" {
		return v == "fast"
	}
	return os.Getenv("IMPL") == "fast"
}

func intParam(r *http.Request, key string, def int) int {
	if v, err := strconv.Atoi(r.URL.Query().Get(key)); err == nil && v >= 0 {
		return v
	}
	return def
}

func main() {
	enableLocks()
	store := work.NewStore()

	// http.DefaultServeMux — на нём net/http/pprof регистрирует /debug/pprof/.
	mux := http.DefaultServeMux

	// CPU-нагрузка: /work?n=38&impl=slow|fast
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		n := intParam(r, "n", 38) // достаточно тяжело для CPU
		var res int
		if fast(r) {
			defer work.TimeIt(fmt.Sprintf("FibFast(%d)", n))()
			res = work.FibFast(n)
		} else {
			defer work.TimeIt(fmt.Sprintf("Fib(%d)", n))()
			res = work.Fib(n)
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmtInt(res)))
	})

	// Нагрузка на память: /report?n=2000&impl=slow|fast
	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		n := intParam(r, "n", 2000)
		var s string
		if fast(r) {
			defer work.TimeIt(fmt.Sprintf("BuildReportFast(%d)", n))()
			s = work.BuildReportFast(n)
		} else {
			defer work.TimeIt(fmt.Sprintf("BuildReport(%d)", n))()
			s = work.BuildReport(n)
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintf(w, "%d bytes\n", len(s))
	})

	// Конкуренция за мьютекс: /lock?impl=slow|fast
	mux.HandleFunc("/lock", func(w http.ResponseWriter, r *http.Request) {
		if fast(r) {
			store.HitFast("lock")
		} else {
			store.HitSlow("lock")
		}
		_, _ = w.Write([]byte(fmtInt(store.Count("lock"))))
	})

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("Go %s, GOMAXPROCS=%d", runtime.Version(), runtime.GOMAXPROCS(0))
	log.Println("Server on " + addr + "; pprof on /debug/pprof/")
	log.Fatal(http.ListenAndServe(addr, mux))
}

func fmtInt(v int) string { return fmt.Sprintf("%d\n", v) }
