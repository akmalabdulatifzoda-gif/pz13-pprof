package work

// Fib — неоптимальный рекурсивный Фибоначчи: O(2^n) вызовов.
// Демонстрация CPU-нагрузки («до» оптимизации).
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

// FibFast — итеративная версия: O(n) времени, O(1) памяти («после»).
func FibFast(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
