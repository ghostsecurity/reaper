package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/tools/fuzz"
	"github.com/ghostsecurity/reaper/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/instructor-ai/instructor-go/pkg/instructor"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

var (
	activeSessions = make(map[uint]*types.AgentSession)
	mutex          = &sync.Mutex{}
)

// AgentManager allows agents to send messages to the chat
type AgentManager struct {
	Ctx        *fiber.Ctx
	Pool       *websocket.Pool
	DB         *gorm.DB
	SessionID  uint
	Author     uint
	AuthorRole types.UserRole
	Prompt     string
}

type Attack struct {
	AttackName string `json:"attack_name"  jsonschema:"title=the attack name,default=BOLA,description=The name of the attack,enum=BOLA,example=BOLA"`
	Domain     string `json:"domain"       jsonschema:"title=the domain,default=ghostbank.net,description=The domain of the attack,example=ghostbank.net"`
	Report     bool   `json:"report"       jsonschema:"title=the report,default=true,description=Should a report be generated?,example=true,example=false"`
}

type WorstKey struct {
	Name string `json:"name" jsonschema:"title=key name,description=The key to focus on in the fuzz attack"`
}

type Report struct {
	Content string `json:"content"          jsonschema:"title=the content,description=The content of the report in markdown as a string"`
}

// NOTE: This is an experimental AI Agent
// This agent is attempting to demonstrate tool usage by GPT4x models
// and report writing. This is not a production-ready agent and should
// not be considered secure or reliable for any purpose.
func (manager *AgentManager) StartAgent() {
	// Ensure only one agent is running at a time for a session
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := activeSessions[manager.SessionID]; !exists {
		activeSessions[manager.SessionID] = &types.AgentSession{
			ID:        manager.SessionID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    "active",
		}
		go manager.runAgent()
	} else {
		fmt.Printf("[***] agent session %d already running\n", manager.SessionID)
		manager.sendAgentMessage("Agent session already running. Please wait until it completes.")
	}
}

func (manager *AgentManager) sendAgentMessage(content string) {
	message := models.AgentSessionMessage{
		AuthorID:       manager.Author,
		AuthorRole:     manager.AuthorRole,
		AgentSessionID: uint(manager.SessionID),
		Content:        content,
	}

	resp := manager.DB.Create(&message)
	if resp.Error != nil {
		fmt.Printf("Error saving agent message: %v\n", resp.Error)
	}

	msg := &types.AgentSessionMessage{
		Type:       types.MessageTypeAgentSessionMessage,
		AuthorID:   manager.Author,
		AuthorRole: manager.AuthorRole,
		SessionID:  uint(manager.SessionID),
		Content:    content,
	}

	manager.Pool.Broadcast <- msg
}

func process_prompt(prompt string) Attack {
	ctx := context.Background()

	client := instructor.FromOpenAI(
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		instructor.WithMode(instructor.ModeJSONSchema),
		instructor.WithMaxRetries(3),
	)

	var attack Attack
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
		&attack,
	)
	_ = resp
	if err != nil {
		panic(err)
	}

	fmt.Printf("Model returned Attack Name: %s, Domain: %s, Report: %t\n", attack.AttackName, attack.Domain, attack.Report)

	return attack
}

