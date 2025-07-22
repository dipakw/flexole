package app

import (
	"embed"
	"fmt"
	"os"
)

func generateConfig(samples *embed.FS) error {
	cli := NewCli(map[string]string{
		"client-config": "client.yml",
		"server-config": "server.yml",
	})

	passed := false

	if cli.Get("client-config").Passed {
		passed = true

		if err := generateClientConfig(cli.Get("client-config").Value(), samples); err != nil {
			return err
		}
	}

	if cli.Get("server-config").Passed {
		passed = true

		if err := generateServerConfig(cli.Get("server-config").Value(), samples); err != nil {
			return err
		}
	}

	if !passed {
		return fmt.Errorf(`Either --client-config or --server-config must be provided`)
	}

	return nil
}

func generateClientConfig(name string, samples *embed.FS) error {
	sampleBytes, err := samples.ReadFile("client.yml")

	if err != nil {
		return err
	}

	return os.WriteFile(name, sampleBytes, 0644)
}

func generateServerConfig(name string, samples *embed.FS) error {
	sampleBytes, err := samples.ReadFile("server.yml")

	if err != nil {
		return err
	}

	return os.WriteFile(name, sampleBytes, 0644)
}
