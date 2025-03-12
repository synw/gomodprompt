package modprompt

import (
	"encoding/json"
	"fmt"
	"strings"
)

var Templates map[string]LmTemplate

func init() {
	var err error
	err = json.Unmarshal([]byte(templates), &Templates)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling templates: %v", err))
	}
}

func InitTemplate(name string) (LmTemplate, error) {
	tpl, ok := Templates[name]
	if !ok {
		return LmTemplate{}, fmt.Errorf("No template named %s", name)
	}
	return tpl, nil
}

func NewPromptTemplate(name string) (*PromptTemplate, error) {
	tmpl, err := InitTemplate(name)
	if err != nil {
		return nil, err
	}
	pt := &PromptTemplate{
		ID:         tmpl.ID,
		Name:       tmpl.Name,
		User:       tmpl.User,
		Assistant:  tmpl.Assistant,
		System:     tmpl.System,
		Shots:      tmpl.Shots,
		Stop:       tmpl.Stop,
		Linebreaks: tmpl.Linebreaks,
		AfterShot:  tmpl.AfterShot,
		Prefix:     tmpl.Prefix,
		ToolsDef:   tmpl.ToolsDef, // Copy ToolsDef
	}
	return pt, nil
}

func (tpl *PromptTemplate) ReplaceSystem(msg string) *PromptTemplate {
	tpl._replaceSystem = msg
	return tpl
}

func (tpl *PromptTemplate) AfterSystem(msg string) *PromptTemplate {
	tpl._extraSystem = msg
	return tpl
}

func (tpl *PromptTemplate) AfterAssistant(msg string) *PromptTemplate {
	tpl._extraAssistant = msg
	return tpl
}

func (tpl *PromptTemplate) ReplacePrompt(msg string) *PromptTemplate {
	tpl._replacePrompt = msg
	return tpl
}

func (tpl *PromptTemplate) AddTool(tool map[string]interface{}) *PromptTemplate {
	tpl.Tools = append(tpl.Tools, tool)
	return tpl
}

func (tpl *PromptTemplate) AddShot(user string, assistant string) *PromptTemplate {
	tpl.Shots = append(tpl.Shots, TurnBlock{User: user, Assistant: assistant})
	return tpl
}

func (tpl *PromptTemplate) AddShots(shots []TurnBlock) *PromptTemplate {
	tpl.Shots = append(tpl.Shots, shots...)
	return tpl
}

func (tpl *PromptTemplate) RenderShot(shot TurnBlock) string {
	var buf []string
	buf = append(buf, tpl.buildUserBlock(&shot.User))
	buf = append(buf, tpl.buildAssistantBlock(&shot.Assistant, true))

	if shot.Tool != "" && tpl.ToolsDef != nil {
		toolCall := tpl.RenderToolCall(shot.Tool)
		toolResponse := strings.Replace(tpl.ToolsDef.Response, "{tools_response}", shot.Tool, 1)
		buf = append(buf, toolCall, toolResponse)
	}

	if tpl.AfterShot != nil {
		buf = append(buf, *tpl.AfterShot)
	}
	return strings.Join(buf, "")
}

func (tpl *PromptTemplate) Render(skipEmptySystem bool) string {
	var res string
	if tpl.System != nil {
		res = tpl.buildSystemBlock(skipEmptySystem)
	}
	for _, shot := range tpl.Shots {
		res += tpl.RenderShot(shot)
	}
	res += tpl.buildUserBlock(nil)
	res += tpl.buildAssistantBlock(nil, false)
	return res
}

func (tpl *PromptTemplate) Prompt(msg string) string {
	txt := tpl.Render(true)
	txt = strings.Replace(txt, "{prompt}", msg, 1)
	return txt
}

func (tpl *PromptTemplate) PushToHistory(turn HistoryTurn) *PromptTemplate {
	tpl.Shots = append(tpl.Shots, TurnBlock{
		User:      turn.User,
		Assistant: turn.Assistant,
		Tool:      turn.Tool,
	})
	return tpl
}

func (tpl *PromptTemplate) buildSystemBlock(skipEmptySystem bool) string {
	var res string
	if tpl.System != nil {
		res = tpl.System.Schema
		if tpl._replaceSystem != "" {
			res = strings.Replace(res, "{system}", tpl._replaceSystem, -1)
		} else if tpl.System.Message != nil {
			res = strings.Replace(res, "{system}", *tpl.System.Message, -1)
		}
		if tpl._extraSystem != "" {
			res += tpl._extraSystem
		}
		if tpl.Linebreaks != nil && tpl.Linebreaks.System != nil {
			res += strings.Repeat("\n", *tpl.Linebreaks.System)
		}
		if strings.Contains(res, "{tools}") && tpl.ToolsDef != nil {
			toolsBlock := tpl.buildToolsBlock()
			res = strings.Replace(res, "{tools}", toolsBlock, 1)
		}
	} else if !skipEmptySystem {
		res = tpl.System.Schema
	}
	return res
}

func (tpl *PromptTemplate) buildUserBlock(msg *string) string {
	var buf []string
	userBlock := tpl.User
	if tpl._replacePrompt != "" {
		userBlock = strings.Replace(userBlock, "{prompt}", tpl._replacePrompt, -1)
	}
	if msg != nil {
		userBlock = strings.Replace(userBlock, "{prompt}", *msg, -1)
	}
	buf = append(buf, userBlock)
	if tpl.Linebreaks != nil && tpl.Linebreaks.User != nil {
		buf = append(buf, strings.Repeat("\n", *tpl.Linebreaks.User))
	}
	return strings.Join(buf, "")
}

func (tpl *PromptTemplate) buildAssistantBlock(msg *string, isShot bool) string {
	var buf []string
	assistantBlock := tpl.Assistant
	buf = append(buf, assistantBlock)
	if tpl.Linebreaks != nil && tpl.Linebreaks.Assistant != nil {
		buf = append(buf, strings.Repeat("\n", *tpl.Linebreaks.Assistant))
	}
	if msg != nil {
		buf = append(buf, *msg)
	}
	if (tpl._extraAssistant != "") && !isShot {
		buf = append(buf, tpl._extraAssistant)
	}
	return strings.Join(buf, "")
}

func (tpl *PromptTemplate) buildToolsBlock() string {
	if tpl.ToolsDef == nil || len(tpl.Tools) == 0 {
		return ""
	}
	toolsJSON, _ := json.Marshal(tpl.Tools)
	return strings.Replace(tpl.ToolsDef.Def, "{tools}", string(toolsJSON), 1)
}

// New method to handle tool calls
func (tpl *PromptTemplate) RenderToolCall(toolName string) string {
	if tpl.ToolsDef == nil || tpl.ToolsDef.Call == "" {
		return ""
	}
	return strings.Replace(tpl.ToolsDef.Call, "{tool}", toolName, 1)
}
