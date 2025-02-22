package signature

import (
	"fmt"
	"strings"

	"github.com/toothbrush/go-pipeline"
)

var _ SignedFielder = (*CommandStepWithInvariants)(nil)

// CommandStepWithInvariants is a CommandStep with PipelineInvariants.
type CommandStepWithInvariants struct {
	pipeline.CommandStep
	RepositoryURL string
}

// SignedFields returns the default fields for signing.
func (c *CommandStepWithInvariants) SignedFields() (map[string]any, error) {
	return map[string]any{
		"command":        c.Command,
		"env":            EmptyToNilMap(c.Env),
		"plugins":        EmptyToNilSlice(c.Plugins),
		"matrix":         EmptyToNilPtr(c.Matrix),
		"repository_url": c.RepositoryURL,
	}, nil
}

// ValuesForFields returns the contents of fields to sign.
func (c *CommandStepWithInvariants) ValuesForFields(fields []string) (map[string]any, error) {
	// Make a set of required fields. As fields is processed, mark them off by
	// deleting them.
	required := map[string]struct{}{
		"command":        {},
		"env":            {},
		"plugins":        {},
		"matrix":         {},
		"repository_url": {},
	}

	out := make(map[string]any, len(fields))
	for _, f := range fields {
		delete(required, f)

		switch f {
		case "command":
			out["command"] = c.Command

		case "env":
			out["env"] = EmptyToNilMap(c.Env)

		case "plugins":
			out["plugins"] = EmptyToNilSlice(c.Plugins)

		case "matrix":
			out["matrix"] = EmptyToNilPtr(c.Matrix)

		case "repository_url":
			out["repository_url"] = c.RepositoryURL

		default:
			// All env:: values come from outside the step.
			if strings.HasPrefix(f, EnvNamespacePrefix) {
				break
			}

			return nil, fmt.Errorf("unknown or unsupported field for signing %q", f)
		}
	}

	if len(required) > 0 {
		missing := make([]string, 0, len(required))
		for k := range required {
			missing = append(missing, k)
		}
		return nil, fmt.Errorf("one or more required fields are not present: %v", missing)
	}
	return out, nil
}
