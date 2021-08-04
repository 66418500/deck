package file

import (
	"encoding/json"
	"fmt"
	"strings"

	ghodss "github.com/ghodss/yaml"
	"github.com/kong/deck/utils"
	"github.com/xeipuuv/gojsonschema"
)

func validate(content []byte) error {
	var c map[string]interface{}
	err := ghodss.Unmarshal(content, &c)
	if err != nil {
		return fmt.Errorf("unmarshaling file content: %w", err)
	}
	c = ensureJSON(c)
	schemaLoader := gojsonschema.NewStringLoader(kongJSONSchema)
	documentLoader := gojsonschema.NewGoLoader(c)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if result.Valid() {
		return nil
	}
	var errs utils.ErrArray
	for _, desc := range result.Errors() {
		err := fmt.Errorf(desc.String())
		errs.Errors = append(errs.Errors, err)
	}

	var kongdefaults *KongDefaults
	if err := json.Unmarshal(content, kongdefaults); err != nil {
		return fmt.Errorf("unmarshaling file into KongDefaults: %w", err)
	}
	if kongdefaults != nil {
		if err := CheckDefaults(*kongdefaults); err != nil {
			return fmt.Errorf("default values are not allowed for these fields: %w", err)
		}
	}
	return errs
}

func Check(item string, t *string) string {
	if t != nil || len(*t) > 0 {
		return item
	}
	return ""
}

func CheckDefaults(kd KongDefaults) error {
	var invalid []string
	if ret := Check("Service.ID", kd.Service.ID); len(ret) > 0 {
		invalid = append(invalid, "Service.ID")
	}

	if ret := Check("Service.Host", kd.Service.Host); len(ret) > 0 {
		invalid = append(invalid, "Service.Host")
	}

	if ret := Check("Service.Name", kd.Service.Name); len(ret) > 0 {
		invalid = append(invalid, "Service.Name")
	}

	if kd.Service.Port != nil {
		invalid = append(invalid, "Service.Port")
	}

	if ret := Check("Route.Name", kd.Route.Name); len(ret) > 0 {
		invalid = append(invalid, "Route.Name")
	}

	if ret := Check("Route.ID", kd.Route.ID); len(ret) > 0 {
		invalid = append(invalid, "Route.ID")
	}

	if ret := Check("Target.Target", kd.Target.Target); len(ret) > 0 {
		invalid = append(invalid, "Target.Target")
	}

	if ret := Check("Target.ID", kd.Target.ID); len(ret) > 0 {
		invalid = append(invalid, "Target.ID")
	}

	if ret := Check("Upstream.Name", kd.Upstream.Name); len(ret) > 0 {
		invalid = append(invalid, "Upstream.Name")
	}

	if ret := Check("Upstream.ID", kd.Upstream.ID); len(ret) > 0 {
		invalid = append(invalid, "Upstream.ID")
	}

	if len(invalid) > 0 {
		return fmt.Errorf("unacceptable fields in defaults: %s", strings.Join(invalid, ", "))
	}
	return nil
}