func generate_report_content(function_model Attack, findings []models.FuzzResult) string {
	ctx := context.Background()

	client := instructor.FromOpenAI(
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		instructor.WithMode(instructor.ModeJSONSchema),
		instructor.WithMaxRetries(3),
	)
	system_prompt := "You are an expert web application security analyst/engineer who has just run a sophisticated attack against an application to validate that a broken object level access (BOLA) flaw exists in a web app. Your job is to strike a balance between ease of understanding and technical accuracy when writing up this report.  Do NOT hallicinate or make up anything that you can't derive from the results presented in this prompt.  You want to impress the reader with prose that is concise and clear."

	user_prompt := "I have provided you with a JSON array of approx 20 results from a successful broken object level access (BOLA) attack. The results are in this code block:\n\n"
	jsonData, err := json.MarshalIndent(findings, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v", err)
		return ""
	}
	jsonStr := string(jsonData)
	user_prompt += fmt.Sprintf("The domain scanned was: %s\n", function_model.Domain)
	user_prompt += "```\n"
	user_prompt += jsonStr
	user_prompt += "```\n\n"

	user_prompt += "I want you to write a 2-3 page technical write-up formatted in markdown with sections like summary, date/time, an explanation of what a BOLA attack is, an detailed explanation of what happened with these 20 requests, and what should be done to fix/remediate the vulnerability (including what types of things to fix in code and what might be done with a web application firewall).\n\n"
	user_prompt += "The following sections and descriptions are desired.\n\n"
	user_prompt += "- **Title** - A short but useful title of the report\n"
	user_prompt += "- **Executive Summary** - Short description of the attack, what the attack is, and its characteristics that summarizes the risk and impact and what the successful attack was able to do.\n"
	user_prompt += "- **Detailed Explanation** - In this section, provides a technical Requests/Response Analysis\n"
	user_prompt += "- **Remediation Guidance** - What to fix in code/logic, web application firewalls/waf, or other mitigations.\n\n"

	user_prompt += "Your response will be in markdown format only with no extra explanation or messages to me. Only the report content should be returned. Be sure to have two newlines before code blocks and lists.\n\n"
	user_prompt += "Ensure that code block start and end markers are on their own lines.\n\n"
	user_prompt += "Do not include a full listing of the requests/responses in the output. I will add them myself.\n"

	fmt.Printf("System Prompt: %s\n", system_prompt)
	fmt.Printf("User Prompt: %s\n", user_prompt)

	var report Report
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system_prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: user_prompt,
				},
			},
		},
		&report,
	)
	_ = resp
	// fmt.Printf("Model returned resp: %v\n", resp)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Model returned Report: %s\n", report.Content)
	content := report.Content

	content += "\n\n# Appendix\n\n"
	content += "The Request/Response Listing:\n\n"
	for _, finding := range findings {
		content += "```\n"
		content += "# Request\n"
		content += fmt.Sprintf("%s\n\n", finding.Request)
		content += "# Response\n"
		content += fmt.Sprintf("%s\n\n", finding.Response)
		content += "```\n\n"
	}

	return content
}

func excludedFuzzKeys(request models.Request) []string {
	// Extract the keys from the body
	ctx := context.Background()

	client := instructor.FromOpenAI(
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		instructor.WithMode(instructor.ModeJSONSchema),
		instructor.WithMaxRetries(3),
	)

	prompt := fmt.Sprintf("The post body from a valid %s API request to host: %s with content-type %s is %v\n\nThe headers are %v\n\n Return the body key name that would be most dangerous/impactful in terms of a succesful BOLA/IDOR fuzzing attack that would likely cause the most impact.", request.Method, request.Host, request.ContentType, request.Body, request.Headers)
	var dangerousKey WorstKey
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
		&dangerousKey,
	)
	_ = resp
	if err != nil {
		panic(err)
	}
	fmt.Printf("Model returned key to fuzz: %v\n", dangerousKey.Name)

	// Get the keys from the body
	var postbody map[string]interface{}
	fmt.Printf("body: %v\n", request.Body)
	if err := json.Unmarshal([]byte(request.Body), &postbody); err != nil {
		fmt.Printf("failed to parse body keys: %v", err)
		return []string{}
	}
	fmt.Printf("Post body: %v\n", postbody)

	// Remove the most dangerous key from the list
	var keys []string
	for key := range postbody {
		if key != dangerousKey.Name {
			keys = append(keys, key)
		}
	}

	return keys
}

