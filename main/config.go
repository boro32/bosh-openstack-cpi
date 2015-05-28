package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bocaction "github.com/frodenas/bosh-openstack-cpi/action"
	bocconfig "github.com/frodenas/bosh-openstack-cpi/openstack/config"
)

type Config struct {
	OpenStack bocconfig.Config
	Actions   bocaction.ConcreteFactoryOptions
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config %s", path)
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	if err = config.Validate(); err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (c Config) Validate() error {
	if err := c.OpenStack.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating OpenStack configuration")
	}

	if err := c.Actions.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Actions configuration")
	}

	return nil
}
