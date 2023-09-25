package templateprompt

// Constants defining the various prompt templates.
const (
	HelperAgentPrompt = `Your are helper agent, provide helpful answers solely using the facts provided below:

	Topic Details: {{.Details}} 
	
	Helpful Facts: {{.SSContent}}`
)
