package packer_azure_image_version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
)

func SetImageVersions(i SetImageVersionInput) error {
	s := session{}

	var err error

	for _, fp := range i.Paths {
		if err = setPackerImageGalleryDestinationImageVersion(&s, fp, i); err != nil {
			return err
		}
	}

	return nil
}

func setPackerImageGalleryDestinationImageVersion(s *session, path string, i SetImageVersionInput) error {
	var err error

	fInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	m := fInfo.Mode()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var jt JSONTemplate

	bb := new(bytes.Buffer)
	_, err = bb.Write(data)

	decoder := json.NewDecoder(bb)
	if err = decoder.Decode(&jt); err != nil {
		return fmt.Errorf("failed to process json: %+v", err)
	}

	var builder Builder
	switch len(jt.Builders) {
	case 0:
		return fmt.Errorf("no builders found in packer file")
	case 1:
		builder = jt.Builders[0]
	default:
		return fmt.Errorf("multiple builders defined but only one is currently supported")
	}

	logrus.Debugf("read SIG destination subscription: %s", builder.SharedGalleryDestination.SigDestinationSubscription)
	logrus.Debugf("read SIG destination resource group: %s", builder.SharedGalleryDestination.SigDestinationResourceGroup)
	logrus.Debugf("read SIG destination gallery name: %s", builder.SharedGalleryDestination.SigDestinationGalleryName)
	logrus.Debugf("read SIG destination image name: %s", builder.SharedGalleryDestination.SigDestinationImageName)
	logrus.Debugf("read SIG destination image version: %s", jt.Builders[0].SharedGalleryDestination.SigDestinationImageVersion)

	if !allDefined(builder.SharedGalleryDestination.SigDestinationSubscription,
		builder.SharedGalleryDestination.SigDestinationResourceGroup,
		builder.SharedGalleryDestination.SigDestinationGalleryName,
		builder.SharedGalleryDestination.SigDestinationImageName) {
		return fmt.Errorf("shared gallery destination invalid or undefined")
	}

	idi := ParseImageDefinitionID(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s",
		builder.SharedGalleryDestination.SigDestinationSubscription,
		builder.SharedGalleryDestination.SigDestinationResourceGroup,
		builder.SharedGalleryDestination.SigDestinationGalleryName,
		builder.SharedGalleryDestination.SigDestinationImageName,
	))

	if err := s.getGalleryImageVersionsClient(idi.SubscriptionID); err != nil {
		return err
	}

	rawVersions, err := getGalleryImageVersions(s, idi)
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

	var newVer semver.Version

	sort.Sort(semver.Collection(vs))
	logrus.Debugf("latest existing version: %s", vs[len(vs)-1].String())

	switch {
	case i.IncMajor:
		newVer = vs[len(vs)-1].IncMajor()
		logrus.Debugf("incremented major with result: %s", newVer.String())
	case i.IncMinor:
		newVer = vs[len(vs)-1].IncMinor()
		logrus.Debugf("incremented minor with result: %s", newVer.String())
	case i.IncPatch:
		newVer = vs[len(vs)-1].IncPatch()
		logrus.Debugf("incremented patch with result: %s", newVer.String())
	}

	if newVer.String() == builder.SharedGalleryDestination.SigDestinationImageVersion {
		fmt.Println("shared gallery destination version is already at desired version:", builder.SharedGalleryDestination.SigDestinationImageVersion)
		return nil
	}

	builder.SharedGalleryDestination.SigDestinationImageVersion = newVer.String()

	logrus.Debugf("setting new image version to: %s", newVer.String())

	jt.Builders[0] = builder

	b := new(bytes.Buffer)
	e := json.NewEncoder(b)

	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err = e.Encode(jt); err != nil {
		return fmt.Errorf("failed to encode json: %+v", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY, m)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing", path)
	}

	_, err = f.Write(b.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write updated json")
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %s", path)
	}

	if ! i.Quiet {
		fmt.Printf("new shared destination gallery image version set to: %s\n", newVer.String())
	}

	return nil
}

type SetImageVersionInput struct {
	IncMajor bool
	IncMinor bool
	IncPatch bool
	Paths    []string
	Quiet	 bool
}
