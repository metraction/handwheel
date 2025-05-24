package integrations

import (
	"log"
	"regexp"

	"github.com/Tiktai/handler/model"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type CraneIntegration struct {
	config     *model.Config
	remoteOpts []remote.Option
}

func NewCraneIntegration(cfg *model.Config) *CraneIntegration {
	var opts []remote.Option
	// Add transport option
	opts = append(opts, remote.WithTransport(NewHttpTransport(cfg)))
	return &CraneIntegration{config: cfg, remoteOpts: opts}
}

// getAuthForRegistry returns the auth for a given registry, or nil if not found
func getAuthForRegistry(cfg *model.Config, registry string) authn.Authenticator {
	for _, reg := range cfg.Crane.Registries {
		if reg.Registry == registry {
			return &authn.Basic{
				Username: reg.Username,
				Password: reg.Password,
			}
		}
	}
	return authn.Anonymous
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

	// Get registry from ref.Context().RegistryStr()
	auth := getAuthForRegistry(ci.config, ref.Context().RegistryStr())
	log.Printf("Using auth for registry %s: %v", ref.Context().RegistryStr(), auth)
	remoteImg, err := remote.Image(ref, append(ci.remoteOpts, remote.WithAuth(auth))...)
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
