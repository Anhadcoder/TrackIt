package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("TrackIT: missing command")
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit(os.Args[2:])
	case "commit":
		handleCommit(os.Args[2:])
	case "log":
		handleLog()
	case "status":
		handleStatus()
	case "revert":
		handleRevert(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func handleInit(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: trackit init <filename>")
		return
	}
	fileName := args[0]

	err := os.Mkdir(".trackit", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating .trackit directory:", err)
		return
	}

	err = os.WriteFile(".trackit/tracked_file.txt", []byte(fileName), 0644)
	if err != nil {
		fmt.Println("Error saving tracked filename:", err)
		return
	}

	os.Mkdir(".trackit/commits", 0755)

	fmt.Println("Initialized tracking for", fileName)
}

func handleCommit(args []string) {
	if len(args) < 2 || args[0] != "-m" {
		fmt.Println("Usage: trackit commit -m \"message\"")
		return
	}

	message := args[1]

	data, err := os.ReadFile(".trackit/tracked_file.txt")
	if err != nil {
		fmt.Println("Could not read tracked file name.")
		return
	}
	filename := strings.TrimSpace(string(data))

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read tracked file:", err)
		return
	}

	hash := sha1.Sum(content)
	hashStr := fmt.Sprintf("%x", hash)

	snapshotPath := ".trackit/commits/" + hashStr + ".txt"
	err = os.WriteFile(snapshotPath, content, 0644)
	if err != nil {
		fmt.Println("Error saving commit:", err)
		return
	}

	entry := fmt.Sprintf("Commit: %s\nTime: %s\nMessage: %s\n\n", hashStr, time.Now().Format(time.RFC3339), message)
	f, err := os.OpenFile(".trackit/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer f.Close()
	f.WriteString(entry)

	fmt.Println("Committed successfully:", hashStr)
}

func handleLog() {
	logData, err := os.ReadFile(".trackit/log.txt")
	if err != nil {
		fmt.Println("No commits found.")
		return
	}
	fmt.Println(string(logData))
}

func handleStatus() {
	data, err := os.ReadFile(".trackit/tracked_file.txt")
	if err != nil {
		fmt.Println("No file is currently being tracked.")
		return
	}
	filename := strings.TrimSpace(string(data))

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read tracked file:", err)
		return
	}

	hash := sha1.Sum(content)
	hashStr := fmt.Sprintf("%x", hash)

	files, err := ioutil.ReadDir(".trackit/commits")
	if err != nil {
		fmt.Println("Error reading commits:", err)
		return
	}

	found := false
	for _, file := range files {
		if file.Name() == hashStr+".txt" {
			found = true
			break
		}
	}

	if found {
		fmt.Println("Status: No changes since last commit.")
	} else {
		fmt.Println("Status: Changes detected since last commit.")
	}
}

func handleRevert(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: trackit revert <commit_hash>")
		return
	}
	commitHash := args[0]

	snapshotPath := ".trackit/commits/" + commitHash + ".txt"
	content, err := os.ReadFile(snapshotPath)
	if err != nil {
		fmt.Println("Commit not found:", commitHash)
		return
	}

	data, err := os.ReadFile(".trackit/tracked_file.txt")
	if err != nil {
		fmt.Println("Could not read tracked file name.")
		return
	}
	filename := strings.TrimSpace(string(data))

	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Println("Error reverting file:", err)
		return
	}

	fmt.Println("Reverted", filename, "to commit", commitHash)
}
