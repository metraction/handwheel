package integrations

import (
	"log"
	"regexp"

	"github.com/Tiktai/handler/model"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/authn"
)

type CraneIntegration struct {
	config model.CraneConfig
	remoteOpts []remote.Option
}

func NewCraneIntegration(cfg model.CraneConfig) *CraneIntegration {
	var opts []remote.Option
	if cfg.RegistryUsername != "" && cfg.RegistryPassword != "" {
		auth := &authn.Basic{
			Username: cfg.RegistryUsername,
			Password: cfg.RegistryPassword,
		}
		opts = append(opts, remote.WithAuth(auth))
	}
	return &CraneIntegration{config: cfg, remoteOpts: opts}
} 

func (ci *CraneIntegration) WhiteListImages() func(elem model.ImageMetric) bool {
	return func(elem model.ImageMetric) bool {
		// Check whitelist
		matched := false
		for _, pattern := range ci.config.Images_whitelist {
			re, err := regexp.Compile(pattern)
			if err != nil {
				log.Printf("invalid whitelist regexp: %s: %v", pattern, err)
				continue
			}
			if re.MatchString(elem.Image_spec) {
				matched = true
				break
			}
		}
		return matched
	}
}

func (ci *CraneIntegration) CraneRetrieveLabels(elem model.ImageMetric) model.Image {
	img := model.Image{Image_spec: elem.Image_spec, Labels: map[string]string{}}
	if elem.Image_spec == "" {
		return img
	}

	ref, err := name.ParseReference(elem.Image_spec)
	if err != nil {
		log.Printf("failed to parse image reference %s: %v", elem.Image_spec, err)
		return img
	}

	remoteImg, err := remote.Image(ref, ci.remoteOpts...)
	if err != nil {
		log.Printf("failed to fetch image %s: %v", elem.Image_spec, err)
		return img
	}
	cfg, err := remoteImg.ConfigFile()
	if err != nil {
		log.Printf("failed to get config for image %s: %v", elem.Image_spec, err)
		return img
	}
	img.Labels = cfg.Config.Labels
	return img
}
