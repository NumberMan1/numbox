package encryp

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type testStruct struct {
	Value string
}

func TestGeneratorAccessKey(t *testing.T) {
	token, err := GeneratorDefaultAccessKey(NewValueClaims(testStruct{
		Value: "test",
	}, 2*time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}

	var testObj testStruct
	err = ParseDefaultAccessToken(token, &testObj)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(testObj)

	time.Sleep(3 * time.Millisecond)

	err = ParseDefaultAccessToken(token, &testObj)
	if err == nil {
		t.Error(errors.New("unexpected valid"))
		return
	}
	fmt.Println(err.Error())
}
