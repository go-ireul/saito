package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Package package
type Package struct {
	Name        string `json:"name"`
	CloneURL    string `json:"clone_url"`
	ImportURL   string `json:"-"`
	Description string `json:"description"`
}

// PackageManager package manager
type PackageManager struct {
	Option PackageManagerOption
	pkgs   []Package
	apiURL string
	mutex  *sync.RWMutex
}

// PackageManagerOption option for package manager
type PackageManagerOption struct {
	Domain       string
	Token        string
	Organization string
}

// NewPackageManager create a new package manager
func NewPackageManager(option PackageManagerOption) *PackageManager {
	m := &PackageManager{
		Option: option,
		pkgs:   make([]Package, 0),
		mutex:  &sync.RWMutex{},
	}
	m.apiURL = fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100", option.Organization)
	return m
}

// StartTicking start update ticking
func (m *PackageManager) StartTicking() {
	tchan := time.Tick(time.Second * 60)
	go func() {
		m.updatePackages()
		for {
			<-tchan
			m.updatePackages()
		}
	}()
}

// List list all packages
func (m *PackageManager) List() []Package {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.pkgs
}

// Get get package by name
func (m *PackageManager) Get(name string) (pkg Package) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, p := range m.pkgs {
		if p.Name == name {
			pkg = p
		}
	}
	return
}

func (m *PackageManager) updatePackages() (err error) {
	// get api url
	var resp *http.Response
	if resp, err = http.Get(m.apiURL); err != nil {
		log.Println("failed to get api url:", m.apiURL, err)
		return
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("failed to read api url:", m.apiURL, err)
		return
	}
	// unmarshal json
	var pkgs []Package
	if err = json.Unmarshal(body, &pkgs); err != nil {
		log.Println("failed to unmarshal json:", m.apiURL, err)
		return
	}
	// update m.pkgs
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i, p := range pkgs {
		pkgs[i].ImportURL = fmt.Sprintf("%s/%s", m.Option.Domain, p.Name)
	}
	m.pkgs = pkgs
	return
}
