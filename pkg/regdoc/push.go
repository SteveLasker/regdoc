package regdoc

import (
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func Push(contents []byte, ref string, resolver remotes.Resolver) error {
	ctx, store := newORASContext()

	desc := store.Add("", ContentLayerMediaType, contents)
	layers := []ocispec.Descriptor{desc}

	pushOpts := []oras.PushOpt{
		oras.WithConfigMediaType(ConfigMediaType),
		oras.WithNameValidation(nil),
	}

	manifest, err := oras.Push(ctx, resolver, ref, store, layers, pushOpts...)
	if err != nil {
		return err
	}

	printStderr("Pushed", ref)
	printStderr("Size:", readableSize(desc.Size))
	printStderr("Digest:", manifest.Digest)

	return nil
}
