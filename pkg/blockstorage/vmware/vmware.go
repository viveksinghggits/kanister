package vmware

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vslm"

	"github.com/kanisterio/kanister/pkg/blockstorage"
)

var _ blockstorage.Provider = (*fcdProvider)(nil)

type fcdProvider struct {
	cli *vslm.Client
}

// NewProvider returns new provider for VMware FCDs.
func NewProvider(config map[string]string) (blockstorage.Provider, error) {
	u, err := soap.ParseURL(config["url"])
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get config")
	}
	cli, err := vim25.NewClient(context.TODO(), soap.NewClient(u, true))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create vim25 client")
	}
	vslmCli, err := vslm.NewClient(context.TODO(), cli)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create VSLM client")
	}
	p := &fcdProvider{
		cli: vslmCli,
	}
	return p, nil
}

func (*fcdProvider) Type() blockstorage.Type {
	return blockstorage.TypeFCD
}

func (*fcdProvider) VolumeCreate(context.Context, blockstorage.Volume) (*blockstorage.Volume, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) VolumeCreateFromSnapshot(ctx context.Context, snapshot blockstorage.Snapshot, tags map[string]string) (*blockstorage.Volume, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) VolumeDelete(context.Context, *blockstorage.Volume) error {
	return errors.New("Not implemented")
}

func (*fcdProvider) VolumeGet(ctx context.Context, id string, zone string) (*blockstorage.Volume, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) SnapshotCopy(ctx context.Context, from blockstorage.Snapshot, to blockstorage.Snapshot) (*blockstorage.Snapshot, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) SnapshotCreate(ctx context.Context, volume blockstorage.Volume, tags map[string]string) (*blockstorage.Snapshot, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) SnapshotCreateWaitForCompletion(context.Context, *blockstorage.Snapshot) error {
	return errors.New("Not implemented")
}

func (*fcdProvider) SnapshotDelete(context.Context, *blockstorage.Snapshot) error {
	return errors.New("Not implemented")
}

func (*fcdProvider) SnapshotGet(ctx context.Context, id string) (*blockstorage.Snapshot, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) SetTags(ctx context.Context, resource interface{}, tags map[string]string) error {
	return errors.New("Not implemented")
}

func (*fcdProvider) VolumesList(ctx context.Context, tags map[string]string, zone string) ([]*blockstorage.Volume, error) {
	return nil, errors.New("Not implemented")
}

func (*fcdProvider) SnapshotsList(ctx context.Context, tags map[string]string) ([]*blockstorage.Snapshot, error) {
	return nil, errors.New("Not implemented")
}
