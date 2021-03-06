// Copyright 2019 The Kanister Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zone

import (
	"context"
	"reflect"
	"testing"

	. "gopkg.in/check.v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	kubevolume "github.com/kanisterio/kanister/pkg/kube/volume"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ZoneSuite struct{}

var _ = Suite(&ZoneSuite{})

func (s ZoneSuite) TestConsistentZone(c *C) {
	// We don't care what the answer is as long as it's consistent.
	for _, tc := range []struct {
		sourceZone string
		nzs        map[string]struct{}
		out        string
	}{
		{
			sourceZone: "",
			nzs: map[string]struct{}{
				"zone1": {},
			},
			out: "zone1",
		},
		{
			sourceZone: "",
			nzs: map[string]struct{}{
				"zone1": {},
				"zone2": {},
			},
			out: "zone2",
		},
		{
			sourceZone: "from1",
			nzs: map[string]struct{}{
				"zone1": {},
				"zone2": {},
			},
			out: "zone1",
		},
	} {
		out, err := consistentZone(tc.sourceZone, tc.nzs, make(map[string]struct{}))
		c.Assert(err, IsNil)
		c.Assert(out, Equals, tc.out)
	}
}

func (s ZoneSuite) TestNodeZoneAndRegionGCP(c *C) {
	ctx := context.Background()
	node1 := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{kubevolume.PVRegionLabelName: "us-west2", kubevolume.PVZoneLabelName: "us-west2-a"},
		},
	}
	node2 := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node2",
			Labels: map[string]string{kubevolume.PVRegionLabelName: "us-west2", kubevolume.PVZoneLabelName: "us-west2-b"},
		},
	}
	node3 := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node3",
			Labels: map[string]string{kubevolume.PVRegionLabelName: "us-west2", kubevolume.PVZoneLabelName: "us-west2-c"},
		},
	}
	expectedZone := make(map[string]struct{})
	expectedZone["us-west2-a"] = struct{}{}
	expectedZone["us-west2-b"] = struct{}{}
	expectedZone["us-west2-c"] = struct{}{}
	cli := fake.NewSimpleClientset(node1, node2, node3)
	z, r, err := NodeZonesAndRegion(ctx, cli)
	c.Assert(err, IsNil)
	c.Assert(reflect.DeepEqual(z, expectedZone), Equals, true)
	c.Assert(r, Equals, "us-west2")
}

func (s ZoneSuite) TestNodeZoneAndRegionEBS(c *C) {
	ctx := context.Background()
	node1 := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{kubevolume.PVRegionLabelName: "us-west-2", kubevolume.PVZoneLabelName: "us-west-2a"},
		},
	}
	node2 := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node2",
			Labels: map[string]string{kubevolume.PVRegionLabelName: "us-west-2", kubevolume.PVZoneLabelName: "us-west-2b"},
		},
	}
	node3 := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node3",
			Labels: map[string]string{kubevolume.PVRegionLabelName: "us-west-2", kubevolume.PVZoneLabelName: "us-west-2c"},
		},
	}
	expectedZone := make(map[string]struct{})
	expectedZone["us-west-2a"] = struct{}{}
	expectedZone["us-west-2b"] = struct{}{}
	expectedZone["us-west-2c"] = struct{}{}
	cli := fake.NewSimpleClientset(node1, node2, node3)
	z, r, err := NodeZonesAndRegion(ctx, cli)
	c.Assert(err, IsNil)
	c.Assert(reflect.DeepEqual(z, expectedZone), Equals, true)
	c.Assert(r, Equals, "us-west-2")
}
