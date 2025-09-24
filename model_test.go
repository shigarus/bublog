package bublog

import (
	// "slices"
	"testing"
)

type testWriter struct {
	written []byte
}

func (w *testWriter) Write(p []byte) (int, error) {
	w.written = append(w.written, p...)
	return len(p), nil
}

func TestModel(t *testing.T) {
	// toWrite := []byte("asdfg")
	//
	// someWriter := testWriter{[]byte("")}
	// m := New(&someWriter)
	// m.SetSize(4, 5)
	// m.Write(toWrite)
	// res := m.View()
	// want := "asdf\ng   \n    \n    \n    "
	// if res != want {
	// 	t.Errorf("Rendered incorrectly: got '%s' want '%s'", res, want)
	// }
	// if !slices.Equal(someWriter.written, toWrite) {
	// 	t.Errorf("Additional writer did not received data. Contains: %v, want: %v", someWriter.written, toWrite)
	// }
}

func TestEmpty(t *testing.T) {
	m := NewModel("")
	m.SetSize(1, 1)
	res := m.View()
	if res != " " {
		t.Errorf("Expect rendr one space, got %s", res)
	}
}
