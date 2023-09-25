
# Intellichunk AI

**Retrieval Augmented Generation (RAG) / Go Based Lightweight Langchain Alternative**

⚠️ **Note:** This repository is in developmental phase and has been made public in the hope that it might serve as a helpful starting point for others. There are numerous areas for improvement and addition, so feel free to fork and submit Pull Requests.

The initial motivation behind Intellichunk was to develop a lightweight langchain alternative that facilitates AI based smart chunking and building a vector database from the Command Line Interface (CLI). Moreover, it aims to harness the Retrieval Augmented Generation (RAG) ability through an API, making it accessible for front-end calls.

Although deploying LLMs on their own offers value, their full potential is unlocked when amalgamated with other computational or informational assets. This library is curated to simplify the foundation for developers aiming to construct applications of this nature.

This repository is engineered to integrate conversational chat functionalities by leveraging a vectorstore. Moreover, it furnishes a Command Line Interface (CLI) tailored for Intellichunk, enabling the streamlined processing, vectorization of voluminous textual data, and formulating vector stores from sources. While the CLI deals with sporadic tasks, the API endpoint persistently resides in the cloud for uninterrupted access.

 
  
## Features
- **Initialization Flexibility:** Initialize LLM with various options including Temperature, TopP/Nucleus Sampling, Model Selection, and ChatHistory.
- **Dynamic Prompt Templating:** Customize and use multiple prompts as per your needs.
- **Serverless Ready:** Easy to deploy on Google Cloud Run, AWS Lambda, Azure Functions, and other serverless providers.
- **Versatile Run Environment:** Operate the application both a Command Line Interface (CLI) and an API server.
- **Tiktoken Tokenizer:** Token counting without making an API call using tiktoken. Choose any encoding. 
- **Efficient Text Chunking:** Optimize the processing of long documents by breaking them down into meaningful, manageable chunks.
- **Metadata-rich Chunks:** Each chunk comes with vital metadata, including questions linked to unique entities and keyword generation.
- **Batch Processing:** Efficiently processes large volumes of textual data and creates vector stores from sources. 
- **CLI AI Chat:** Initiate and conduct conversations with stored chat history directly via the terminal.
- **API Capabilities:** Offers API endpoints for direct querying, extracting chunks, embedding, and vectorizing large text inputs.
- **Clear Documentation:** Provides explicit CLI command instructions, API references, and source file templates.
- **Environment Configuration:** Leverage `.env` files for convenient and secure management of API keys and related settings.


## + Intellichunk

This package optimizes text chunking for long factual documents and articles by dividing them into smaller, manageable pieces using LLM. It includes essential metadata with each chunk, such as relevant questions extracted from unique entities, and generates keywords, though they are not vectorized due to low performance within LLM. The metadata supports standard database integration if needed, and the package efficiently performs batch embeddings generation for faster vector generation and data processing.

'

  

'

  

## CLI Commands
 
#### Conversation
The `conversation` command takes a class name and a query question.
It starts a conversation with chat history through the terminal.
```shell

go  run  .  conversation "ClassID"  "Tell me about x"

```

##### Example
```text
{"message":"Using RUN_ENV=local environment variables","severity":"INFO"}
LLM: The challenges in addressing public health risks associated with climate change are numerous. Some of the key challenges include: 1. Socioeconomic factors, Age, Marginalized groups
You: explain #3 a bit
LLM: Certainly! When we talk about vulnerable populations in the context of climate change and public health, we are referring to groups of people who are disproportionately affected by the health risks associated with climate change. Here are some key points to consider:
.
.
.
in energy and climate research.
You: exit
Exiting the conversation...
```



#### Intellichunk Add

The 'intellichunk add' command takes a class name and the path to a folder containing text files as input.

It iterates over the text files within the folder, reads articles/sources from each file.

This function is useful for batch processing of large text files and storing their context in a structured and accessible format.

path is relative to the project folder.

  

```shell

go  run  .  intellichunk  add  "ClassID"  "/files"

```

  

Sample response in the terminal from this command:

  

```text

:: Analysing... > sourcesCUT.txt

::::: Processing Article... > The economic transformation: What would change in the net-zero transition

--------> Intellichunked into 6 nodes successfully.

--------> Added to the vectorstore. Vector Object IDs: 39082230-871f-462c-a136-09377eef5b26, ...

::::: Processing Article... > Decarbonization and the Benefits of Tackling Climate Change

--------> Intellichunked into 3 nodes successfully.

--------> Added to the vectorstore. Vector Object IDs: 309969ee-a5c0-491c-a390-c4763eb05e2b, ...

:: Analysing... > anotherSource.txt

```




##### Source File Template

The code expects the following format within `.txt` files:

  

```Text

Title:The economic transformation: What would change in the net-zero transition

RefURL:https://www.mckinsey.com/capabilities/sustainability/our-insights/the-economic-transformation-what-would-change-in-the-net-zero-transition

Content:Long long text content...

  

Title:Another article title

RefURL:https://wwwss

Content:Long long text content...

```
#### Run API
The `runapi` command starts the api server. It's useful in local environments.
```shell

go  run  .  runapi

```
  

## API Reference

  

#### Ask a question about given class/topic.

  

```http

GET /conversation

```

  

| Parameter | Type | Description |
| :-------- | :------- | :------------------------- |
| `ConversationID` | `string` | ConversationID on the frontend for back reference. This will be returned. |
| `ClassID` | `string` | **Required**. Classes are queried on this ID within the vector database. |
| `ChatHistory` | `array` | Conversation history in string array format. Can be empty. Even indexed strings are the user’s input;  and the odd index strings are the llm response |
| `Query` | `string` | **Required**. The question |

  
  

Sample Response:

```

{

"ConversationID": "1fa23fdaa45",

"ClassID": "camp001",

"Query": "How this decarbonizes the economy??",

"Answer": {

"Response": "This Class reduces our reliance on fossil fuels and shift towards cleaner energy sources. This is important because it helps combat climate change by doing XYZ, which is one of the most significant challenges facing our planet.........",

"Sources": ["http://somearticle.com/decarbonization", "http://wikipedia.com/Climate_change"]

},

"Suggestions": ["What steps can we take to decarbonize the economy?", "What are the primary sources of carbon emissions?", "What are the impacts of climate change?"]

}

```

  

#### Split large document into meaningful chunks, embed and vectorize them. Returns ids from vector database.

  

```http

POST /intellichunk/add

```

  

| Parameter | Type | Description |
| :-------- | :------- | :-------------------------------- |
| `ClassName` | `string` | **Required**. ClassID / Class Name within vector database. |
| `LongText` | `string` | **Required**. Large chunk of any text data/document.|

  
  

  
  

### Prerequisites

  
#### .env file
```file
OPENAI_API_KEY=sk-f9Lxxxx

WEAVIATE_API_KEY=Fndpxxxxx

WEAVIATE_URL=wfac-bhpx6tjb.weaviate.network
```

#### Sample VSCode launch.json file for debugging
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch intellichunk Add",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "args": [
                "intellichunk",
                "add",
                "Class_test",
                "files/tests/",
                "-s"
            ],
            "env": {},
            "showLog": true
        },
        {
            "name": "Launch - Conversation",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "args": [
                "conversation"
            ],
            "env": {},
            "showLog": true,
            "console": "integratedTerminal"
        },
        {
            "name": "Launch - API Run",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "args": [
                "runapi"
            ],
            "env": {},
            "showLog": true,
            "console": "integratedTerminal"
        }
    ]
}
```
