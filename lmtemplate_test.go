package modprompt

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestInitTemplate tests the InitTemplate function.
func TestInitTemplate(t *testing.T) {
	// Initialize the Templates map
	err := json.Unmarshal([]byte(templates), &Templates)
	if err != nil {
		t.Fatalf("Error unmarshalling templates: %v", err)
	}

	// Test with a valid template name
	tpl, err := InitTemplate("chatml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if tpl.ID != "chatml" {
		t.Errorf("Expected template ID to be 'chatml', got: %s", tpl.ID)
	}

	// Test with an invalid template name
	_, err = InitTemplate("nonexistent")
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	expectedErr := "No template named nonexistent"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message to be '%s', got: %s", expectedErr, err.Error())
	}
}

// TestPromptTemplate_Render tests the Render method of PromptTemplate.
func TestPromptTemplate_Render(t *testing.T) {
	// Initialize the Templates map
	err := json.Unmarshal([]byte(templates), &Templates)
	if err != nil {
		t.Fatalf("Error unmarshalling templates: %v", err)
	}

	// Get a template
	tpl, err := InitTemplate("chatml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Create a PromptTemplate instance
	promptTpl := &PromptTemplate{
		ID:         tpl.ID,
		Name:       tpl.Name,
		User:       tpl.User,
		Assistant:  tpl.Assistant,
		System:     tpl.System,
		Shots:      tpl.Shots,
		Stop:       tpl.Stop,
		Linebreaks: tpl.Linebreaks,
		AfterShot:  tpl.AfterShot,
		Prefix:     tpl.Prefix,
	}

	// Render the template
	rendered := promptTpl.Render(true)
	expected := `<|im_start|>system
{system}<|im_end|>
<|im_start|>user
{prompt}<|im_end|>
<|im_start|>assistant`
	expected = expected + "\n"
	//fmt.Println(rendered)
	if rendered != expected {
		t.Errorf("Expected template rendered to:\n\n|%s| \n\n Got:\n\n|%s|", expected, rendered)
	}

	// Prompt
	prompt := "List the planets of the solar system"
	rendered = promptTpl.Prompt(prompt)
	expected = `<|im_start|>system
{system}<|im_end|>
<|im_start|>user
List the planets of the solar system<|im_end|>
<|im_start|>assistant`
	expected = expected + "\n"
	if rendered != expected {
		t.Errorf("Expected template rendered to:\n\n|%s| \n\n Got:\n\n|%s|", expected, rendered)
	}
}

// TestPromptTemplate_Prompt test options
func TestPromptTemplate_Options(t *testing.T) {
	// Get a template
	tpl, err := InitTemplate("chatml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Create a PromptTemplate instance
	promptTpl := &PromptTemplate{
		ID:         tpl.ID,
		Name:       tpl.Name,
		User:       tpl.User,
		Assistant:  tpl.Assistant,
		System:     tpl.System,
		Shots:      tpl.Shots,
		Stop:       tpl.Stop,
		Linebreaks: tpl.Linebreaks,
		AfterShot:  tpl.AfterShot,
		Prefix:     tpl.Prefix,
	}

	system := "You are a helpful assistant"
	promptTpl.ReplaceSystem(system)
	// Render the template
	rendered := promptTpl.Render(true)
	expected := `<|im_start|>system
You are a helpful assistant<|im_end|>
<|im_start|>user
{prompt}<|im_end|>
<|im_start|>assistant`
	expected = expected + "\n"
	//fmt.Println(rendered)
	if rendered != expected {
		t.Errorf("Expected template rendered to:\n\n|%s| \n\n Got:\n\n|%s|", expected, rendered)
	}
}

// TestPromptTemplate_AddShot tests the AddShot method of PromptTemplate.
func TestPromptTemplate_AddShot(t *testing.T) {
	// Initialize the Templates map
	err := json.Unmarshal([]byte(templates), &Templates)
	if err != nil {
		t.Fatalf("Error unmarshalling templates: %v", err)
	}

	// Get a template
	tpl, err := InitTemplate("chatml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Create a PromptTemplate instance
	promptTpl := &PromptTemplate{
		ID:         tpl.ID,
		Name:       tpl.Name,
		User:       tpl.User,
		Assistant:  tpl.Assistant,
		System:     tpl.System,
		Shots:      tpl.Shots,
		Stop:       tpl.Stop,
		Linebreaks: tpl.Linebreaks,
		AfterShot:  tpl.AfterShot,
		Prefix:     tpl.Prefix,
	}

	// Add a shot
	user := "What is the capital of France?"
	assistant := "The capital of France is Paris."
	promptTpl.AddShot(user, assistant)

	// Check if the shot was added correctly
	if len(promptTpl.Shots) != 1 {
		t.Errorf("Expected 1 shot, got %d", len(promptTpl.Shots))
	}
	if promptTpl.Shots[0].User != user {
		t.Errorf("Expected user to be '%s', got '%s'", user, promptTpl.Shots[0].User)
	}
	if promptTpl.Shots[0].Assistant != assistant {
		t.Errorf("Expected assistant to be '%s', got '%s'", assistant, promptTpl.Shots[0].Assistant)
	}
	expected := `<|im_start|>system
{system}<|im_end|>
<|im_start|>user
What is the capital of France?<|im_end|>
<|im_start|>assistant
The capital of France is Paris. <|im_end|>
<|im_start|>user
{prompt}<|im_end|>
<|im_start|>assistant
`
	rendered := promptTpl.Render(true)
	if rendered != expected {
		t.Errorf("Expected template rendered to:\n\n|%s| \n\n Got:\n\n|%s|", expected, rendered)
	}
}

// TestPromptTemplate_AddShots tests the AddShots method of PromptTemplate.
func TestPromptTemplate_AddShots(t *testing.T) {
	// Initialize the Templates map
	err := json.Unmarshal([]byte(templates), &Templates)
	if err != nil {
		t.Fatalf("Error unmarshalling templates: %v", err)
	}

	// Get a template
	tpl, err := InitTemplate("chatml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Create a PromptTemplate instance
	promptTpl := &PromptTemplate{
		ID:         tpl.ID,
		Name:       tpl.Name,
		User:       tpl.User,
		Assistant:  tpl.Assistant,
		System:     tpl.System,
		Shots:      tpl.Shots,
		Stop:       tpl.Stop,
		Linebreaks: tpl.Linebreaks,
		AfterShot:  tpl.AfterShot,
		Prefix:     tpl.Prefix,
	}

	// Add multiple shots
	shots := []TurnBlock{
		{User: "What is the capital of France?", Assistant: "The capital of France is Paris."},
		{User: "What is the capital of Germany?", Assistant: "The capital of Germany is Berlin."},
	}
	promptTpl.AddShots(shots)

	// Check if the shots were added correctly
	if len(promptTpl.Shots) != 2 {
		t.Errorf("Expected 2 shots, got %d", len(promptTpl.Shots))
	}
	if promptTpl.Shots[0].User != shots[0].User {
		t.Errorf("Expected user to be '%s', got '%s'", shots[0].User, promptTpl.Shots[0].User)
	}
	if promptTpl.Shots[0].Assistant != shots[0].Assistant {
		t.Errorf("Expected assistant to be '%s', got '%s'", shots[0].Assistant, promptTpl.Shots[0].Assistant)
	}
	if promptTpl.Shots[1].User != shots[1].User {
		t.Errorf("Expected user to be '%s', got '%s'", shots[1].User, promptTpl.Shots[1].User)
	}
	if promptTpl.Shots[1].Assistant != shots[1].Assistant {
		t.Errorf("Expected assistant to be '%s', got '%s'", shots[1].Assistant, promptTpl.Shots[1].Assistant)
	}
	rendered := promptTpl.Render(true)
	fmt.Println(rendered)
}
