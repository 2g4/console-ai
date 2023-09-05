package data

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
)

// Build our new spinner
var s = spinner.New(spinner.CharSets[7], 100*time.Millisecond)

func PrintAiAnswer(input string) {
	colorPurple := "\033[36m"
	colorReset := "\033[0m"
	result := markdown.Render(input, 80, 7)

	fmt.Println(string(colorPurple), "  AI: ", string(colorReset))
	fmt.Println(string(result))
}

// Start delay spinner
func StartAiThinking() {
	s.Prefix = "   AI: "
	s.Color("yellow")
	s.Start()                   // Start the spinner
	time.Sleep(1 * time.Second) // Run for some time to simulate work
}

// Stop delay spinner
func StopAiThinking() {
	s.Stop()
}

// Converts string to int
func StrToInt(s string) int {
	marks, _ := strconv.ParseInt(s, 0, 64)
	return int(marks)
}

// Converts string to float
func StrToFloat(s string) float32 {
	marks, _ := strconv.ParseFloat(s, 32)
	return float32(marks)
}

func PromptInit() {
	answer := PromptCustomOrDefault("Settings not found. \nDo you want to initialize the app? (Y/n)", "Y")
	fmt.Println()
	if strings.HasPrefix(strings.ToLower(answer), "y") {
		PopulateSettings()
	} else {
		fmt.Println("No? I'm out, bye!")
		// Exit the program
		os.Exit(0)
	}
}

func PromptCustomOrDefault(question string, defaultValue string) string {
	validate := func(input string) error {
		if input == "" && defaultValue == "" {
			return errors.New("Invalid string")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | green }} ",
	}

	fmt.Println()
	fmt.Println(question)

	prompt := promptui.Prompt{
		Label:     ">>",
		Validate:  validate,
		Templates: templates,
		Default:   defaultValue,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	if result != "" {
		return result
	}

	return defaultValue
}

func PopulateSettings() {
	// Create settings table
	CreateTable()
	// Get default settings to give the user a hint or auto-fill
	settings := GetDefaultSettings()
	// Ask the user for the settings

	// Get and encrypt API key
	apiKey := PromptCustomOrDefault("Secret API key (required)", "")
	encryptedKey, err := EncryptMessage(GetEncryptionKey(), apiKey)
	if err != nil {
		log.Fatalln(err)
	}
	InsertKey("apiKey", encryptedKey)

	InsertKey("apiUrl", PromptCustomOrDefault("API URL (optional)", settings.ApiURL))
	InsertKey("model", PromptCustomOrDefault("Model, default is GPT4, you can set 'gpt-3.5-turbo' instead (optional)", settings.Model))
	InsertKey("systemMessage", PromptCustomOrDefault("System message (optional)", settings.SystemMessage))
	InsertKey("maxTokens", PromptCustomOrDefault("Max tokens (optional)", fmt.Sprint(settings.MaxTokens)))
	InsertKey("temperature", PromptCustomOrDefault("Temperature (optional)", fmt.Sprint(settings.Temperature)))
	InsertKey("frequencyPenalty", PromptCustomOrDefault("Frequency penalty (optional)", fmt.Sprint(settings.FrequencyPenalty)))
	InsertKey("presencePenalty", PromptCustomOrDefault("Presence penalty (optional)", fmt.Sprint(settings.PresencePenalty)))
	fmt.Println("\nAll set, now you can ask AI.")
}

func OpenAiRequest(query string) string {

	settings := FetchSettings()

	// Add system message if messages is empty
	if len(GetMessages()) == 0 {
		AppendMessage("system", settings.SystemMessage)
	}

	// Populates user message
	AppendMessage("user", query)

	pb := &ApiPostBody{
		Messages:         GetMessages(),
		MaxTokens:        settings.MaxTokens,
		Temperature:      settings.Temperature,
		FrequencyPenalty: settings.FrequencyPenalty,
		PresencePenalty:  settings.PresencePenalty,
		Model:            settings.Model,
	}

	jsonData, err := json.Marshal(pb)
	if err != nil {
		log.Fatalf("Could not marshal JSON: %s", err)
	}

	request, error := http.NewRequest("POST", settings.ApiURL, bytes.NewBuffer(jsonData))
	authorization := "Bearer " + settings.ApiKey
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authorization)

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		fmt.Println("Request error")
		return ""
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body")
		return ""
	}

	d := OpenApiResponse{}

	if err := json.Unmarshal(bodyBytes, &d); err != nil {
		fmt.Println("Failed to unmarshal response body")
		return ""
	}

	if d.Choices == nil || len(d.Choices) < 0 {
		fmt.Println("Failed to get choices")
		return ""
	}

	// Populates assistant message'
	AppendMessage(d.Choices[0].Message.Role, d.Choices[0].Message.Content)

	return d.Choices[0].Message.Content
}

func InitOrDie() {
	answer := PromptCustomOrDefault("\nDo you want to initialize the app? (Y/n)", "Y")
	fmt.Println()
	if strings.HasPrefix(strings.ToLower(answer), "y") {
		PopulateSettings()
	} else {
		fmt.Println("No? I'm out, bye!")
		// Exit the program
		os.Exit(0)
	}
}

func ReadQuestion() {
	validate := func(input string) error {
		if input == "" {
			return errors.New("Invalid string")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | green }} ",
	}

	fmt.Println()

	prompt := promptui.Prompt{
		Label:     "Human:",
		Validate:  validate,
		Templates: templates,
	}

	input, err := prompt.Run()

	if err != nil {
		fmt.Printf("Bye")
		return
	}

	RunQuery(input)

}

func RunQuery(input string) {
	// Start the spinner
	StartAiThinking()
	// Gets response
	response := OpenAiRequest(input)
	// Removes spinner
	StopAiThinking()
	PrintAiAnswer(response)
	// Read next question
	ReadQuestion()
}

func GetEncryptionKey() []byte {
	hostname, _ := os.Hostname()
	base64Hostname := base64.StdEncoding.EncodeToString([]byte(hostname + "just-in-case"))

	return []byte(
		// Gets 16 characters from the base64 encoded hostname
		base64Hostname[:16])
}
