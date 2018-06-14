package main

import (
	"os"
	"fmt"
	"encoding/json"
	"os/exec"
)

func makeResponse(status, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func Init() interface{} {
	resp := makeResponse("Success", "No Initialisation required")
	resp["capabilities"] = map[string]interface{}{
		"attach": false,
	}
	return resp
}

func Mount(target string, options map[string]string) interface{} {
	bucket := options["bucket"]
	// subPath := options["subPath"]

	dirMode, ok := options["dirMode"]
	if !ok {
		dirMode = "0755"
	}

	fileMode, ok := options["fileMode"]
	if !ok {
		fileMode = "0644"
	}

	// Remove the target
	err := os.Remove(target)
	if err != nil {
		return makeResponse("Failure", err.Error())
	}

	// Use the target as the mount point for gcsfuse
	args := []string{
		"-o",
		"nonempty",
		"--dir-mode",
		dirMode,
		"--file-mode",
		fileMode,
		bucket,
		target,
	}
	mountCmd := exec.Command("gcsfuse", args...)
	err = mountCmd.Start()
	if err != nil {
		return makeResponse("Failure", err.Error())
	}

	return makeResponse("Success", "Mount completed!")
}

func Unmount(target string) interface{} {
	umountCmd := exec.Command("umount", target)
	err := umountCmd.Start()
	if err != nil {
		return makeResponse("Failure", err.Error())
	}
	return makeResponse("Success", "Successfully unmounted")
}

func printJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(jsonBytes))
}

func main() {
	switch action := os.Args[1]; action {
	case "init":
		printJSON(Init())
	case "mount":
		optsString := os.Args[3]
		opts := make(map[string]string)
		json.Unmarshal([]byte(optsString), &opts)
		printJSON(Mount(os.Args[2], opts))
	case "unmount":
		printJSON(Unmount(os.Args[2]))
	default:
		printJSON(makeResponse("Not supported", fmt.Sprintf("Operation %s is not supported", action)))
	}
}
