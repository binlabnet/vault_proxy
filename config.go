package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	VaultConfig         *VaultConfig           `yaml:"vaultConfig"`
	CookieEncryptionKey string                 `yaml:"cookieEncryptionKey"`
	CookieName          string                 `yaml:"cookieName"`
	HeaderName          string                 `yaml:"headerName"`
	PublicURLRaw        string                 `yaml:"publicURL"`
	UpstreamURLRaw      string                 `yaml:"upstreamURL"`
	Meta                map[string]interface{} `yaml:"meta"`
	AccessList          []*AccesItem           `yaml:"accessList"`
	routeRegExpMap      map[string]*regexp.Regexp
	publicURL           *url.URL
	upstreamURL         *url.URL
}

type AccesItem struct {
	Name      string   `yaml:"name"`
	Path      string   `yaml:"path"`
	Policies  []string `yaml:"policies"`
	Methods   []string `yaml:"methods"`
	methodMap map[string]bool
	policyMap map[string]bool
	re        *regexp.Regexp
}

func LoadConfig(filename string) (*Config, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to load configuration file. %v", err)
	}
	c := &Config{}
	err = yaml.Unmarshal(out, &c)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse configuration file. %v", err)
	}
	if err := c.Parse(); err != nil {
		return nil, err
	}
	return c, nil
}

func (a *AccesItem) Parse() error {
	a.policyMap = make(map[string]bool, len(a.Policies))
	for _, policy := range a.Policies {
		a.policyMap[policy] = true
	}
	a.methodMap = make(map[string]bool, len(a.Methods))
	for _, method := range a.Methods {
		a.methodMap[strings.ToUpper(method)] = true
	}
	re, err := regexp.Compile(a.Path)
	if err != nil {
		return fmt.Errorf("Unable to parse '%s' as regular expression. %v", a.Path, err)
	}
	a.re = re
	return nil
}

func (c *Config) Parse() error {
	publicURL, err := url.Parse(c.PublicURLRaw)
	if err != nil {
		return fmt.Errorf("Unable to parse PublicURL. %v", err)
	}
	c.publicURL = publicURL
	upstreamURL, err := url.Parse(c.UpstreamURLRaw)
	if err != nil {
		return fmt.Errorf("Unable to parse UpstreamURL. %v", err)
	}
	c.upstreamURL = upstreamURL
	c.routeRegExpMap = make(map[string]*regexp.Regexp, len(c.AccessList))
	for _, item := range c.AccessList {
		err := item.Parse()
		if err != nil {
			return err
		}
	}
	if c.VaultConfig == nil {
		return errors.New("Vault configuration must be specified")
	}
	if err := c.VaultConfig.Parse(); err != nil {
		return err
	}
	return nil
}