package ossync_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/strider2038/ossync"
)

func TestNewGroup(t *testing.T) {
	group := ossync.NewGroup(context.Background())
	group.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	group.Go(func(ctx context.Context) error {
		return fmt.Errorf("error")
	})

	err := group.Wait()

	if err.Error() != "error" {
		t.Error("error is expected")
	}
}
