package integrations

import (
	"log"
	"regexp"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/metraction/handwheel/model"
)

type CraneIntegration struct {
	config     *model.Config
	remoteOpts []remote.Option
}

type MultiRegistryKeychain struct {
	cfg *model.Config
}

func (kc *MultiRegistryKeychain) Resolve(resource authn.Resource) (authn.Authenticator, error) {
	for _, reg := range kc.cfg.Crane.Registries {
		if reg.Registry == resource.RegistryStr() {
			return &authn.Basic{
				Username: reg.Username,
				Password: reg.Password,
			}, nil
		}
	}
	return authn.Anonymous, nil
}

func NewCraneIntegration(cfg *model.Config) *CraneIntegration {
	var opts []remote.Option
	// Add transport option
	opts = append(opts, remote.WithTransport(NewHttpTransport(cfg)))
	opts = append(opts, remote.WithAuthFromKeychain(&MultiRegistryKeychain{cfg}))
	return &CraneIntegration{config: cfg, remoteOpts: opts}
}

func (ci *CraneIntegration) WhiteListImages() func(elem model.ImageMetric) bool {
	// Collect all patterns from devlake.projects.images
	var patterns []string
	for _, project := range ci.config.DevLake.Projects {
		patterns = append(patterns, project.Images...)
	}
	return func(elem model.ImageMetric) bool {
		matched := false
		for _, pattern := range patterns {
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
