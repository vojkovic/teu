package deployment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

func DeploymentComposeUp(composeFilePath string, projectName string) error {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{composeFilePath},
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithName(projectName),
	)
	if err != nil {
		return err
	}

	project, err := cli.ProjectFromOptions(ctx, options)
	if err != nil {
		return err
	}

	// Use the MarshalYAML method to get YAML representation
	projectYAML, err := project.MarshalYAML()
	if err != nil {
		return err
	}

	var dockerComposeBody string = string(projectYAML)

	p := createDockerProject(ctx, dockerComposeBody, projectName)

	srv, err := createDockerService()

	if err != nil {
		return fmt.Errorf("error create docker service: %v", err)
	}

	err = srv.Up(ctx, p, api.UpOptions{})
	if err != nil {
		return fmt.Errorf("error up: %v", err)
	}

	return nil
}

func DeploymentComposeDown(composeFilePath string, projectName string) error {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{composeFilePath},
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithName(projectName),
	)
	if err != nil {
		return err
	}

	project, err := cli.ProjectFromOptions(ctx, options)
	if err != nil {
		return err
	}

	// Use the MarshalYAML method to get YAML representation
	projectYAML, err := project.MarshalYAML()
	if err != nil {
		return err
	}

	var dockerComposeBody string = string(projectYAML)

	p := createDockerProject(ctx, dockerComposeBody, projectName)

	srv, err := createDockerService()

	if err != nil {
		return fmt.Errorf("error create docker service: %v", err)
	}

	err = srv.Down(ctx, p.Name, api.DownOptions{})
	if err != nil {
		return fmt.Errorf("error down: %v", err)
	}

	return nil
}

func createDockerProject(ctx context.Context, dockerComposeBody, projectName string) *types.Project {
	configDetails := types.ConfigDetails{
		// Fake path, doesn't need to exist.
		WorkingDir: "/in-memory/",
		ConfigFiles: []types.ConfigFile{
			{Filename: "docker-compose.yaml", Content: []byte(dockerComposeBody)},
		},
		Environment: nil,
	}

	p, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(projectName, true)
	})

	if err != nil {
		log.Fatalln("error load:", err)
	}
	addServiceLabels(p)
	return p
}

// createDockerService creates a docker service which can be
// used to interact with docker-compose.
func createDockerService() (api.Service, error) {
	var srv api.Service
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return srv, err
	}

	dockerContext := "default"

	//Magic line to fix error:
	//Failed to initialize: unable to resolve docker endpoint: no context store initialized
	myOpts := &flags.ClientOptions{Context: dockerContext, LogLevel: "error"}
	err = dockerCli.Initialize(myOpts)
	if err != nil {
		return srv, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, nil
}

/*
addServiceLabels adds the labels docker compose expects to exist on services.
This is required for future compose operations to work, such as finding
containers that are part of a service.
*/
func addServiceLabels(project *types.Project) {
	for i, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False", // default, will be overridden by `run` command
		}
		project.Services[i] = s
	}
}