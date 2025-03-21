# Modprompt

A collection of prompt templates for language models.

## Features

- Classic template formats for different models.
- Easily modify and adapt templates on-the-fly.
- Few shots and conversation history support.

:books: [Api doc](https://pkg.go.dev/github.com/synw/gomodprompt)

## Installation

To use Modprompt, you'll need to install it first. You can do this by running:

```bash
go get -u github.com/synw/gomodprompt
```

## Basic Usage

Here's a basic example of how to use Modprompt to create a chat template and generate a prompt:

```go
package main

import (
	"fmt"
	"log"

	modprompt "github.com/synw/gomodprompt"
)

func main() {
	tpl, err := modprompt.InitTemplate("chatml")
	if err != nil {
		log.Fatal(err)
	}

	prompt := tpl.Prompt("Hello, how are you?")
	fmt.Println(prompt)
}
```

## Customization

Modprompt allows you to customize the chat templates to fit your specific use case. You can replace system messages, add extra text after user or assistant messages, and more. Here's an example of how to customize the Deepseek3 template:

```go
tpl, err := modprompt.InitTemplate("deepseek3")
if err != nil {
	log.Fatal(err)
}

// Replace the system message
tpl.ReplaceSystem("You are a helpful assistant.")

// Add extra text after the assistant message
tpl.AfterAssistant("\n\n")

// Add a shot to the template history
tpl.AddShot("What is the capital of France?", "The capital of France is Paris.")

prompt := tpl.Prompt("What is the largest planet in our solar system?")
fmt.Println(prompt)
```

In this example, we're replacing the system message, adding extra text after the assistant message, and adding a shot to the template history.

## Advanced Usage

### Handling Tool Calls

Modprompt supports tool calls, allowing you to integrate external tools into your prompts. Here's an example:

```go
tpl, err := modprompt.InitTemplate("mistral-system-tools")
if err != nil {
	log.Fatal(err)
}

// Add a tool definition
toolDef := map[string]interface{}{
	"name": "search",
	"description": "Search the web for information.",
}
tpl.AddTool(toolDef)

// Render a prompt with a tool call
prompt := tpl.Prompt("Find information about the capital of France.")
fmt.Println(prompt)
```

### Managing Multiple Shots

You can add multiple shots to the template history to simulate a conversation:

```go
tpl, err := modprompt.InitTemplate("chatml")
if err != nil {
	log.Fatal(err)
}

// Add multiple shots
shots := []modprompt.TurnBlock{
	{User: "What is the capital of France?", Assistant: "The capital of France is Paris."},
	{User: "Who wrote 'To Kill a Mockingbird'?", Assistant: "'To Kill a Mockingbird' was written by Harper Lee."},
}
tpl.AddShots(shots)

prompt := tpl.Prompt("What is the largest planet in our solar system?")
fmt.Println(prompt)
```

## History and Images

Modprompt also supports chat history and image data. You can push turns to the history and include image data in the turns. Here's an example:

```go
tpl, err := modprompt.InitTemplate("chatml")
if err != nil {
	log.Fatal(err)
}

// Push a turn to the history with images
tpl.PushToHistory(modprompt.HistoryTurn{
	User:      "What is the capital of France?",
	Assistant: "The capital of France is Paris.",
	Images: []modprompt.ImgData{
		{
			ID:   1,
			Data: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
		},
	},
})

prompt := tpl.Prompt("What is the largest planet in our solar system?")
fmt.Println(prompt)
```

In this example, we're pushing a turn with a user message, an assistant response, and an image to the chat history.

## Error Handling

Always check for errors when initializing templates or rendering prompts. Here's how you can handle errors:

```go
tpl, err := modprompt.InitTemplate("chatml")
if err != nil {
	// Handle error
	log.Fatalf("Failed to initialize template: %v", err)
}

prompt := tpl.Prompt("Hello, world!")
if prompt == "" {
	// Handle empty prompt
	log.Println("Generated an empty prompt.")
}
```
