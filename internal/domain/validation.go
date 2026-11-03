package domain

import "github.com/google/uuid"

func ValidateCategory(category *Category) error {
	fieldsErr := make(map[string]string)

	if category.Name == "" {
		fieldsErr["name"] = "name must be not empty"
	}

	if len(fieldsErr) > 0 {
		return &ValidationErr{Fields: fieldsErr}
	}
	return nil
}

func ValidationAssistant(assistant *Assistant) error {
	fieldsErr := make(map[string]string)

	if assistant.Name == "" {
		fieldsErr["name"] = "name must be not empty"
	}

	if assistant.CategoryID == uuid.Nil {
		fieldsErr["categoryId"] = "category id must be not empty"
	}

	if assistant.Model == "" {
		fieldsErr["model"] = "model must be not empty"
	}

	if assistant.Description == "" {
		fieldsErr["description"] = "description must be not empty"
	}

	if assistant.SystemPrompt == "" {
		fieldsErr["system_prompt"] = "system prompt must be not empty"
	}

	if len(fieldsErr) > 0 {
		return &ValidationErr{Fields: fieldsErr}
	}
	return nil
}
