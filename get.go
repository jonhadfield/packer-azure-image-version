package packer_azure_image_version

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
	"sort"
	"strings"
)

func getGalleryImageVersions(s *session, idi ImageDefinitionResourceID) (versions []string, err error) {
	ctx := context.Background()
	var it compute.GalleryImageVersionListIterator
	logrus.Debugf("getting image versions for: %s", idi.Raw)
	c := s.galleryImageVersionsClient[idi.SubscriptionID]
	it, err = c.ListByGalleryImageComplete(ctx, idi.ResourceGroup, idi.Gallery, idi.ImageName)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return nil, err
	}

	for it.NotDone() {
		f := it.Value()
		versions = append(versions, *f.Name)

		err = it.NextWithContext(ctx)
		if err != nil {
			return
		}
	}

	logrus.Debugf("read existing image versions: %s", strings.Join(versions, ", "))

	return
}

func GetImageVersions(i GetImageVersionsInput) error {
	s := session{}

	idi := ParseImageDefinitionID(i.ImageDefinitionID)

	if err := s.getGalleryImageVersionsClient(idi.SubscriptionID); err != nil {
		return err
	}

	rawVersions, err := getGalleryImageVersions(&s, idi)
	if err != nil {
		return err
	}

	vs := make([]*semver.Version, len(rawVersions))
	for i, r := range rawVersions {
		v, err := semver.NewVersion(r)
		if err != nil {
			return fmt.Errorf("error parsing version: %s", err)
		}

		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))
	switch {
	case i.Oldest:
		fmt.Println(vs[0])
	case i.Latest:
		fmt.Println(vs[len(vs)-1])
	case i.IncMajor:
		fmt.Println(vs[len(vs)-1].IncMajor())
	case i.IncMinor:
		fmt.Println(vs[len(vs)-1].IncMinor())
	case i.IncPatch:
		fmt.Println(vs[len(vs)-1].IncPatch())
	default:
		for x := range vs {
			fmt.Println(vs[x])
		}
	}

	return nil
}


type GetImageVersionsInput struct {
	SubscriptionID    string
	Latest            bool
	Oldest            bool
	IncMajor          bool
	IncMinor          bool
	IncPatch          bool
	ImageDefinitionID string
}