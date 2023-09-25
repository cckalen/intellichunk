package intellichunk_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apsystole/log"
	"github.com/cckalen/intellichunk/config"
	"github.com/cckalen/intellichunk/internal/intellichunk"
	"github.com/cckalen/intellichunk/internal/llm"
	"github.com/cckalen/intellichunk/internal/models"

	"github.com/hlindberg/testutils"
	"github.com/sashabaranov/go-openai"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Errorf("Load env Error: %s", err)

	}
}

func Test_Tokenizer(t *testing.T) {
	te := intellichunk.NewTokenizer("gpt-3.5-turbo")
	text := "Hello, world!"
	tokens, err := te.Tokenize(text)
	testutils.CheckNotError(err, t)
	testutils.CheckNotNil(tokens, t)
}

func Test_CountTokens(t *testing.T) {
	tik := intellichunk.NewTokenizer("gpt-3.5-turbo")
	longtext := "Ladakh is a region administered by India as a union territory[1]"

	tokenCount, err := tik.CountTokens(longtext)
	testutils.CheckNumericGreater(0, tokenCount, t)
	testutils.CheckNotError(err, t)
}

func Test_EmbeddingsWithTokens(t *testing.T) {
	te := intellichunk.NewTokenizer("gpt-3.5-turbo")
	text := "This text going to be tokenized first and then embedded"

	tokens, err := te.Tokenize(text)
	testutils.CheckNotError(err, t)
	testutils.CheckNotNil(tokens, t)

	// Create a new OpenAI language model instance.
	languageModel := llm.NewOpenAI()

	embedResp, err := languageModel.GenerateEmbeddings(context.Background(), tokens, openai.AdaEmbeddingV2, "system")
	testutils.CheckNotError(err, t)
	log.Print(embedResp)
}

func Test_EmbeddingsFromText(t *testing.T) {

	// Create a new OpenAI language model instance.
	languageModel := llm.NewOpenAI()
	textToEmbed := []string{"Soccer players", "interesting facts", "some other random topic to think about"}
	embedResp, err := languageModel.GenerateMultipleEmbeddingsFromText(context.Background(), textToEmbed)
	fmt.Print(embedResp)
	testutils.CheckNotError(err, t)
}

func Test_SplitTextIntoContainer(t *testing.T) {

	//longtext := "Ladakh is a region administered by India as a union territory[1] and constitutes an eastern portion of the larger Kashmir region that has been the subject of a dispute between India and Pakistan since 1947 and India and China since 1959.[2] Ladakh is bordered by the Tibet Autonomous Region to the east, the Indian state of Himachal Pradesh to the south, both the Indian-administered union territory of Jammu and Kashmir and the Pakistan-administered Gilgit-Baltistan to the west, and the southwest corner of Xinjiang across the Karakoram Pass in the far north. It extends from the Siachen Glacier in the Karakoram range to the north to the main Great Himalayas to the south.[11][12] The eastern end, consisting of the uninhabited Aksai Chin plains, is claimed by the Indian Government as part of Ladakh, and has been under Chinese control since 1962.[13]		In the past, Ladakh gained importance from its strategic location at the crossroads of important trade routes,[14] but as Chinese authorities closed the borders between Tibet Autonomous Region and Ladakh in the 1960s, international trade dwindled. Since 1974, the Government of India has successfully encouraged tourism in Ladakh. As Ladakh is strategically important, the Indian military maintains a strong presence in the region. "
	file, err := os.Open("cache/longText.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	text := string(bytes)

	resp, _ := intellichunk.SplitTextIntoContainerNodes(text)

	fmt.Println(resp)
	//testutils.CheckNotNil(resp, t)
}

func Test_GenerateContainers(t *testing.T) {

	splittedText := `{
		"abstract_description": "The input provides information about Ladakh, a region administered by India as a union territory. It discusses the disputed status of Ladakh between India and Pakistan, as well as India and China. The input also mentions the geographical borders of Ladakh and its strategic importance.",
		"nodes": [
		  {
			"content": "Ladakh is a region administered by India as a union territory and constitutes an eastern portion of the larger Kashmir region that has been the subject of a dispute between India and Pakistan since 1947 and India and China since 1959. Ladakh is bordered by the Tibet Autonomous Region to the east, the Indian state of Himachal Pradesh to the south, both the Indian-administered union territory of Jammu and Kashmir and the Pakistan-administered Gilgit-Baltistan to the west, and the southwest corner of Xinjiang across the Karakoram Pass in the far north. It extends from the Siachen Glacier in the Karakoram range to the north to the main Great Himalayas to the south.",
			"keywords": ["Ladakh", "India", "union territory", "Kashmir region", "dispute"],
			"questions": ["What is the disputed status of Ladakh?", "What are the geographical borders of Ladakh?"]
		  },
		  {
			"content": "The eastern end, consisting of the uninhabited Aksai Chin plains, is claimed by the Indian Government as part of Ladakh and has been under Chinese control since 1962. In the past, Ladakh gained importance from its strategic location at the crossroads of important trade routes, but as Chinese authorities closed the borders between Tibet Autonomous Region and Ladakh in the 1960s, international trade dwindled. Since 1974, the Government of India has successfully encouraged tourism in Ladakh. As Ladakh is strategically important, the Indian military maintains a strong presence in the region.",
			"keywords": ["Aksai Chin plains", "Chinese control", "trade routes", "tourism", "Indian military"],
			"questions": ["Who claims the uninhabited Aksai Chin plains?", "What is the status of international trade in Ladakh?"]
		  }
		],
		"summary": "Ladakh is a region administered by India as a union territory and has been a disputed territory between India and Pakistan since 1947 and India and China since 1959. It is bordered by the Tibet Autonomous Region, Himachal Pradesh, Jammu and Kashmir, Gilgit-Baltistan, and Xinjiang. Ladakh extends from the Siachen Glacier to the Great Himalayas. The uninhabited Aksai Chin plains at the eastern end are claimed by the Indian Government but have been under Chinese control since 1962. Ladakh has gained strategic importance due to its location at trade routes, although international trade decreased after the closure of borders with Tibet. The Indian military maintains a strong presence in the region.",
		"title": "Ladakh: A Disputed Region Administered by India"
	  }`

	intellichunk.GenerateContainerNodes(splittedText, "Original title of ..", "https://wwwww")

}
func Test_addTestFromCLI(t *testing.T) {
	relFolderPath := "../../files"
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting project root: %v", err)
	}

	// Joining the project root with the relative folder path provided
	absFolderPath := filepath.Join(projectRoot, relFolderPath)

	intellichunk.AddFromFolder("Class_testRun2", absFolderPath, false)

}

