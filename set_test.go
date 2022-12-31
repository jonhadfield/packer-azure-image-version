package packer_azure_image_version

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestParseJSONTemplate(t *testing.T) {
	data, err := os.ReadFile("testdata/example-one.json")
	require.NoError(t, err)
	require.NotEmpty(t, data)

	jt, err := parseJSONTemplate(data)
	require.NoError(t, err)
	require.Len(t, jt.Builders, 1)
	b := jt.Builders[0]
	require.True(t, b.UseAzureCLIAuth)
	require.Equal(t, "azure-arm", b.Type)
	require.Equal(t, "Linux", b.OSType)
	pi := b.PlanInfo
	require.Equal(t, "example-plan-name", pi.PlanName)
	require.Equal(t, "example-plan-product", pi.PlanProduct)
	require.Equal(t, "example-plan-publisher", pi.PlanPublisher)
	require.Equal(t, "example-plan-promotion-code", pi.PlanPromotionCode)
	require.Equal(t, "example-image-publisher", b.ImagePublisher)
	require.Equal(t, "example-image-sku", b.ImageSku)
	require.Equal(t, "example-image-offer", b.ImageOffer)
	require.Equal(t, "example-location", b.Location)
	require.Equal(t, "example-vm-size", b.VMSize)
	require.Equal(t, "example-managed-image-resource-group-name", b.ManagedImageResourceGroupName)
	require.Equal(t, "example-managed-image-name", b.ManagedImageName)
	tags := b.AzureTags
	require.Equal(t, "tag-one-value", tags["tag-one"])
	require.Equal(t, "tag-two-value", tags["tag-two"])
	require.Equal(t, "tag-three-value", tags["tag-three"])
	require.Equal(t, "tag-four-value", tags["tag-four"])
	require.True(t, b.PrivateVirtualNetworkWithPublicIp)
	require.Equal(t, "example-vnet-name", b.VirtualNetworkName)
	require.Equal(t, "example-subnet-name", b.VirtualNetworkSubnetName)
	require.Equal(t, "example-vnet-resource-group", b.VirtualNetworkResourceGroupName)
	sig := b.SharedGalleryDestination
	require.Equal(t, "b24988ac-6180-42a0-ab88-20f7382dd24c", sig.SigDestinationSubscription)
	require.Equal(t, "example-destination-resource-group", sig.SigDestinationResourceGroup)
	require.Equal(t, "example-destination-gallery", sig.SigDestinationGalleryName)
	require.Equal(t, "example-new-image-name", sig.SigDestinationImageName)
	require.Equal(t, "0.0.1", sig.SigDestinationImageVersion)
	require.Len(t, sig.SigDestinationReplicationRegions, 1)
	require.Equal(t, "canadacentral", sig.SigDestinationReplicationRegions[0])
	sigs := b.SharedGallery
	require.Equal(t, "c35099bd-32c2-6180-bc89-85t7382bb24c", sigs.Subscription)
	require.Equal(t, "example-source-resource-group", sigs.ResourceGroup)
	require.Equal(t, "example-source-gallery", sigs.GalleryName)
	require.Equal(t, "example-source-image-name", sigs.ImageName)
	require.Equal(t, "1.2.3", sigs.ImageVersion)
}
