package percent

import (
	"testing"
)

func TestPercentageNonFiniteAndOverflow(t *testing.T) {
	for _, s := range []string{"NaN%", "+Inf%", "-Inf%", "1e100%"} {
		if _, err := ParsePercentageToPermille(s); err == nil {
			t.Fatal(s)
		}
	}
	if value, err := ParsePercentageToPermille("12.34%"); err != nil || value != 123 {
		t.Fatal(value, err)
	}
}
