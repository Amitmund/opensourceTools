/*
Author: Amit Mund
Email: amitmund@gmail.com
Purpose: To Map between the "pid" and the container, if any.

// Few Notes:
// cat /proc/{{pid}}/cwd/var/run/docker/runtime-runc/moby/{{fullContainerID}}/state.json
// cat /proc/{{pid}}/cgroups
// cat /sys/fs/cgroup/memory/docker/{{fullContainerID}}/memory.stat
// explore stuffs within:: /sys/fs/cgroup/cpu/{{fullContainerID}}
*/

package main

// -------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------
import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------

// amund Notes
// ArgLen := len(os.Args) // This way don't work. You get following error.
// syntax error: non-declaration statement outside function body

// -------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------

// For text color
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var argLen = len(os.Args)

// -------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------

// main function.
func main() {

	// Print help(), if this program is executated without any argument.
	if argLen != 2 {
		help()
	} else {
		var _, givenPidErr = strconv.Atoi(os.Args[1])
		if givenPidErr != nil {
			// handle error
			// fmt.Println(givenPidErr)
			fmt.Println("Please provide Number as an argument.")
			os.Exit(2)
		}

		// To run the command from the given pid.
		var givenPidString = os.Args[1]
		pidMapping(givenPidString)
	}

}

// -------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------

// help function.
func help() {

	// To fetch the program name.
	programName := os.Args[0]

	fmt.Printf("\n%sUsage:%s", colorGreen, colorReset)
	fmt.Printf("\n\tThis program expect one and only one numerical pid number as an command-line argument.\n")
	fmt.Printf("\tSo that, it will try to map if the given pid# belongs to any container...\n")
	fmt.Printf("\n%sExample:%s\n", colorGreen, colorReset)
	fmt.Printf("\t%s pid_number\n", programName)
	fmt.Printf("\t%s 1234\n\n", programName)
	fmt.Printf("\n%sQuickNote:%s\n", colorGreen, colorReset)
	fmt.Printf("\t%sThis tool best goes with the following command output...%s\n", colorPurple, colorReset)
	fmt.Printf("\t%sps -eo pcpu,pid,user,size,priority,args | sort -k1 -r -n | head -10%s\n\n", colorPurple, colorReset)

}

// -------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------

//
func pidMapping(givenPid string) {

	var givenPidString = os.Args[1]
	fileName := "/" + "proc" + "/" + givenPidString + "/" + "cgroup"
	// fileName := givenPidString

	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		// log.Panicf("failed reading data from file: %s", err)
		fmt.Printf("The given pid %s\"%s\"%s doesn't belongs to any running container process.\n", colorRed, os.Args[1], colorReset)
		os.Exit(2)
	}

	fileData := string(file)

	if strings.Contains(fileData, "pids:/docker/") {

		myRegex := regexp.MustCompile("pids:/docker/(.*)")
		containerID := myRegex.FindStringSubmatch(fileData) // Returns both the regex.
		containerID_12Char := containerID[1][0:12]          // from groun 0:12 characters.

		fmt.Printf("Found the given PID belongs to: container_id: %s %v %s\n", colorBlue, containerID_12Char, colorReset)

		// Docker related command for example:
		// docker inspect --format='{{.Config.Image}}' {{container_id}}
		// docker inspect --format='{{.State.Status}}' {{container_id}}
		// docker inspect --format='{{.State.Pid}}' {{container_id}}
		// docker inspect --format='{{.State.StartedAt}}' {{container_id}}
		// docker inspect --format='{{.Name}}' {{container_id}}
		// docker inspect --format='{{.HostConfig.ShmSize}}' {{container_id}}
		// docker inspect --format='{{.Config.Labels.maintainer}}' {{container_id}}

		// Image/Application Name
		cmdOutputRaw, _ := exec.Command("docker", "inspect", "--format='{{.Config.Image}}'", containerID_12Char).Output()
		cmdOutput := string(cmdOutputRaw[:])
		fmt.Printf("\nApplicaiton Image Name: %s%s%s", colorRed, cmdOutput, colorReset)

		// Image Maintainer
		cmdOutputRaw3, _ := exec.Command("docker", "inspect", "--format='{{.Config.Labels.maintainer}}'", containerID_12Char).Output()
		ImageMaintainer := string(cmdOutputRaw3[:])
		fmt.Printf("Image Maintainer: %s%s%s", colorRed, ImageMaintainer, colorReset)

		// Stated Time
		cmdOutputRaw1, _ := exec.Command("docker", "inspect", "--format='{{.State.StartedAt}}'", containerID_12Char).Output()
		containerStartedAt := string(cmdOutputRaw1[:])
		fmt.Printf("Container Started At: %s%s%s", colorRed, containerStartedAt, colorReset)

		// Shared memory
		cmdOutputRaw2, _ := exec.Command("docker", "inspect", "--format='{{.HostConfig.ShmSize}}'", containerID_12Char).Output()
		containerSharedMemory := string(cmdOutputRaw2[:])
		fmt.Printf("Container Shared Memory in Bytes: %s%s%s\n", colorRed, containerSharedMemory, colorReset)

		// Other related commands for further details
		fmt.Printf("You might like to run %s docker container top %s %s to see further details.\n", colorRed, containerID_12Char, colorReset)
		fmt.Printf("You might like to run %s docker ps | grep %s %s to see further details.\n", colorRed, containerID_12Char, colorReset)
		fmt.Printf("You might like to run %s docker stats %s %s to see further details.\n", colorRed, containerID_12Char, colorReset)

	} else {
		fmt.Printf("The given pid %s\"%s\"%s doesn't belongs to any running container process.\n", colorRed, os.Args[1], colorReset)
		os.Exit(3)
	}

	//-----------------------------------------------------------------------
	//-----------------------------------------------------------------------

}

