package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
)

const defaultRegion = "US"

var regionRe = regexp.MustCompile(`^[A-Z]{2}$`)

// configuration mirrors settings_schema in plugin.json.
type configuration struct {
	Region string
}

func (c configuration) region() (string, error) {
	r := strings.ToUpper(strings.TrimSpace(c.Region))
	if r == "" {
		return defaultRegion, nil
	}
	if !regionRe.MatchString(r) {
		return "", fmt.Errorf("Region must be a two-letter country code such as ES or US, got %q", c.Region)
	}
	return r, nil
}

type settings struct {
	regionCode atomic.Value
}

func (s *settings) region() string {
	if r, ok := s.regionCode.Load().(string); ok {
		return r
	}
	return defaultRegion
}

func (p *Plugin) OnConfigurationChange() error {
	var c configuration
	if err := p.API.LoadPluginConfiguration(&c); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	r, err := c.region()
	if err != nil {
		return err
	}
	p.settings.regionCode.Store(r)
	return nil
}