// runAgent performs the agent's tasks
func (manager *AgentManager) runAgent() {
	defer stopAgent(manager.SessionID)

	// Agent logic goes here
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("[***] OPENAI_API_KEY not set")
		manager.sendAgentMessage("OPENAI_API_KEY not set as an environment variable")
		return
	}

	// process the chat message prompt from the user into a structured model so we can determine the attack type,
	// the domain to run the attack on, and whether to generate a report
	manager.sendAgentMessage("Processing prompt...")
	function_model := process_prompt(manager.Prompt)
	manager.sendAgentMessage("Determining attack to run...")

	if function_model.AttackName == "BOLA" {
		content := fmt.Sprintf("Running %s attack on %s and generating a report: %t", function_model.AttackName, function_model.Domain, function_model.Report)
		manager.sendAgentMessage(content)

		manager.sendAgentMessage("Fetching potential endpoints...")
		// Fetch the requests in the domain with the POST method and application/json content type
		var requests []models.Request
		var target_request models.Request

		result := manager.DB.Where("host LIKE ? AND method = ? AND content_type = ?", "%"+function_model.Domain, "POST", "application/json").Order("created_at DESC").Find(&requests)
		if result.Error != nil {
			manager.sendAgentMessage("Sorry, no potential endpoints found.")
			return
		}
		if len(requests) > 0 {
			target_request = requests[0]
			manager.sendAgentMessage("Susceptible endpoint found!")
		} else {
			manager.sendAgentMessage("Sorry, no potential endpoints found.")
			return
		}
		fmt.Printf("[***] Targeting endpoint %d\n", target_request.ID)

		// Ask the model for the attack parameters for the request
		excludedKeys := excludedFuzzKeys(target_request)
		fmt.Printf("[***] Excluded keys: %v\n", excludedKeys)
		time.Sleep(8 * time.Second)

		// Run the BOLA attack
		manager.sendAgentMessage("Running the attack...")
		err := fuzz.CreateAttack(function_model.Domain, excludedKeys, manager.Pool, manager.DB, 0, 0, 0)
		if err != nil {
			slog.Error("Failed to create fuzz attack", "error", err)
		}
		manager.sendAgentMessage("Attack completed successfully")

		// Fetch the findings
		manager.sendAgentMessage("Fetching attack findings...")
		var findings []models.FuzzResult
		findingsResults := manager.DB.Where("hostname = ?", target_request.Host).Order("created_at DESC").Limit(5).Find(&findings)
		if findingsResults.Error != nil {
			fmt.Printf("[***] Error fetching attack findings: %v\n", findingsResults.Error)
			manager.sendAgentMessage("Error fetching findings.")
			return
		}
		if len(findings) == 0 {
			fmt.Printf("[***] No findings found.\n")
			manager.sendAgentMessage("Sorry, no findings found.")
			return
		}
		fmt.Printf("Findings result %v\n", findingsResults)

		// Generate a report
		if function_model.Report {
			manager.sendAgentMessage("Generating attack analysis report. This may take some time to process..")
			report_content := generate_report_content(function_model, findings)

			manager.sendAgentMessage("Saving the report...")
			report := models.Report{
				Domain:   function_model.Domain,
				Markdown: report_content,
			}
			resp := manager.DB.Create(&report)
			if resp.Error != nil {
				fmt.Printf("Error saving report: %v\n", resp.Error)
				manager.sendAgentMessage("Error saving report")
				return
			}
			manager.sendAgentMessage("Report saved successfully")
		} else {
			// No report requested
			manager.sendAgentMessage("No findings or report generation was not requested")
		}
		manager.sendAgentMessage("Done.")
	} else {
		// Attack not supported
		msg := fmt.Sprintf("Attack not supported: %s", function_model.AttackName)
		manager.sendAgentMessage(msg)
	}
}

// stopAgent cleans up the session
func stopAgent(sessionID uint) {
	mutex.Lock()
	defer mutex.Unlock()

	if session, exists := activeSessions[sessionID]; exists {
		session.UpdatedAt = time.Now()
		session.Status = "stopped"
		delete(activeSessions, sessionID)
	}
}
