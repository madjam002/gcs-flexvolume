package main

import (
	"os"
	"fmt"
	"encoding/json"
	"path"
	"os/exec"
	"strconv"
)

func makeResponse(status, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func isMountPoint(path string) bool {
	cmd := exec.Command("mountpoint", path)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
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
	subPath := options["subPath"]

	printDebug("target", target)
	printDebug("bucket", bucket)
	printDebug("subPath", subPath)

	dirMode, ok := options["dirMode"]
	if !ok {
		dirMode = "0755"
	}
	printDebug("dirMode", dirMode)

	fileMode, ok := options["fileMode"]
	if !ok {
		fileMode = "0644"
	}
	printDebug("fileMode", fileMode)

	mountPath := path.Join("/mnt/gcsfuse", bucket)

	if !isMountPoint(mountPath) {
		os.MkdirAll(mountPath, 0755)
		args := []string{
			"-o",
			"nonempty",
			"--dir-mode",
			dirMode,
			"--file-mode",
			fileMode,
			"--debug_fuse",
			"--debug_gcs",
			"--debug_http",
			bucket,
			mountPath,
		}
		printDebug("args", args)
		mountCmd := exec.Command("gcsfuse", args...)
		mountCmd.Start()
	}

	srcPath := path.Join(mountPath, subPath)
	printDebug("srcPath", srcPath)

	// Create subpath if it does not exist
	intDirMode, _ := strconv.ParseUint(dirMode, 8, 32)
	os.MkdirAll(srcPath, os.FileMode(intDirMode))

	// Now we rmdir the target, and then make a symlink to it!
	err := os.Remove(target)
	if err != nil {
		return makeResponse("Failure", err.Error())
	}

	err = os.Symlink(srcPath, target)
	if err != nil {
		return makeResponse("Failure", err.Error())
	}

	return makeResponse("Success", "Mount completed!")
}

func Unmount(target string) interface{} {
	err := os.Remove(target)
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

func printDebug(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a)
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
	default:
		printJSON(makeResponse("Not supported", fmt.Sprintf("Operation %s is not supported", action)))
	}
}