func Test_IntelliChunkAdd(t *testing.T) {
	//longtext := `A coalition in California's Bay Area seeks to remove a multi-use path on the Richmond-San Rafael Bridge, asserting it worsens traffic and pollution. The group argues that the path diverts road space from cars, exacerbating congestion and air quality issues. Yet, studies suggest the toll plaza, not the path, is the primary cause of traffic woes. Despite the path reducing lanes, peak-hour driving times have increased only slightly. A separate study shows removing the bike lane could shift bottlenecks, increasing overall travel times.In Texas, the Great Springs Project aims to create a 100-mile trail network between San Antonio and Austin. However, 95% of Texan land is privately owned, posing challenges in securing routes and agreements with landowners. The project's timeline extends to 2036 due to these complexities.Minnesota's progressive transportation bill legalizes an "Idaho Stop" for cyclists and allocates 40% of new vehicle sales tax to transit. The bill also establishes a state-wide e-bike rebate and directs funds toward walking, biking to school, and pedestrian infrastructure.Conversely, New York faced a unique challenge with tactical urbanism: outdoor dining structures encroaching on bike lanes. An Upper West Side restaurant owner installed rubber speed bumps to slow e-bike food deliveries. However, the New York DOT demanded their removal. Traffic calming measures, such as neckdowns, were proposed to address the issue.Berlin's bike lane diversity showcases various designs influencing cyclist speed and environment awareness. Some lanes are road-grade, encouraging quicker cycling, while others at sidewalk level promote a more cautious approach. The varying infrastructure quality underscores the importance of well-designed and connected systems.In sum, these scenarios highlight the complex interactions between urban planning, infrastructure, and transportation choices, underscoring the need for data-driven decision-making and a holistic perspective on urban mobility.`

	longtext := "Ladakh is a region administered by India as a union territory[1] and constitutes an eastern portion of the larger Kashmir region that has been the subject of a dispute between India and Pakistan since 1947 and India and China since 1959.[2] Ladakh is bordered by the Tibet Autonomous Region to the east, the Indian state of Himachal Pradesh to the south, both the Indian-administered union territory of Jammu and Kashmir and the Pakistan-administered Gilgit-Baltistan to the west, and the southwest corner of Xinjiang across the Karakoram Pass in the far north. It extends from the Siachen Glacier in the Karakoram range to the north to the main Great Himalayas to the south.[11][12] The eastern end, consisting of the uninhabited Aksai Chin plains, is claimed by the Indian Government as part of Ladakh, and has been under Chinese control since 1962.[13]		In the past, Ladakh gained importance from its strategic location at the crossroads of important trade routes,[14] but as Chinese authorities closed the borders between Tibet Autonomous Region and Ladakh in the 1960s, international trade dwindled. Since 1974, the Government of India has successfully encouraged tourism in Ladakh. As Ladakh is strategically important, the Indian military maintains a strong presence in the region. "
	objIDs, err := intellichunk.Add("Class_testRun", longtext)
	if err != nil {
		log.Println("Error Test_IntelliChunkAdd:", err)
	}
	fmt.Println(objIDs)
}

func Test_AddWithoutNodes(t *testing.T) {
	// a test file with fields.
	jsonFile, err := os.Open("cache/purposes2.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	data := make(map[string][]string)

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	dataHolders := make([]models.GeneralDataHolder, 0)

	for title, contents := range data {
		holder := models.GeneralDataHolder{
			Title:     "." + title,
			Content:   strings.Join(contents, ", "), // Join all strings in the slice with a comma and a space
			Embedding: []float32{},
		}
		dataHolders = append(dataHolders, holder)
	}

	vecIDs, err := intellichunk.AddWithoutNodes("DomainList", dataHolders)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(vecIDs)
}
