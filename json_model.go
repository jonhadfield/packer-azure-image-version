package packer_azure_image_version

import "time"

type ConfigWrapper struct {
	Builders     []Builder
	Provisioners interface{}
}

type ResourceID struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	Name           string
	Raw            string
}

type ImageDefinitionResourceID struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	Gallery        string
	ImageName      string
	Raw            string
}

type PackerConfig struct {
	PackerBuildName     string            `json:"packer_build_name,omitempty"`
	PackerBuilderType   string            `json:"packer_builder_type,omitempty"`
	PackerCoreVersion   string            `json:"packer_core_version,omitempty"`
	PackerDebug         bool              `json:"packer_debug,omitempty"`
	PackerForce         bool              `json:"packer_force,omitempty"`
	PackerOnError       string            `json:"packer_on_error,omitempty"`
	PackerUserVars      map[string]string `json:"packer_user_variables,omitempty"`
	PackerSensitiveVars []string          `json:"packer_sensitive_variables,omitempty"`
}

type Config struct {
	CloudEnvironmentName    string        `json:"cloud_environment_name,omitempty"`
	MetadataHost            string        `json:"metadata_host,omitempty"`
	ClientID                string        `json:"client_id,omitempty"`
	ClientSecret            string        `json:"client_secret,omitempty"`
	ClientCertPath          string        `json:"client_cert_path,omitempty"`
	ClientCertExpireTimeout time.Duration `json:"client_cert_token_timeout,omitempty"`
	ClientJWT               string        `json:"client_jwt,omitempty"`
	ObjectID                string        `json:"object_id,omitempty"`
	TenantID                string        `json:"tenant_id,omitempty"`
	SubscriptionID          string        `json:"subscription_id,omitempty"`
	UseAzureCLIAuth         bool          `json:"use_azure_cli_auth,omitempty"`
}

type SharedImageGallery struct {
	Subscription  string `json:"subscription,omitempty"`
	ResourceGroup string `json:"resource_group,omitempty"`
	GalleryName   string `json:"gallery_name,omitempty"`
	ImageName     string `json:"image_name,omitempty"`
	ImageVersion  string `json:"image_version,omitempty"`
}

type SharedImageGalleryDestination struct {
	SigDestinationSubscription       string   `json:"subscription,omitempty"`
	SigDestinationResourceGroup      string   `json:"resource_group,omitempty"`
	SigDestinationGalleryName        string   `json:"gallery_name,omitempty"`
	SigDestinationImageName          string   `json:"image_name,omitempty"`
	SigDestinationImageVersion       string   `json:"image_version,omitempty"`
	SigDestinationReplicationRegions []string `json:"replication_regions,omitempty"`
	SigDestinationStorageAccountType string   `json:"storage_account_type,omitempty"`
}

type PlanInformation struct {
	PlanName          string `json:"plan_name,omitempty"`
	PlanProduct       string `json:"plan_product,omitempty"`
	PlanPublisher     string `json:"plan_publisher,omitempty"`
	PlanPromotionCode string `json:"plan_promotion_code,omitempty"`
}

