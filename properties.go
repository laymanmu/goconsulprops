package goconsulprops

import (
	"log"
	"path/filepath"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// VersionedValue ...
type VersionedValue struct {
	Value   string
	Version uint64
}

// Properties holds key/value strings from consul for a given prefix.
type Properties struct {
	name        string
	prefix      string
	consul      *consul.KV
	props       map[string]VersionedValue
	refreshedAt time.Time
}

// NewProperties creates a Properties struct based on a given prefix & refreshes its key/values.
//func NewProperties(prefix string, kv *consul.KV) *Properties {
func NewProperties(consulAddress string, consulPrefix string) *Properties {
	// connection:
	config := consul.DefaultConfig()
	config.Address = consulAddress
	client, err := consul.NewClient(config)
	if err != nil {
		panic(err)
	}

	// strip trailing slash to get this base property name:
	if consulPrefix[len(consulPrefix)-1:] == "/" {
		consulPrefix = consulPrefix[:len(consulPrefix)-1]
	}
	name := filepath.Base(consulPrefix)

	// add trailing slash back:
	consulPrefix = consulPrefix + "/"

	props := make(map[string]VersionedValue)
	p := &Properties{name: name, prefix: consulPrefix, consul: client.KV(), props: props}
	p.Refresh()
	return p
}

// Refresh will set/update the stored key/value strings from a consul server.
func (p *Properties) Refresh() {
	pairs, _, err := p.consul.List(p.prefix, nil)
	if err != nil {
		panic(err)
	}
	for _, pair := range pairs {
		key := pair.Key[len(p.prefix):]
		if len(key) == 0 {
			continue
		}
		version := pair.ModifyIndex
		value := string(pair.Value)
		p.props[key] = VersionedValue{Value: value, Version: version}
		log.Printf("[goconsulprops] set %v.%v: %v (version: %v)\n", p.name, key, value, version)
	}
	p.refreshedAt = time.Now()
}

// GetValue returns a property value.
func (p *Properties) GetValue(key string) string {
	return p.props[key].Value
}

// GetVersion returns a property value.
func (p *Properties) GetVersion(key string) uint64 {
	return p.props[key].Version
}

// RefreshedAt returns the time when key/values were last refreshed.
func (p *Properties) RefreshedAt() time.Time {
	return p.refreshedAt
}
