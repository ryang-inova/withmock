// Copyright 2013 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lib

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

func readPackages(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pkgs := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		pkgs = append(pkgs, line)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return pkgs, nil
}

type MockConfig struct {
	// Local configuration
	MockPrototypes   bool // Mock prototypes (i.e. functions without bodies)
	IgnoreInits      bool // Don't call the original init functions
	MatchOSArch      bool // only use files for GOOS & GOARCH
	IgnoreNonGoFiles bool // Don't copy non-go files into the mocked package

	// File based configuration
	MOCK      string `yaml:"MOCK"`
	EXPECT    string `yaml:"EXPECT"`
	ObjEXPECT string `yaml:"obj.EXPECT"`
}

type Config struct {
	Mocks map[string]*MockConfig
}

func (c *Config) Mock(path string) *MockConfig {
	m := &MockConfig{
		MOCK:      "MOCK",
		EXPECT:    "EXPECT",
		ObjEXPECT: "EXPECT",
	}

	dc, found := c.Mocks["DEFAULT"]
	if !found {
		dc = &MockConfig{}
	}

	mc, found := c.Mocks[path]
	if !found {
		mc = &MockConfig{}
	}

	switch {
	case mc.MOCK != "":
		m.MOCK = mc.MOCK
	case dc.MOCK != "":
		m.MOCK = dc.MOCK
	}

	switch {
	case mc.EXPECT != "":
		m.EXPECT = mc.EXPECT
	case dc.EXPECT != "":
		m.EXPECT = dc.EXPECT
	}

	switch {
	case mc.ObjEXPECT != "":
		m.ObjEXPECT = mc.ObjEXPECT
	case dc.ObjEXPECT != "":
		m.ObjEXPECT = dc.ObjEXPECT
	}

	return m
}

func ReadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = goyaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
