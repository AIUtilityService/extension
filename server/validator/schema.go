package validator

import "fmt"

type ResourceGraphDefinition struct {
	Version     string       `yaml:"version"`
	Resources   []Resource   `yaml:"resources"`
	Connections []Connection `yaml:"connections"`
}

type Resource struct {
	Name       string                 `yaml:"name"`
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties"`
}

type Connection struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	Type string `yaml:"type"`
}

func ValidateResourceGraph(data map[string]interface{}) []error {
	var errors []error

	// Validate version
	if version, ok := data["version"].(string); !ok || version == "" {
		errors = append(errors, fmt.Errorf("invalid or missing version"))
	}

	// Validate resources
	if resources, ok := data["resources"].([]interface{}); ok {
		for _, res := range resources {
			if resource, ok := res.(map[string]interface{}); ok {
				if err := validateResource(resource); err != nil {
					errors = append(errors, err)
				}
			}
		}
	}

	return errors
}

func validateResource(resource map[string]interface{}) error {
	if _, ok := resource["name"].(string); !ok {
		return fmt.Errorf("resource missing name")
	}

	resourceType, ok := resource["type"].(string)
	if !ok {
		return fmt.Errorf("resource missing type")
	}

	// Add validation for resource type
	validTypes := []string{"compute", "storage", "network", "database"}
	isValidType := false
	for _, validType := range validTypes {
		if resourceType == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid resource type: %s", resourceType)
	}

	return nil
}
