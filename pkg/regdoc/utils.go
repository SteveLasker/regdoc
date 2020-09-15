package regdoc

import (
	"context"
	"fmt"
	"os"

	orascontent "github.com/deislabs/oras/pkg/content"
	orascontext "github.com/deislabs/oras/pkg/context"
)

func printStderr(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func printStdout(a ...interface{}) {
	fmt.Fprintln(os.Stdout, a...)
}

func newORASContext() (context.Context, *orascontent.Memorystore) {
	ctx := orascontext.Background()
	memoryStore := orascontent.NewMemoryStore()
	return ctx, memoryStore
}

func readableSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
