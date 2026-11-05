HEY ?= $(shell go env GOPATH)/bin/hey
BENCHSTAT ?= $(shell go env GOPATH)/bin/benchstat

.PHONY: tools run build test bench profile-slow profile-fast profile-all pprof-cpu pprof-heap clean

## tools: утилиты нагрузки и сравнения бенчмарков
tools:
	go install github.com/rakyll/hey@latest
	go install golang.org/x/perf/cmd/benchstat@latest

build:
	go build -o bin/api ./cmd/api

## run: сервер на :8080, pprof на /debug/pprof/
run:
	go run ./cmd/api

test:
	go test ./... -count=1

## bench: бенчмарки -benchmem (5 повторов) -> bench.txt
bench:
	go test ./internal/work -run=^$$ -bench=. -benchmem -count=5 | tee bench.txt

## bench-compare: сравнение «до/после» через benchstat
bench-compare: bench
	$(BENCHSTAT) bench.txt

## profile-slow / profile-fast: профили под одинаковой нагрузкой (сервер должен быть запущен)
profile-slow:
	HEY=$(HEY) ./scripts/profile.sh slow
profile-fast:
	HEY=$(HEY) ./scripts/profile.sh fast
profile-all: profile-slow profile-fast

## pprof-cpu / pprof-heap: веб-интерфейс pprof (Graph, Flame Graph, Top)
pprof-cpu:
	go tool pprof -http=:9999 profiles/slow/cpu.pb.gz
pprof-heap:
	go tool pprof -http=:9998 -sample_index=alloc_space profiles/slow/allocs.pb.gz

clean:
	rm -rf bin