type NameValue struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Builder struct {
	PackerConfig
	Config
	Type                                       string                        `json:"type,omitempty"`
	UserAssignedManagedIdentities              []string                      `json:"user_assigned_managed_identities,omitempty"`
	CaptureNamePrefix                          string                        `json:"capture_name_prefix,omitempty"`
	CaptureContainerName                       string                        `json:"capture_container_name,omitempty"`
	SharedGallery                              SharedImageGallery            `json:"shared_image_gallery,omitempty"`
	SharedGalleryDestination                   SharedImageGalleryDestination `json:"shared_image_gallery_destination,omitempty"`
	SharedGalleryTimeout                       time.Duration                 `json:"shared_image_gallery_timeout,omitempty"`
	SharedGalleryImageVersionEndOfLifeDate     string                        `json:"shared_gallery_image_version_end_of_life_date,omitempty"`
	SharedGalleryImageVersionReplicaCount      int32                         `json:"shared_image_gallery_replica_count,omitempty"`
	SharedGalleryImageVersionExcludeFromLatest bool                          `json:"shared_gallery_image_version_exclude_from_latest,omitempty"`
	ImagePublisher                             string                        `json:"image_publisher"`
	ImageOffer                                 string                        `json:"image_offer"`
	ImageSku                                   string                        `json:"image_sku"`
	ImageVersion                               string                        `json:"image_version,omitempty"`
	ImageUrl                                   string                        `json:"image_url,omitempty"`
	CustomManagedImageName                     string                        `json:"custom_managed_image_name,omitempty"`
	CustomManagedImageResourceGroupName        string                        `json:"custom_managed_image_resource_group_name,omitempty"`
	Location                                   string                        `json:"location,omitempty"`
	VMSize                                     string                        `json:"vm_size,omitempty"`
	ManagedImageResourceGroupName              string                        `json:"managed_image_resource_group_name,omitempty"`
	ManagedImageName                           string                        `json:"managed_image_name,omitempty"`
	ManagedImageStorageAccountType             string                        `json:"managed_image_storage_account_type,omitempty"`
	ManagedImageOSDiskSnapshotName             string                        `json:"managed_image_os_disk_snapshot_name,omitempty"`
	ManagedImageDataDiskSnapshotPrefix         string                        `json:"managed_image_data_disk_snapshot_prefix,omitempty"`
	KeepOSDisk                                 bool                          `json:"keep_os_disk,omitempty,omitempty"`
	ManagedImageZoneResilient                  bool                          `json:"managed_image_zone_resilient,omitempty"`
	AzureTags                                  map[string]string             `json:"azure_tags,omitempty"`
	AzureTag                                   []NameValue                   `json:"azure_tag,omitempty"`
	ResourceGroupName                          string                        `json:"resource_group_name,omitempty"`
	StorageAccount                             string                        `json:"storage_account,omitempty"`
	TempComputeName                            string                        `json:"temp_compute_name,omitempty"`
	TempNicName                                string                        `json:"temp_nic_name,omitempty"`
	TempResourceGroupName                      string                        `json:"temp_resource_group_name,omitempty"`
	BuildResourceGroupName                     string                        `json:"build_resource_group_name,omitempty"`
	BuildKeyVaultName                          string                        `json:"build_key_vault_name,omitempty"`
	BuildKeyVaultSKU                           string                        `json:"build_key_vault_sku,omitempty"`
	PrivateVirtualNetworkWithPublicIp          bool                          `json:"private_virtual_network_with_public_ip,omitempty"`
	VirtualNetworkName                         string                        `json:"virtual_network_name,omitempty"`
	VirtualNetworkSubnetName                   string                        `json:"virtual_network_subnet_name,omitempty"`
	VirtualNetworkResourceGroupName            string                        `json:"virtual_network_resource_group_name,omitempty"`
	CustomDataFile                             string                        `json:"custom_data_file,omitempty"`
	UserDataFile                               string                        `json:"user_data_file,omitempty"`
	PlanInfo                                   PlanInformation               `json:"plan_info,omitempty"`
	PollingDurationTimeout                     time.Duration                 `json:"polling_duration_timeout,omitempty"`
	OSType                                     string                        `json:"os_type,omitempty"`
	TempOSDiskName                             string                        `json:"temp_os_disk_name,omitempty"`
	OSDiskSizeGB                               int32                         `json:"os_disk_size_gb,omitempty"`
	AdditionalDiskSize                         []int32                       `json:"disk_additional_size,omitempty"`
	DiskCachingType                            string                        `json:"disk_caching_type,omitempty"`
	AllowedInboundIpAddresses                  []string                      `json:"allowed_inbound_ip_addresses,omitempty"`
	BootDiagSTGAccount                         string                        `json:"boot_diag_storage_account,omitempty"`
	CustomResourcePrefix                       string                        `json:"custom_resource_build_prefix,omitempty"`
	UserName                                   string                        `json:"username,omitempty"`
	Password                                   string                        `json:"password,omitempty"`
	Comm                                       interface{}                   `json:"comm,omitempty"`
	AsyncResourceGroupDelete                   bool                          `json:"async_resourcegroup_delete,omitempty"`
}

