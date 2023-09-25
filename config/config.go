package config

import (
	"os"
	"regexp"

	"github.com/apsystole/log"
	"github.com/joho/godotenv"
)

const projectDirName = "intellichunk"

// LoadEnv is a function that loads environment variables from a .env file when the program is run locally.
// It checks if the environment variable "RUN_ENV" is set to "local". If it is not, it assumes the function is running in a Google Cloud environment
// where dynamic loading of environment variables is not allowed, and therefore the function does nothing and returns.
// If "RUN_ENV" is set to "local", it logs this information and proceeds to locate the .env file.
// It uses a regular expression to locate the project directory and then appends the .env filename to the directory path.
// Then it attempts to load the environment variables from the .env file using the godotenv package's Load function.
// If loading fails, it logs the error, along with the current working directory, and exits the program. Otherwise, it completes successfully.
func LoadEnv() error {
	if os.Getenv("RUN_ENV") != "local" {
		return nil
	} else {
		log.Info("Using RUN_ENV=local environment variables")
	}

	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Warningf("CWD: %s  err: %s", cwd, err)

		os.Exit(-1)
		return err
	}

	return nil
}
