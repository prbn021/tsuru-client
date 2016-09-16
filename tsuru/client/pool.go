// Copyright 2016 tsuru-client authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/tsuru/tsuru/cmd"
)

type PoolList struct{}

type Pool struct {
	Name        string
	Teams       []string
	Public      bool
	Default     bool
	Provisioner string
}

func (p *Pool) Kind() string {
	if p.Public {
		return "public"
	}
	if p.Default {
		return "default"
	}
	return ""
}

func (p *Pool) GetProvisioner() string {
	if p.Provisioner == "" {
		return "default"
	}
	return p.Provisioner
}

type poolEntriesList []Pool

func (l poolEntriesList) Len() int      { return len(l) }
func (l poolEntriesList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l poolEntriesList) Less(i, j int) bool {
	cmp := strings.Compare(l[i].Kind(), l[j].Kind())
	if cmp == 0 {
		return l[i].Name < l[j].Name
	}
	return cmp < 0
}

func (PoolList) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL("/pools")
	if err != nil {
		return err
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var pools []Pool
	err = json.NewDecoder(resp.Body).Decode(&pools)
	if err != nil {
		return err
	}
	sort.Sort(poolEntriesList(pools))
	t := cmd.Table{Headers: cmd.Row([]string{"Pool", "Kind", "Provisioner", "Teams"})}
	for _, pool := range pools {
		t.AddRow(cmd.Row([]string{pool.Name, pool.Kind(), pool.GetProvisioner(), strings.Join(pool.Teams, ", ")}))
	}
	context.Stdout.Write(t.Bytes())
	return nil
}

func (PoolList) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "pool-list",
		Usage:   "pool-list",
		Desc:    "List all pools available for deploy.",
		MinArgs: 0,
	}
}
