package main

import (
	"fmt"
	"os"

	"github.com/mohanrajaaa/go_crash_course/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <action> [filename] [content]")
		os.Exit(1)
	}

	action := os.Args[1]
	filename := ""
	content := ""

	if len(os.Args) > 2 {
		filename = os.Args[2]
	}

	if len(os.Args) > 3 {
		content = os.Args[3]
	}

	switch action {
	case "create":
		err := cli.CreateFile(filename, content)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		fmt.Println("File created successfully.")

	case "read":
		fileContent, err := cli.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		fmt.Println("File content:", fileContent)

	case "update":
		err := cli.UpdateFile(filename, content)
		if err != nil {
			fmt.Println("Error updating file:", err)
			return
		}
		fmt.Println("File updated successfully.")

	case "delete":
		err := cli.DeleteFile(filename)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
		fmt.Println("File deleted successfully.")

	default:
		fmt.Println("Invalid action. Please use one of: create, read, update, delete.")
		os.Exit(1)
	}
}
