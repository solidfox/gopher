package spider

import (
	"testing"
)

func BenchmarkParsePage(t *testing.B) {
	Get30Pages()
}
