package intellichunk

import (
	"fmt"
	"os"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

const (
	_defaultTokenEncoding = "cl100k_base"
)

// Tokenizer is a structure that contains the encoding information.
type Tokenizer struct {
	EncodingName string
}

// NewTokenizer initializes a new tokenizer with the encoding provided, if no encoding is provided default is used.
func NewTokenizer(encodingName ...string) Tokenizer {
	encoding := _defaultTokenEncoding
	if len(encodingName) > 0 {
		encoding = encodingName[0]
	}
	return Tokenizer{
		EncodingName: encoding,
	}
}

func (t Tokenizer) CountTokens(text string) (int, error) {
	tokens, err := t.Tokenize(text)
	if err != nil {
		return 0, err
	}
	return len(tokens), nil
}

// Tokenize text into tokens using Tiktoken.
func (t Tokenizer) Tokenize(text string) ([]int, error) {
	// Setting the cache directory
	err := os.Setenv("TIKTOKEN_CACHE_DIR", "cache/")
	if err != nil {
		return nil, fmt.Errorf("error setting environment variable: %w", err)
	}

	// Get the encoding
	tk, err := tiktoken.EncodingForModel(t.EncodingName)
	if err != nil {
		return nil, fmt.Errorf("tiktoken.EncodingForModel: %w", err)
	}

	// Tokenize the text
	tokens := tk.Encode(text, nil, nil)
	//tokensDecoded := tk.Decode(tokens)

	// If we need to Convert each token to string
	// tokenStrings := make([]string, len(tokens))
	// for i, token := range tokens {
	// 	tokenStrings[i] = fmt.Sprint(token)
	// }

	return tokens, nil
}