type WinRM struct {
	WinRMUser     string        `json:"winrm_username,omitempty"`
	WinRMPassword string        `json:"winrm_password,omitempty"`
	WinRMHost     string        `json:"winrm_host,omitempty"`
	WinRMNoProxy  bool          `json:"winrm_no_proxy,omitempty"`
	WinRMPort     int           `json:"winrm_port,omitempty"`
	WinRMTimeout  time.Duration `json:"winrm_timeout,omitempty"`
	WinRMUseSSL   bool          `json:"winrm_use_ssl,omitempty"`
	WinRMInsecure bool          `json:"winrm_insecure,omitempty"`
	WinRMUseNTLM  bool          `json:"winrm_use_ntlm,omitempty"`
}

type SSH struct {
	SSHHost                   string        `json:"ssh_host,omitempty"`
	SSHPort                   int           `json:"ssh_port,omitempty"`
	SSHUsername               string        `json:"ssh_username,omitempty"`
	SSHPassword               string        `json:"ssh_password,omitempty"`
	SSHCiphers                []string      `json:"ssh_ciphers,omitempty"`
	SSHClearAuthorizedKeys    bool          `json:"ssh_clear_authorized_keys,omitempty"`
	SSHKEXAlgos               []string      `json:"ssh_key_exchange_algorithms,omitempty"`
	SSHPrivateKeyFile         string        `json:"ssh_private_key_file,omitempty"`
	SSHCertificateFile        string        `json:"ssh_certificate_file,omitempty"`
	SSHPty                    bool          `json:"ssh_pty,omitempty"`
	SSHTimeout                time.Duration `json:"ssh_timeout,omitempty"`
	SSHWaitTimeout            time.Duration `json:"ssh_wait_timeout,omitempty"`
	SSHAgentAuth              bool          `json:"ssh_agent_auth,omitempty"`
	SSHDisableAgentForwarding bool          `json:"ssh_disable_agent_forwarding,omitempty"`
	SSHHandshakeAttempts      int           `json:"ssh_handshake_attempts,omitempty"`
	SSHBastionHost            string        `json:"ssh_bastion_host,omitempty"`
	SSHBastionPort            int           `json:"ssh_bastion_port,omitempty"`
	SSHBastionAgentAuth       bool          `json:"ssh_bastion_agent_auth,omitempty"`
	SSHBastionUsername        string        `json:"ssh_bastion_username,omitempty"`
	SSHBastionPassword        string        `json:"ssh_bastion_password,omitempty"`
	SSHBastionInteractive     bool          `json:"ssh_bastion_interactive,omitempty"`
	SSHBastionPrivateKeyFile  string        `json:"ssh_bastion_private_key_file,omitempty"`
	SSHBastionCertificateFile string        `json:"ssh_bastion_certificate_file,omitempty"`
	SSHFileTransferMethod     string        `json:"ssh_file_transfer_method,omitempty"`
	SSHProxyHost              string        `json:"ssh_proxy_host,omitempty"`
	SSHProxyPort              int           `json:"ssh_proxy_port,omitempty"`
	SSHProxyUsername          string        `json:"ssh_proxy_username,omitempty"`
	SSHProxyPassword          string        `json:"ssh_proxy_password,omitempty"`
	SSHKeepAliveInterval      time.Duration `json:"ssh_keep_alive_interval,omitempty"`
	SSHReadWriteTimeout       time.Duration `json:"ssh_read_write_timeout,omitempty"`
	SSHRemoteTunnels          []string      `json:"ssh_remote_tunnels,omitempty"`
	SSHLocalTunnels           []string      `json:"ssh_local_tunnels,omitempty"`
	SSHPublicKey              []byte        `json:"ssh_public_key,omitempty"`
	SSHPrivateKey             []byte        `json:"ssh_private_key,omitempty"`
}

type JSONTemplate struct {
	MinPackerVersion   string            `json:"min_packer_version,omitempty"`
	Variables          map[string]string `json:"variables,omitempty"`
	SensitiveVariables []string          `json:"sensitive-variables,omitempty"`
	Builders           []Builder         `json:"builders,omitempty"`
	Provisioners       []interface{}     `json:"provisioners,omitempty"`
	PostProcessors     []interface{}     `json:"post-processors,omitempty"`
}
