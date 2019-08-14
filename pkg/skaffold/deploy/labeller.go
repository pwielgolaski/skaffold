/*
Copyright 2019 The Skaffold Authors

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

package deploy

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/version"
)

const (
	K8sManagedByLabelKey = "app.kubernetes.io/managed-by"
	RunIDLabel           = "skaffold.dev/run-id"
	unknownVersion       = "unknown"
	empty                = ""
)

var (
	runID     string
	runIDOnce sync.Once
)

// DefaultLabeller adds K8s style managed-by label and a run-specific UUID label
type DefaultLabeller struct {
	version string
	runID   string
}

func NewLabeller(verStr string) *DefaultLabeller {
	runIDOnce.Do(func() {
		runID = uuid.New().String()
	})
	if verStr == empty {
		verStr = version.Get().Version
	}
	if verStr == empty {
		verStr = unknownVersion
	}
	return &DefaultLabeller{
		version: verStr,
		runID:   runID,
	}
}

func (d *DefaultLabeller) Labels() map[string]string {
	return map[string]string{
		K8sManagedByLabelKey: d.skaffoldVersion(),
		RunIDLabel:           d.runID,
	}
}

func (d *DefaultLabeller) RunIDKeyValueString() string {
	return fmt.Sprintf("%s=%s", RunIDLabel, d.runID)
}

func (d *DefaultLabeller) K8sManagedByLabelKeyValueString() string {
	return fmt.Sprintf("%s=%s", K8sManagedByLabelKey, d.skaffoldVersion())
}

func (d *DefaultLabeller) skaffoldVersion() string {
	return fmt.Sprintf("skaffold-%s", d.version)
}
