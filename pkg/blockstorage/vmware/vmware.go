package vmware

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vslm"

	"github.com/kanisterio/kanister/pkg/blockstorage"
)

var _ blockstorage.Provider = (*fcdProvider)(nil)

type fcdProvider struct {
	gom *vslm.GlobalObjectManager
}

const (
	defaultWaitTime = 10 * time.Minute

	noDescription = ""
)

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
	gom := vslm.NewGlobalObjectManager(vslmCli)
	p := &fcdProvider{
		gom: gom,
	}
	return p, nil
}

func (p *fcdProvider) Type() blockstorage.Type {
	return blockstorage.TypeFCD
}

func (p *fcdProvider) VolumeCreate(ctx context.Context, volume blockstorage.Volume) (*blockstorage.Volume, error) {
	spec := types.VslmCreateSpec{
		Name:         "disk-name",
		CapacityInMB: volume.Size / 1024,
	}
	task, err := p.gom.CreateDisk(ctx, spec)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create disk")
	}
	res, err := task.Wait(ctx, defaultWaitTime)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to wait on task")
	}
	obj, ok := res.(types.VStorageObject)
	if !ok {
		return nil, errors.New("Unknown type for result")
	}
	return p.VolumeGet(ctx, obj.Config.Id.Id, "")
}

func (p *fcdProvider) VolumeCreateFromSnapshot(ctx context.Context, snapshot blockstorage.Snapshot, tags map[string]string) (*blockstorage.Volume, error) {
	volID := ""
	name := ""
	task, err := p.gom.CreateDiskFromSnapshot(ctx, ID(volID), ID(snapshot.ID), name, nil, nil, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create disk from snapshot")
	}
	res, err := task.Wait(ctx, defaultWaitTime)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to wait on task")
	}
	_ = res
	return nil, errors.New("Not implemented")
}

func (p *fcdProvider) VolumeDelete(ctx context.Context, volume *blockstorage.Volume) error {
	task, err := p.gom.Delete(ctx, ID(volume.ID))
	if err != nil {
		return errors.Wrap(err, "Failed to delete the disk")
	}
	_, err = task.Wait(ctx, defaultWaitTime)
	return err
}

func (p *fcdProvider) VolumeGet(ctx context.Context, id string, zone string) (*blockstorage.Volume, error) {
	obj, err := p.gom.Retrieve(ctx, ID(id))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query the disk")
	}
	return convertFromObjectToVolume(obj), nil
}

func (p *fcdProvider) SnapshotCopy(ctx context.Context, from blockstorage.Snapshot, to blockstorage.Snapshot) (*blockstorage.Snapshot, error) {
	return nil, errors.New("Not implemented")
}

func (p *fcdProvider) SnapshotCreate(ctx context.Context, volume blockstorage.Volume, tags map[string]string) (*blockstorage.Snapshot, error) {
	task, err := p.gom.CreateSnapshot(ctx, ID(volume.ID), noDescription)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create snapshot")
	}
	res, err := task.Wait(ctx, defaultWaitTime)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to wait on task")
	}
	_ = res
	return nil, errors.New("Not implemented")
}

func (p *fcdProvider) SnapshotCreateWaitForCompletion(context.Context, *blockstorage.Snapshot) error {
	return errors.New("Not implemented")
}

func (p *fcdProvider) SnapshotDelete(ctx context.Context, snapshot *blockstorage.Snapshot) error {
	task, err := p.gom.DeleteSnapshot(ctx, ID(snapshot.Volume.ID), ID(snapshot.ID))
	if err != nil {
		return errors.Wrap(err, "Failed to delete snapshot")
	}
	_, err = task.Wait(ctx, defaultWaitTime)
	return err
}

func (p *fcdProvider) SnapshotGet(ctx context.Context, id string) (*blockstorage.Snapshot, error) {
	results, err := p.gom.RetrieveSnapshotInfo(ctx, ID(id))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get snapshot info")
	}
	if len(results) != 1 {
		return nil, errors.New("Wrong number of results")
	}
	return convertFromObjectToSnapshot(&results[0]), nil
}

func (p *fcdProvider) SetTags(ctx context.Context, resource interface{}, tags map[string]string) error {
	return errors.New("Not implemented")
}

func (p *fcdProvider) VolumesList(ctx context.Context, tags map[string]string, zone string) ([]*blockstorage.Volume, error) {
	return nil, errors.New("Not implemented")
}

func (p *fcdProvider) SnapshotsList(ctx context.Context, tags map[string]string) ([]*blockstorage.Snapshot, error) {
	return nil, errors.New("Not implemented")
}
