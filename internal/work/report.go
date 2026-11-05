package work

import (
	"fmt"
	"strconv"
	"strings"
)

// BuildReport — «тяжёлая по памяти» функция для Heap-профиля.
// Склеивает строки через += и fmt.Sprintf: на каждой итерации
// создаётся новая строка, а старая копируется целиком — O(n²) байт аллокаций.
func BuildReport(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += fmt.Sprintf("note-%d;", i)
	}
	return s
}

// BuildReportFast — та же строка через strings.Builder с заранее
// выделенным буфером и strconv.AppendInt без промежуточных строк.
func BuildReportFast(n int) string {
	var sb strings.Builder
	sb.Grow(n * 12)
	var num [20]byte
	for i := 0; i < n; i++ {
		sb.WriteString("note-")
		sb.Write(strconv.AppendInt(num[:0], int64(i), 10))
		sb.WriteByte(';')
	}
	return sb.String()
}
