package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/driftcheck/internal/compose"
	"github.com/yourorg/driftcheck/internal/docker"
	"github.com/yourorg/driftcheck/internal/drift"
)

func main() {
	composeFile := flag.String("file", "docker-compose.yml", "path to docker-compose file")
	service := flag.String("service", "", "service name to check (empty = all services)")
	jsonOutput := flag.Bool("json", false, "output results as JSON")
	flag.Parse()

	project, err := compose.ParseFile(*composeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing compose file: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	inspector, err := docker.NewInspector(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to Docker: %v\n", err)
		os.Exit(1)
	}
	defer inspector.Close()

	detector := drift.NewDetector()
	reporter := drift.NewReporter(os.Stdout, *jsonOutput)

	services := project.Services
	if *service != "" {
		svc, ok := project.GetService(*service)
		if !ok {
			fmt.Fprintf(os.Stderr, "service %q not found in compose file\n", *service)
			os.Exit(1)
		}
		services = map[string]compose.Service{*service: svc}
	}

	var results []drift.Result
	for name, svc := range services {
		container, err := inspector.InspectByServiceName(ctx, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not inspect service %q: %v\n", name, err)
			continue
		}
		result := detector.Check(name, svc, container)
		results = append(results, result)
	}

	if err := reporter.Report(results); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	for _, r := range results {
		if len(r.Diffs) > 0 {
			os.Exit(2)
		}
	}
}
