package main

import (
	"bufio"
	"errors"
	"fmt"
	strftime "github.com/lestrrat-go/strftime"
	sftp "github.com/pkg/sftp"
	ssh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"strings"
	"time"
)

// Colours for logging to the terminal
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// A line of the script to be run
type line_to_execute struct {
	file    string
	line    int
	command string
	args    string
}

// The global variables
var options = make(map[string]string)
var logfile *os.File
var lines = []line_to_execute{}
var current_line line_to_execute

func do_cmd(command string) {
	command = interpolate(command)
	do_log(">", "CMD "+command)

	if options["dry-run"] == "true" {
		do_log(":", "Pretend: "+command)
	} else {
		client, session, err := makeSshConnection(options["username"], options["password"], options["host"]+":"+options["port"])
		if err != nil {
			dropdead(fmt.Sprintf("%s", err))
		}

		out, err := session.CombinedOutput(command)
		prefix := "<"
		if err != nil {
			prefix = "!"
		}

		for _, v := range strings.Split(string(out), "\n") {
			if len(v) > 0 {
				do_log(prefix, v)
			}
		}

		client.Close()
	}
}

func do_copy(command string) {
	command = interpolate(command)
	parts := strings.Fields(command)
	src := parts[0]
	dst := parts[1]

	do_log(">", "COPY ["+src+"] to ["+dst+"]")

	if options["dry-run"] == "true" {
		cmd := "sftp " + src + " " + options["host"] + ":" + dst
		do_log(":", "Pretend: "+cmd)
	} else {
		// https://docs.couchdrop.io/walkthroughs/using-sftp-clients/using-sftp-with-golang

		config := &ssh.ClientConfig{
			User:            options["username"],
			Auth:            []ssh.AuthMethod{ssh.Password(options["password"])},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// Connect to the SSH server
		conn, err := ssh.Dial("tcp", options["host"]+":"+options["port"], config)
		if err != nil {
			dropdead(fmt.Sprintf("Failed to connect to SSH server: %s", err))
		}
		defer conn.Close()

		// Open SFTP session
		sftpClient, err := sftp.NewClient(conn)
		if err != nil {
			dropdead(fmt.Sprintf("Failed to open SFTP session: %s", err))
		}
		defer sftpClient.Close()

		localFile, err := os.Open(src)
		if err != nil {
			dropdead(fmt.Sprintf("Failed to open local file: %s", err))
		}
		defer localFile.Close()

		remoteFile, err := sftpClient.Create(dst)
		if err != nil {
			dropdead(fmt.Sprintf("Failed to create remote file: %s", err))
		}
		defer remoteFile.Close()

		_, err = io.Copy(remoteFile, localFile)
		if err != nil {
			dropdead(fmt.Sprintf("Failed to upload file: %s", err))
		}
	}
}

func do_log(prefix, message string) {
	ts, _ := strftime.Format("%Y-%m-%d %H:%M:%S", time.Now())
	logfile.WriteString(ts + " " + prefix + " " + message + "\n")

	color_text := message
	if prefix == ":" {
		color_text = Green + message + Reset
	} else if prefix == ">" {
		color_text = Blue + message + Reset
	} else if prefix == "!" {
		color_text = Red + message + Reset
	} else if prefix == "#" {
		color_text = Yellow + message + Reset
	} else if prefix == "?" {
		color_text = Magenta + message + Reset
	}

	fmt.Println(ts + " " + prefix + " " + color_text)
}

func do_config(k, v string) {
	k = strings.ToLower(k)
	val, ok := options[k]

	if ok {
		if v == val {
			do_log("#", fmt.Sprintf("Setting [%s] to [%s] (no change)", k, v))
		} else {
			do_log("?", fmt.Sprintf("Re-setting [%s] to [%s] from [%s]", k, v, val))
		}
	} else {
		do_log("#", fmt.Sprintf("Setting [%s] to [%s]", k, v))
	}

	options[k] = v
}

func do_debug() {
	do_log("#", "START CONFIG")
	for k, v := range options {
		do_log("#", "["+k+"] = ["+v+"]")
	}
	do_log("#", "END CONFIG")
}

func do_echo(message string) {
	do_log("#", interpolate(message))
}

func do_include(filename string) {
	if fileExists(filename) {
		dropdead(fmt.Sprintf("Include file %s not found", filename))
	} else {
		process_file(filename)
	}
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

// https://github.com/Scalingo/go-ssh-examples/blob/master/client.go
func makeSshConnection(user, pass, host string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func dropdead(message string) {
	do_log("!", message)
	os.Exit(3)
}

func check_logdir() {
	_, err := os.Stat("log")
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		err := os.Mkdir("log", 0755)
		if err != nil {
			dropdead("Unable to create the log directory")
		}
	}
}

func create_logfile() *os.File {
	check_logdir()

	f, _ := strftime.Format("%Y%m%d-%H%M", time.Now())

	logfile := "log/legion." + f + ".log"
	log, err := os.Create(logfile)

	if err != nil {
		dropdead("Unable to create the log file")
	}

	return log
}

func interpolate(message string) string {
	var t []string

	for {
		i := strings.Index(message, "{{")

		if i == -1 {
			t = append(t, message)
			break
		}

		t = append(t, message[:i])
		message = message[i+2:]

		i = strings.Index(message, "}}")
		if i == -1 {
			dropdead(fmt.Sprintf("Unbalanced template in %s at line %d: %s %s", current_line.file, current_line.line, current_line.command, current_line.args))
		}

		k := strings.ToLower(message[:i])
		val, ok := options[k]
		if ok {
			t = append(t, val)
			message = message[i+2:]
		} else {
			dropdead(fmt.Sprintf("Unable to find a substitute for %s in %s at line %d: %s %s", k, current_line.file, current_line.line, current_line.command, current_line.args))
		}
	}

	return strings.Join(t, "")
}

// https://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func process_file(filename string) {
	lines = append(lines, line_to_execute{filename, 0, "ECHO", "Reading commands from " + filename})

	line_number := 0

	readFile, err := os.Open(filename)

	if err != nil {
		dropdead(fmt.Sprintf("Unable to read %s", filename))
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line_number++

		line := fileScanner.Text()
		line = standardizeSpaces(line)

		if strings.HasPrefix(line, "#") == false && len(line) > 0 {
			parts := strings.Fields(line)
			c := strings.ToUpper(parts[0])
			a := strings.Join(parts[1:], " ")

			if c == "INCLUDE" {
				// lines = append(lines, line_to_execute{filename, line_number, "ECHO", a})
				do_include(a)
				lines = append(lines, line_to_execute{filename, line_number, "ECHO", "Resuming " + filename})
			} else {
				lines = append(lines, line_to_execute{filename, line_number, c, a})
			}
		}
	}

	readFile.Close()

	lines = append(lines, line_to_execute{filename, line_number, "ECHO", "Completed " + filename})
}

func process(scripts []string) {
	do_log("#", fmt.Sprintf("Legion command line %s", scripts))

	for _, script := range scripts {
		process_file(script)
	}

	for _, current_line = range lines {
		if current_line.command == "CMD" {
			do_cmd(current_line.args)
		} else if current_line.command == "COPY" {
			do_copy(current_line.args)
		} else if current_line.command == "CONFIG" {
			parts := strings.Fields(current_line.args)
			do_config(parts[0], parts[1])
		} else if current_line.command == "ECHO" {
			do_echo(current_line.args)
		} else if current_line.command == "DEBUG" {
			do_debug()
		} else if current_line.command == "HALT" {
			dropdead(fmt.Sprintf("%s commits suicide at line %d", current_line.file, current_line.line))
		} else if current_line.command == "INCLUDE" {
			// Nothing to do in this pass
		} else {
			dropdead(fmt.Sprintf("Unknown command [%s] at line %d of %s", current_line.command, current_line.line, current_line.file))
		}
	}

	do_log("#", "Done")
}

func opts() []string {
	options["dry-run"] = "false"
	scripts := []string{}

	for i := 0; i < len(os.Args[1:]); i++ {
		k := os.Args[(1 + i)]

		if k == "--dry-run" {
			options["dry-run"] = "true"
		} else if k == "--config" {
			i++
			v1 := os.Args[(1 + i)]
			v2 := ""
			if strings.Contains(v1, "=") {
				parts := strings.Split(v1, "=")
				v1 = parts[0]
				v2 = parts[1]
			} else {
				i++
				v2 = os.Args[(1 + i)]
			}
			do_config(v1, v2)
		} else {
			if fileExists(k) {
				dropdead(fmt.Sprintf("[%s] is not a real file", k))
			} else {
				scripts = append(scripts, k)
			}
		}
	}

	return scripts
}

func main() {
	logfile = create_logfile()
	defer logfile.Close()

	scripts := opts()

	process(scripts)
}
