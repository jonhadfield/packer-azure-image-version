package packer_azure_image_version

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
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

func getLatestImageVersion(s *session, builder Builder) (newVer *semver.Version, err error) {
	idi := ParseImageDefinitionID(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s",
		builder.SharedGalleryDestination.SigDestinationSubscription,
		builder.SharedGalleryDestination.SigDestinationResourceGroup,
		builder.SharedGalleryDestination.SigDestinationGalleryName,
		builder.SharedGalleryDestination.SigDestinationImageName,
	))

	if err := s.getGalleryImageVersionsClient(idi.SubscriptionID); err != nil {
		return newVer, err
	}

	rawVersions, err := getGalleryImageVersions(s, idi)
	if err != nil {
		return newVer, err
	}

	// default to 0.0.0 if no versions exist
	if len(rawVersions) == 0 {
		rawVersions = []string{"0.0.0"}
		return semver.MustParse("0.0.0"), err
	}

	vs := make([]*semver.Version, len(rawVersions))
	for i, r := range rawVersions {
		v, err := semver.NewVersion(r)
		if err != nil {
			return newVer, fmt.Errorf("error parsing version: %s", err)
		}

		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))
	logrus.Debugf("latest existing version: %s", vs[len(vs)-1].String())

	return vs[len(vs)-1], nil

}

func parseJSONTemplate(data []byte) (jt JSONTemplate, err error) {
	b := new(bytes.Buffer)
	_, err = b.Write(data)
	if err != nil {
		return jt, err
	}

	decoder := json.NewDecoder(b)
	if err = decoder.Decode(&jt); err != nil {
		return jt, fmt.Errorf("failed to process json: %+v", err)
	}

	var builder Builder
	switch len(jt.Builders) {
	case 0:
		return jt, fmt.Errorf("no builders found in packer file")
	case 1:
		builder = jt.Builders[0]
	default:
		return jt, fmt.Errorf("multiple builders defined but only one is currently supported")
	}

	logrus.Debugf("read SIG destination subscription: %s", builder.SharedGalleryDestination.SigDestinationSubscription)
	logrus.Debugf("read SIG destination resource group: %s", builder.SharedGalleryDestination.SigDestinationResourceGroup)
	logrus.Debugf("read SIG destination gallery name: %s", builder.SharedGalleryDestination.SigDestinationGalleryName)
	logrus.Debugf("read SIG destination image name: %s", builder.SharedGalleryDestination.SigDestinationImageName)
	logrus.Debugf("read SIG destination image version: %s", builder.SharedGalleryDestination.SigDestinationImageVersion)

	if !allDefined(builder.SharedGalleryDestination.SigDestinationSubscription,
		builder.SharedGalleryDestination.SigDestinationResourceGroup,
		builder.SharedGalleryDestination.SigDestinationGalleryName,
		builder.SharedGalleryDestination.SigDestinationImageName) {
		return jt, fmt.Errorf("shared gallery destination invalid or undefined")
	}

	return jt, nil
}

func incremenentSemVer(i SetImageVersionInput, t string, c semver.Version) (n semver.Version) {
	switch {
	case i.IncMajor:
		n = c.IncMajor()
		logrus.Debugf("incremented major with result: %s", n.String())
	case i.IncMinor:
		n = c.IncMinor()
		logrus.Debugf("incremented minor with result: %s", n.String())
	case i.IncPatch:
		n = c.IncPatch()
		logrus.Debugf("incremented patch with result: %s", n.String())
	}

	if n.String() == t {
		fmt.Printf("shared gallery destination version is already at desired version: %s\n", t)

		return n
	}

	return n
}

func updateJSONTemplate(v semver.Version, t *JSONTemplate, i SetImageVersionInput) {
	// we need to remove the builder's subscription id to prevent interactive oauth authentication
	if i.Unattended && t.Builders[0].SubscriptionID != "" {
		logrus.Print("stripping subscription_id from builder to allow for unattended (no oauth) build")
		t.Builders[0].SubscriptionID = ""
	}

	if i.CLIAuth {
		t.Builders[0].UseAzureCLIAuth = true
	}

	if i.publicIP {
		t.Builders[0].PrivateVirtualNetworkWithPublicIp = true
	}

	t.Builders[0].SharedGalleryDestination.SigDestinationImageVersion = v.String()

	logrus.Debugf("setting new image version to: %s", v.String())

	return
}

func encodeJSONTemplate(t JSONTemplate) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)
	e := json.NewEncoder(b)

	if b == nil {
		return b, errors.New("something went wrong")
	}

	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err = e.Encode(t); err != nil {
		return b, fmt.Errorf("failed to encode json: %+v", err)
	}

	return b, nil
}

func setPackerImageGalleryDestinationImageVersion(s *session, path string, i SetImageVersionInput) error {
	var err error

	fInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	m := fInfo.Mode()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var b *bytes.Buffer

	// if filename ends with json, then use json
	var newVer string
	switch filepath.Ext(path) {
	case ".json":
		jt, err := parseJSONTemplate(data)
		if err != nil {
			return err
		}
		liv, err := getLatestImageVersion(s, jt.Builders[0])
		if err != nil {
			return err
		}

		nv := incremenentSemVer(i, jt.Builders[0].SharedGalleryDestination.SigDestinationImageVersion, *liv)

		updateJSONTemplate(nv, &jt, i)

		b, err = encodeJSONTemplate(jt)
		if err != nil {
			return err
		}

		newVer = jt.Builders[0].SharedGalleryDestination.SigDestinationImageVersion
	case ".hcl":
		err = errors.New("not implemented")
	default:
		err = fmt.Errorf("unexpected file extension '%s'", filepath.Ext(path))
	}
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, m)
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

	if !i.Quiet {
		fmt.Printf("new shared destination gallery image version set to: %s\n", newVer)
	}

	return nil
}

type SetImageVersionInput struct {
	IncMajor   bool
	IncMinor   bool
	IncPatch   bool
	Unattended bool
	CLIAuth    bool
	publicIP   bool
	Paths      []string
	Quiet      bool
}
