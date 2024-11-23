package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigData() map[string]string {
	file, err := os.Open("rider-config.txt")
	logError(err)
	defer file.Close()

	data := bufio.NewScanner(file)

	configData := make(map[string]string)
	for data.Scan() {
		command := data.Text()
		c := strings.SplitN(command, "=", 2) // Use 2 to split key=value correctly
		if len(c) == 2 {
			configData[c[0]] = strings.TrimSpace(c[1]) // Trim spaces around the value
		}
	}
	return configData
}

func ValidateArgs(args []string) {
	if len(args) < 1 { // Expect at least one command argument
		log.Fatal("invalid arguments. specify the command: rider run <command-name>")
	}
}

func executeCommand(ctx context.Context, command string, args []string) {
	fmt.Println("Executing command:", command, args)
	parts := strings.Fields(command) // Split command and its inline arguments
	baseCmd := parts[0]
	cmdArgs := append(parts[1:], args...) 
	cmd := exec.CommandContext(ctx, baseCmd, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Could not run command:", err)
	}
}

func getOrCreateFile(filePath string) string {
	filePath = filepath.Join(filePath, "rider-config.txt")
	if _, err := os.Stat(filePath); err != nil {
		if _, err := os.Create(filePath); err != nil {
			panic("Unknown error. run rider --help to get more idea.")
		}
	}
	return filePath
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	allArgs := os.Args[1:] // Skip the program name

	ValidateArgs(allArgs)

	baseCommand := allArgs[0]
	config := getConfigData()

	command, exists := config[baseCommand]
	if !exists {
		log.Fatal("No such command. Try adding command using rider add <command> <value>")
		return
	}

	args := allArgs[1:] // Use remaining CLI arguments as args for the command
	executeCommand(ctx, command, args)
	
	currentUser, err := user.Current()

	if err != nil {
		panic("Current User not found")
	}

	fileCreatePath := currentUser.HomeDir

	fileName := getOrCreateFile(fileCreatePath)
	fmt.Println(fileName)
}
