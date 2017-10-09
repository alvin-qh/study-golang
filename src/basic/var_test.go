package basic

import (
	"testing"
)

func TestVar(t *testing.T) {
	t.Log("Hello")
	{
		t.Logf("OO")
		{
			t.Fatal()
		}
	}
}
