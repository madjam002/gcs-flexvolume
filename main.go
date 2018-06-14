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

	dirMode, ok := options["dirMode"]
	if !ok {
		dirMode = "0755"
	}

	fileMode, ok := options["fileMode"]
	if !ok {
		fileMode = "0644"
	}

	mountPath := path.Join("/home/kubernetes/mounts/", bucket)

	if !isMountPoint(mountPath) {
		os.MkdirAll(mountPath, 0777)
		args := []string{
			"-o",
			"nonempty",
			"--dir-mode",
			dirMode,
			"--file-mode",
			fileMode,
			bucket,
			mountPath,
		}
		mountCmd := exec.Command("/home/kubernetes/bin/gcsfuse", args...)
		mountCmd.Env = append(os.Environ(), "PATH=$PATH:/home/kubernetes/bin")
		if err := mountCmd.Start(); err != nil {
			return makeResponse("Failure", err.Error())
		}
	}

	srcPath := path.Join(mountPath, subPath)

	// Create subpath if it does not exist
	intDirMode, _ := strconv.ParseUint(dirMode, 8, 32)
	os.MkdirAll(srcPath, os.FileMode(intDirMode))

	// Now we rmdir the target, and then make a symlink to it!
	if err := os.Remove(target); err != nil {
		return makeResponse("Failure", err.Error())
	}

	if err := os.Symlink(srcPath, target); err != nil {
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