package errutil_test

import (
	"context"
	"log"
	"testing"

	"github.com/kontora13-go/errutil"
	"golang.org/x/sync/errgroup"
)

const (
	codeCritical = "CRITICAL"
)

func TestWrap(t *testing.T) {
	err := errutil.WithDevMessage(nil, "test error", "one")
	log.Print("err := ", err)
	log.Print("---")

	err = errutil.WithCode(err, codeCritical)
	log.Print("err := ", err)
	log.Print("---")

	err = errutil.WithMessage(err, "test user error", "one")
	log.Print("err := ", err)
	log.Print("---")

	log.Print("err := ", err.Error())
	log.Print("code := ", errutil.Code(err))
	log.Print("msg := ", errutil.Message(err))
	log.Print("dev := ", errutil.DevMessage(err))
	log.Print("cause := ", errutil.Cause(err))
	log.Print("stack_trace := ", string(errutil.Stack(err)))
}

func TestWrapInGoroutine(t *testing.T) {
	err := errutil.New("test user error", "one")

	g, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 100; i++ {
		g.Go(func() error {
			return errutil.New("test")
		})
	}

	if err = g.Wait(); err != nil {
		err = errutil.WithStack(err)
	}

	log.Print("err := ", err.Error())
	log.Print("stack_trace := ", errutil.Stack(err))
}
