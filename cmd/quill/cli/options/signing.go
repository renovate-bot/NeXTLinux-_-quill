package options

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/nextlinux/quill/internal/log"
)

var _ Interface = &Signing{}

type Signing struct {
	// bound options
	Identity             string `yaml:"identity" json:"identity" mapstructure:"identity"`
	P12                  string `yaml:"p12" json:"p12" mapstructure:"p12"`
	TimestampServer      string `yaml:"timestamp-server" json:"timestamp-server" mapstructure:"timestamp-server"`
	AdHoc                bool   `yaml:"ad-hoc" json:"ad-hoc" mapstructure:"ad-hoc"`
	FailWithoutFullChain bool   `yaml:"fail-without-full-chain" json:"fail-without-full-chain" mapstructure:"fail-without-full-chain"`

	// unbound options
	Password string `yaml:"password" json:"password" mapstructure:"password"`
}

func DefaultSigning() Signing {
	return Signing{
		TimestampServer:      "http://timestamp.apple.com/ts01",
		FailWithoutFullChain: true,
	}
}

func (o *Signing) Redact() {
	log.Redact(o.Password)
	redactNonFileOrEnvHint(o.P12)
}

func (o *Signing) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(
		&o.Identity,
		"identity", "", o.Identity,
		"identifier to encode into the code directory of the code signing super block (default is derived from the name of the binary being solved)",
	)

	flags.StringVarP(
		&o.P12,
		"p12", "", o.P12,
		"path to a PKCS12 file containing the private key, (leaf) signing certificate, remaining certificate chain. This can also be the base64-encoded contents of the p12 file, or 'env:ENV_VAR_NAME' to read the p12 from a different environment variable",
	)

	flags.StringVarP(
		&o.TimestampServer,
		"timestamp-server", "", o.TimestampServer,
		"URL to a timestamp server to use for timestamping the signature",
	)

	flags.BoolVarP(
		&o.AdHoc,
		"ad-hoc", "", o.AdHoc,
		"perform ad-hoc signing. No cryptographic signature is included and --p12 key and certificate input are not needed. Do NOT use this option for production builds.",
	)
}

func (o *Signing) BindFlags(flags *pflag.FlagSet, v *viper.Viper) error {
	if err := Bind(v, "sign.override-identity", flags.Lookup("identity")); err != nil {
		return err
	}
	if err := Bind(v, "sign.p12", flags.Lookup("p12")); err != nil {
		return err
	}
	if err := Bind(v, "sign.timestamp-server", flags.Lookup("timestamp-server")); err != nil {
		return err
	}
	if err := Bind(v, "sign.ad-hoc", flags.Lookup("ad-hoc")); err != nil {
		return err
	}

	// set default values for non-bound struct items
	v.SetDefault("sign.password", o.Password)
	v.SetDefault("sign.fail-without-full-chain", o.FailWithoutFullChain)

	return nil
}
