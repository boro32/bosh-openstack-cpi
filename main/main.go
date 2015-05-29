package main

import (
	"flag"
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	bocaction "github.com/frodenas/bosh-openstack-cpi/action"
	bocdisp "github.com/frodenas/bosh-openstack-cpi/api/dispatcher"
	boctrans "github.com/frodenas/bosh-openstack-cpi/api/transport"

	"github.com/frodenas/bosh-openstack-cpi/openstack/client"
)

const mainLogTag = "main"

var (
	configFileOpt = flag.String("configFile", "", "Path to configuration file")
)

func main() {
	logger, fs, cmdRunner, uuidGen := basicDeps()

	defer logger.HandlePanic("Main")

	flag.Parse()

	config, err := NewConfigFromPath(*configFileOpt, fs)
	if err != nil {
		logger.Error(mainLogTag, "Loading config - %s", err.Error())
		os.Exit(1)
	}

	dispatcher, err := buildDispatcher(config, logger, fs, cmdRunner, uuidGen)
	if err != nil {
		logger.Error(mainLogTag, "Building Dispatcher - %s", err)
		os.Exit(1)
	}

	cli := boctrans.NewCLI(os.Stdin, os.Stdout, dispatcher, logger)

	if err = cli.ServeOnce(); err != nil {
		logger.Error(mainLogTag, "Serving once %s", err)
		os.Exit(1)
	}
}

func basicDeps() (boshlog.Logger, boshsys.FileSystem, boshsys.CmdRunner, boshuuid.Generator) {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr, os.Stderr)

	fs := boshsys.NewOsFileSystem(logger)

	cmdRunner := boshsys.NewExecCmdRunner(logger)

	uuidGen := boshuuid.NewGenerator()

	return logger, fs, cmdRunner, uuidGen
}

func buildDispatcher(
	config Config,
	logger boshlog.Logger,
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	uuidGen boshuuid.Generator,
) (bocdisp.Dispatcher, error) {
	openstackClient, err := client.NewOpenStackClient(config.OpenStack)
	if err != nil {
		return nil, err
	}

	actionFactory := bocaction.NewConcreteFactory(
		openstackClient,
		uuidGen,
		config.Actions,
		logger,
	)

	caller := bocdisp.NewJSONCaller()

	return bocdisp.NewJSON(actionFactory, caller, logger), nil
}
