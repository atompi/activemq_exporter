/*
Copyright © 2024 Atom Pi <coder.atompi@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

type JolokiaOptions struct {
	URL                string `yaml:"url"`
	BasicAuth          bool   `yaml:"basic_auth"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Cert               string `yaml:"cert"`
	Key                string `yaml:"key"`
	CA                 string `yaml:"ca"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

type BrokerOptions struct {
	Name string `yaml:"name"`
}

type Options struct {
	Jolokia JolokiaOptions `yaml:"jolokia"`
	Broker  BrokerOptions  `yaml:"broker"`
}

var Opts = &Options{}
