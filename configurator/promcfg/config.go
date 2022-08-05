package promcfg

import (
	"net/url"
	"time"

	"github.com/alecthomas/units"
)

// ExtraConfig represents some configuration which will be included in prometheus as it is.

// TLSConfig represents tls configuration, `prometheusCommonConfig.TLSConfig` cannot be used directly
// because it does not Marshal to yaml properly.
type TLSConfig struct {
	CAFile             string `yaml:"ca_file,omitempty" json:"ca_file,omitempty"`
	CertFile           string `yaml:"cert_file,omitempty" json:"cert_file,omitempty"`
	KeyFile            string `yaml:"key_file,omitempty" json:"key_file,omitempty"`
	ServerName         string `yaml:"server_name,omitempty" json:"server_name,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
	MinVersion         string `yaml:"min_version,omitempty" json:"min_version,omitempty"`
}

// Authorization holds prometheus authorization information.
type Authorization struct {
	Type            string `yaml:"type,omitempty"`
	Credentials     string `yaml:"credentials,omitempty"`
	CredentialsFile string `yaml:"credentials_file,omitempty"`
}

// BasicAuth defines the config for the `Authorization` header on every scrape request.
type BasicAuth struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password,omitempty"`
	PasswordFile string `yaml:"password_file,omitempty"`
}

// OAuth2 defines the config for prometheus to gather a token from the endpoint.
type OAuth2 struct {
	ClientID         string            `yaml:"client_id"`
	ClientSecret     string            `yaml:"client_secret,omitempty"`
	ClientSecretFile string            `yaml:"client_secret_file,omitempty"`
	Scopes           []string          `yaml:"scopes,omitempty"`
	TokenURL         string            `yaml:"token_url"`
	EndpointParams   map[string]string `yaml:"endpoint_params,omitempty"`
	TLSConfig        *TLSConfig        `yaml:"tls_config,omitempty"`
	ProxyURL         string            `yaml:"proxy_url,omitempty"`
}

// Job holds fields which do not change from input and output jobs.
type Job struct {
	JobName               string           `yaml:"job_name"`
	HonorLabels           bool             `yaml:"honor_labels,omitempty"`
	HonorTimestamps       *bool            `yaml:"honor_timestamps,omitempty"`
	Params                url.Values       `yaml:"params,omitempty"`
	Scheme                string           `yaml:"scheme,omitempty"`
	BodySizeLimit         units.Base2Bytes `yaml:"body_size_limit,omitempty"`
	SampleLimit           uint             `yaml:"sample_limit,omitempty"`
	TargetLimit           uint             `yaml:"target_limit,omitempty"`
	LabelLimit            uint             `yaml:"label_limit,omitempty"`
	LabelNameLengthLimit  uint             `yaml:"label_name_length_limit,omitempty"`
	LabelValueLengthLimit uint             `yaml:"label_value_length_limit,omitempty"`
	MetricsPath           string           `yaml:"metrics_path,omitempty"`
	ScrapeInterval        time.Duration    `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout         time.Duration    `yaml:"scrape_timeout,omitempty"`
	TLSConfig             *TLSConfig       `yaml:"tls_config,omitempty"`
	BasicAuth             *BasicAuth       `yaml:"basic_auth,omitempty"`
	Authorization         Authorization    `yaml:"authorization,omitempty"`
	OAuth2                OAuth2           `yaml:"oauth2,omitempty"`

	StaticConfigs        []StaticConfig       `yaml:"static_configs,omitempty"`
	RelabelConfigs       []RelabelConfig      `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs []RelabelConfig      `yaml:"metric_relabel_configs,omitempty"`
	KubernetesSdConfigs  []KubernetesSdConfig `yaml:"kubernetes_sd_configs,omitempty"`
}

// StaticConfig defines each of the static_configs for the prometheus config.
type StaticConfig struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

// GlobalConfig configures values that are used across other configuration
// objects.
type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval time.Duration `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout time.Duration `yaml:"scrape_timeout,omitempty"`
	// The labels to add to any timeseries that this Prometheus instance scrapes.
	ExternalLabels map[string]string `yaml:"external_labels,omitempty"`
}

// RelabelConfig defines relabel config rules which can be used in other configuration objects.
type RelabelConfig struct {
	SourceLabels []string `yaml:"source_labels,omitempty"`
	Separator    string   `yaml:"separator,omitempty"`
	TargetLabel  string   `yaml:"target_label,omitempty"`
	Regex        string   `yaml:"regex,omitempty"`
	Modulus      int      `yaml:"modulus,omitempty"`
	Replacement  string   `yaml:"replacement,omitempty"`
	Action       string   `yaml:"action,omitempty"`
}

// KubernetesSdConfig defines the kubernetes service discovery config.
type KubernetesSdConfig struct {
	Role           string                  `yaml:"role,omitempty"`
	KubeconfigFile string                  `yaml:"kubeconfig_file,omitempty"`
	Namespaces     *KubernetesSdNamespace  `yaml:"namespaces,omitempty"`
	Selectors      *[]KubernetesSdSelector `yaml:"selectors,omitempty"`
	AttachMetadata *AttachMetadata         `yaml:"attach_metadata,omitempty"`
}

type AttachMetadata struct {
	Node *bool `yaml:"node,omitempty"`
}

// KubernetesSdNamespace defines the kubernetes service discovery namespace entity.
type KubernetesSdNamespace struct {
	OwnNamespace *bool    `yaml:"own_namespace,omitempty"`
	Names        []string `yaml:"names,omitempty"`
}

// KubernetesSdSelector defines the kubernetes service discovery selector entity.
type KubernetesSdSelector struct {
	Role  string `yaml:"role,omitempty"`
	Label string `yaml:"label,omitempty"`
	Field string `yaml:"field,omitempty"`
}

// QueueConfig represents the remote-write queue config.
type QueueConfig struct {
	Capacity          int           `yaml:"capacity"`
	MaxShards         int           `yaml:"max_shards"`
	MinShards         int           `yaml:"min_shards"`
	MaxSamplesPerSend int           `yaml:"max_samples_per_send"`
	BatchSendDeadLine time.Duration `yaml:"batch_send_deadline"`
	MinBackoff        time.Duration `yaml:"min_backoff"`
	MaxBackoff        time.Duration `yaml:"max_backoff"`
	RetryOnHTTP429    bool          `yaml:"retry_on_http_429"`
}

// RemoteWrite represents a prometheus remote_write config.
type RemoteWrite struct {
	URL                 string          `yaml:"url"`
	RemoteTimeout       time.Duration   `yaml:"remote_timeout,omitempty"`
	Authorization       Authorization   `yaml:"authorization"`
	TLSConfig           *TLSConfig      `yaml:"tls_config,omitempty"`
	ProxyURL            string          `yaml:"proxy_url,omitempty"`
	QueueConfig         *QueueConfig    `yaml:"queue_config,omitempty"`
	WriteRelabelConfigs []RelabelConfig `yaml:"write_relabel_configs,omitempty"`
}
