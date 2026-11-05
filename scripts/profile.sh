#!/usr/bin/env bash
# Снимает профили под одинаковой нагрузкой для режима IMPL (slow | fast).
# Требует запущенный сервер (make run) и утилиту hey.
# Использование: scripts/profile.sh slow   или   scripts/profile.sh fast
set -euo pipefail

IMPL=${1:-slow}
BASE=${BASE:-http://localhost:8080}
DUR=${DUR:-15}               # длительность нагрузки и CPU-профиля, сек
CONC=${CONC:-8}              # параллельных клиентов
OUT=profiles/$IMPL
HEY=${HEY:-$(go env GOPATH)/bin/hey}
mkdir -p "$OUT"

echo "== [$IMPL] CPU: нагрузка /work, профиль ${DUR}s"
"$HEY" -z "$((DUR + 2))s" -c "$CONC" "$BASE/work?n=38&impl=$IMPL" > "$OUT/hey_work.txt" &
sleep 1
curl -s -o "$OUT/cpu.pb.gz" "$BASE/debug/pprof/profile?seconds=$DUR"
wait

# allocs/block/mutex накапливаются с запуска процесса, поэтому снимаем
# дельта-профили (?seconds=N) — только за время своей нагрузки.
echo "== [$IMPL] Heap: нагрузка /report, дельта allocs за 10s"
"$HEY" -z 12s -c "$CONC" "$BASE/report?n=2000&impl=$IMPL" > "$OUT/hey_report.txt" &
sleep 1
curl -s -o "$OUT/allocs.pb.gz" "$BASE/debug/pprof/allocs?seconds=10"
curl -s -o "$OUT/heap.pb.gz" "$BASE/debug/pprof/heap"
wait

echo "== [$IMPL] Block/Mutex: нагрузка /lock, дельта за 10s"
"$HEY" -z 12s -c 32 "$BASE/lock?impl=$IMPL" > "$OUT/hey_lock.txt" &
sleep 1
curl -s -o "$OUT/block.pb.gz" "$BASE/debug/pprof/block?seconds=10" &
curl -s -o "$OUT/mutex.pb.gz" "$BASE/debug/pprof/mutex?seconds=10"
curl -s -o "$OUT/goroutine.txt" "$BASE/debug/pprof/goroutine?debug=1"
wait

echo "== [$IMPL] Отчёты pprof"
go tool pprof -top -nodecount=15 "$OUT/cpu.pb.gz"                  > "$OUT/cpu_top.txt"   2>/dev/null
go tool pprof -list='work\.Fib' "$OUT/cpu.pb.gz"                   > "$OUT/cpu_list.txt"  2>/dev/null || true
go tool pprof -png -nodecount=20 "$OUT/cpu.pb.gz"                  > "$OUT/cpu_graph.png" 2>/dev/null
go tool pprof -sample_index=alloc_space -top -nodecount=12 "$OUT/allocs.pb.gz" > "$OUT/alloc_space_top.txt" 2>/dev/null
go tool pprof -sample_index=inuse_space -top -nodecount=12 "$OUT/heap.pb.gz"   > "$OUT/inuse_space_top.txt" 2>/dev/null
go tool pprof -sample_index=alloc_space -png -nodecount=15 "$OUT/allocs.pb.gz" > "$OUT/alloc_graph.png" 2>/dev/null
go tool pprof -top -nodecount=10 "$OUT/block.pb.gz"                > "$OUT/block_top.txt" 2>/dev/null
go tool pprof -top -nodecount=10 "$OUT/mutex.pb.gz"                > "$OUT/mutex_top.txt" 2>/dev/null

grep -E "Requests/sec|Average" "$OUT"/hey_*.txt
echo "Готово: $OUT/"
