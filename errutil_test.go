package errutil_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/kontora13-go/errutil"
	"golang.org/x/sync/errgroup"
)

func TestError(t *testing.T) {
	var err error
	var code = "INTERNAL"

	err = errutil.New("test user error", "one")
	log.Print("err := ", err.Error())
	log.Print("---")

	err = errutil.Newf("test user error: %v", "one")
	log.Print("err := ", err.Error())
	log.Print("---")

	err = errutil.NewWithCode(code, "test user error", "one")
	log.Print("err := ", err.Error())
	log.Print("---")

	err = errutil.NewWithCodef(code, "test user error: %v", "one")
	log.Print("err := ", err.Error())
	log.Print("---")

}

func TestErrorInGoroutine(t *testing.T) {
	err := errutil.New("test user error", "one")
	log.Print("usr := ", err)

	g, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 100; i++ {
		g.Go(func() error {
			_ = errutil.WithDevMessagef(err, fmt.Sprintf("err%v", i))
			//log.Print("err.dev := ", err.Error())

			return nil
		})
	}
	for i := 0; i < 100; i++ {
		g.Go(func() error {
			_ = errutil.WithMessagef(err, fmt.Sprintf("err%v", i))
			//log.Print("err.msg := ", err.Error())

			return nil
		})
	}

	if err = g.Wait(); err != nil {
		t.Error(err)
	}
}

func BenchmarkErrorDevMessage(b *testing.B) {
	err := errutil.New("test user error", "one")

	for i := 0; i < b.N; i++ {
		err = errutil.WithDevMessagef(err, "error")
	}
}

func BenchmarkErrorMessage(b *testing.B) {
	err := errutil.New("test user error", "one")

	for i := 0; i < b.N; i++ {
		err = errutil.WithMessagef(err, "error")
	}
}

func BenchmarkErrorStack(b *testing.B) {
	err := errutil.New("test user error", "one")

	for i := 0; i < b.N; i++ {
		err = errutil.WithStack(err)
	}
}
