package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Template represents a provider configuration template
type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	EnvVars     map[string]string `json:"envVars"`
}

// templateDir is the directory where custom templates are stored
var templateDir string

// builtInTemplates contains the built-in provider templates
var builtInTemplates = map[string]Template{
	"glm": {
		Name:        "glm",
		Description: "GLM (Zhipu AI) provider configuration",
		EnvVars: map[string]string{
			"ANTHROPIC_BASE_URL":             "https://open.bigmodel.cn/api/anthropic",
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":  "glm-4.5-air",
			"ANTHROPIC_DEFAULT_SONNET_MODEL": "glm-5-turbo",
			"ANTHROPIC_DEFAULT_OPUS_MODEL":   "glm-5.1",
		},
	},
	"deepseek": {
		Name:        "deepseek",
		Description: "DeepSeek provider configuration",
		EnvVars: map[string]string{
			"ANTHROPIC_BASE_URL":                       "https://api.deepseek.com/anthropic",
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":            "deepseek-v4-flash",
			"ANTHROPIC_DEFAULT_SONNET_MODEL":           "deepseek-v4-pro[1m]",
			"ANTHROPIC_DEFAULT_OPUS_MODEL":             "deepseek-v4-pro[1m]",
			"CLAUDE_CODE_SUBAGENT_MODEL":               "deepseek-v4-flash",
			"CLAUDE_CODE_EFFORT_LEVEL":                 "max",
			"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
		},
	},
	"mimo": {
		Name:        "mimo",
		Description: "Mimo provider configuration",
		EnvVars: map[string]string{
			"ANTHROPIC_BASE_URL":             "https://token-plan-cn.xiaomimimo.com/anthropic",
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":  "mimo-v2.5-pro",
			"ANTHROPIC_DEFAULT_SONNET_MODEL": "mimo-v2.5-pro",
			"ANTHROPIC_DEFAULT_OPUS_MODEL":   "mimo-v2.5-pro",
		},
	},
}

// initTemplateDir initializes the template directory
func initTemplateDir() error {
	if templateDir == "" {
		templateDir = filepath.Join(cfgDir, "templates")
	}
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}
	return nil
}

// getTemplate returns a template by name, checking built-in templates first
func getTemplate(name string) (*Template, error) {
	// Check built-in templates first
	if tmpl, ok := builtInTemplates[name]; ok {
		return &tmpl, nil
	}

	// Check custom templates
	if err := initTemplateDir(); err != nil {
		return nil, err
	}

	templatePath := filepath.Join(templateDir, name+".json")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("template '%s' not found", name)
	}

	var tmpl Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template '%s': %w", name, err)
	}

	return &tmpl, nil
}

// listTemplates returns all available templates (built-in + custom)
func listTemplates() ([]Template, error) {
	if err := initTemplateDir(); err != nil {
		return nil, err
	}

	var templates []Template

	// Add built-in templates
	for _, tmpl := range builtInTemplates {
		templates = append(templates, tmpl)
	}

	// Add custom templates
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		return templates, nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		templatePath := filepath.Join(templateDir, entry.Name())
		data, err := os.ReadFile(templatePath)
		if err != nil {
			continue
		}

		var tmpl Template
		if err := json.Unmarshal(data, &tmpl); err != nil {
			continue
		}

		templates = append(templates, tmpl)
	}

	return templates, nil
}

// saveCustomTemplate saves a custom template
func saveCustomTemplate(tmpl Template) error {
	if err := initTemplateDir(); err != nil {
		return err
	}

	templatePath := filepath.Join(templateDir, tmpl.Name+".json")
	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	if err := os.WriteFile(templatePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	return nil
}

// deleteCustomTemplate deletes a custom template
func deleteCustomTemplate(name string) error {
	if err := initTemplateDir(); err != nil {
		return err
	}

	// Prevent deleting built-in templates
	if _, ok := builtInTemplates[name]; ok {
		return fmt.Errorf("cannot delete built-in template '%s'", name)
	}

	templatePath := filepath.Join(templateDir, name+".json")
	if err := os.Remove(templatePath); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}
