package compress

import (
	"fmt"
	"testing"
)

func stringAction(s *string) {
	*s = "ibas"
}
func TestStringAction(t *testing.T) {
	s := "haha"
	stringAction(&s)
	fmt.Println(s)
}
