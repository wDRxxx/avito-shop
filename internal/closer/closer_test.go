package closer

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloser(t *testing.T) {
	exitFunc = func(code int) {}
	var wg sync.WaitGroup

	wg.Add(1)
	fn1 := func() error {
		defer wg.Done()
		return nil
	}
	fn2 := func() error {
		return nil
	}

	cl := New(&wg, os.Interrupt)
	SetGlobalCloser(cl)

	Add(1, fn1)
	Add(2, fn2)

	go CloseAll()

	wg.Wait()

	require.Len(t, cl.funcsStageOne, 0)
	require.Len(t, cl.funcsStageTwo, 0)
}
