package models

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Parameters is an object describing the function.
	// You can pass a []byte describing the schema,
	// or you can pass in a struct which serializes to the proper JSON schema.
	Parameters any `json:"parameters"`
}

// Deprecated: use FunctionDefinition instead.
type FunctionDefine = FunctionDefinition

type DataType string

const (
	Object  DataType = "object"
	Number  DataType = "number"
	Integer DataType = "integer"
	String  DataType = "string"
	Array   DataType = "array"
	Null    DataType = "null"
	Boolean DataType = "boolean"
)

// Definition is a struct for describing a JSON Schema.
type Definition struct {
	// Type specifies the data type of the schema.
	Type DataType `json:"type,omitempty"`
	// Description is the description of the schema.
	Description string `json:"description,omitempty"`
	// Enum is used to restrict a value to a fixed set of values. It must be an array with at least
	// one element, where each element is unique. You will probably only use this with strings.
	Enum []string `json:"enum,omitempty"`
	// Properties describes the properties of an object, if the schema type is Object.
	Properties map[string]Definition `json:"properties"`
	// Required specifies which properties are required, if the schema type is Object.
	Required []string `json:"required,omitempty"`
	// Items specifies which data type an array contains, if the schema type is Array.
	Items *Definition `json:"items,omitempty"`
	// MinLength specifies the minimum length of the string.
	MinLength int `json:"minLength,omitempty"`
	// MaxLength specifies the maximum length of the string.
	MaxLength int `json:"maxLength,omitempty"`
}

// ConversationRequest handles a conversation request.
type ConversationRequest struct {
	ConversationID string   `json:"ConversationID"`
	ClassID        string   `json:"ClassID"`
	ChatHistory    []string `json:"ChatHistory"`
	Query          string   `json:"Query"`
}

// ConversationResponse is sent back as response from the API.
type ConversationResponse struct {
	ConversationID string   `json:"ConversationID"`
	ClassID        string   `json:"ClassID"`
	Query          string   `json:"Query"`
	Answer         Answer   `json:"Answer"`
	Suggestions    []string `json:"Suggestions"`
}

// Answer struct holds answer and relevant references
type Answer struct {
	Answer  string   `json:"Answer"`
	Sources []string `json:"Sources"`
}

// Class struct to hold static Class info.
type Class struct {
	Description string
}

// ContainerNode is smaller chunks that compose the whole container.
type ContainerNode struct {
	Content    string   `json:"content"`
	Keywords   []string `json:"keywords"`
	Questions  []string `json:"questions"`
	NodeNumber int      `json:"sectionNumber"`
}

// Container Holds information about a particular context, this can be a document or any other text.
type DataContainer struct {
	Title      string          `json:"title"`
	Summary    string          `json:"summary"`
	AbstactSum string          `json:"abstract_description"`
	Nodes      []ContainerNode `json:"nodes"`
}

type ContainerNodeVector struct {
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	AbstactSum   string    `json:"abstract_sum"`
	Content      string    `json:"content"`
	Keywords     []string  `json:"keywords"`
	Questions    []string  `json:"questions"`
	NodeNumber   int       `json:"section_number"`
	RefTitle     string    `json:"reference_title"`
	ReferenceURL string    `json:"reference_url"`
	Embedding    []float32 `json:"embedding"`
}

type IntellichunkRequest struct {
	ClassName string
	LongText  string
}

// More generic data holders
type GeneralDataHolder struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"`
}
