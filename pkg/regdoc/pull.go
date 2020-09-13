package regdoc

import (
	"github.com/deislabs/oras/pkg/oras"
)

func Pull(ref string) error {
	ctx, resolver, store := newORASContext()

	pullOpts := []oras.PullOpt{
		oras.WithAllowedMediaType(ContentLayerMediaType),
		oras.WithPullEmptyNameAllowed(),
	}

	_, layers, err := oras.Pull(ctx, resolver, ref, store, pullOpts...)
	if err != nil {
		return err
	}

	desc := layers[0]
	manifest, contents, _ := store.Get(desc)

	printStderr("Pulled", ref)
	printStderr("Size:", readableSize(desc.Size))
	printStderr("Digest:", manifest.Digest)

	printStdout(string(contents))

	return nil
}