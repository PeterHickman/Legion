package main

// TODO: need to use the flag package

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	ac "github.com/PeterHickman/ansi_colours"
	ep "github.com/PeterHickman/expand_path"
	toolbox "github.com/PeterHickman/toolbox"
	"github.com/lestrrat-go/strftime"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const TRUE = "true"

type lineToExecute struct {
	file    string
	line    int
	command string
	args    string
}

var options = make(map[string]string)
var logfile *os.File
var lines = []lineToExecute{}
var currentLine lineToExecute

const defaultPrivateKey string = "~/.ssh/id_rsa"

// https://blog.ralch.com/articles/golang-ssh-connection/
func publicKeyFile(file string) ssh.AuthMethod {
	path, _ := ep.ExpandPath(file)

	if !toolbox.FileExists(path) {
		doLog("?", fmt.Sprintf("Unable to read %s using username and password", file))
		return nil
	}

	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

func makeSSHConfig() *ssh.ClientConfig {
	auths := []ssh.AuthMethod{ssh.Password(options["password"])}

	p := publicKeyFile(defaultPrivateKey)
	if p != nil {
		auths = append(auths, p)
	}

	sshConfig := &ssh.ClientConfig{
		User:            options["username"],
		Auth:            auths,
		HostKeyCallback: ssh.HostKeyCallback(func(string, net.Addr, ssh.PublicKey) error { return nil }),
	}
	sshConfig.SetDefaults()

	return sshConfig
}

// https://github.com/Scalingo/go-ssh-examples/blob/master/client.go
func makeSSHConnection() (*ssh.Client, *ssh.Session, error) {
	sshConfig := makeSSHConfig()

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", options["host"], options["port"]), sshConfig)
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

func doCmd(command string) {
	command = interpolate(command)
	doLog(">", "CMD "+command)

	if options["dry-run"] == TRUE {
		doLog(":", "Pretend: "+command)
	} else {
		client, session, err := makeSSHConnection()
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
				doLog(prefix, v)
			}
		}

		client.Close()
	}
}

func doCopy(command string) {
	command = interpolate(command)
	parts := strings.Fields(command)
	src := parts[0]
	dst := parts[1]

	doLog(">", "COPY ["+src+"] to ["+dst+"]")

	if options["dry-run"] == TRUE {
		cmd := "sftp " + src + " " + options["host"] + ":" + dst
		doLog(":", "Pretend: "+cmd)
	} else {
		// https://docs.couchdrop.io/walkthroughs/using-sftp-clients/using-sftp-with-golang

		sshConfig := makeSSHConfig()

		// Connect to the SSH server
		conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", options["host"], options["port"]), sshConfig)
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

func doLog(prefix, message string) {
	ts, _ := strftime.Format("%Y-%m-%d %H:%M:%S", time.Now())
	_, _ = logfile.WriteString(ts + " " + prefix + " " + message + "\n")

	colorText := message
	switch prefix {
	case ":":
		colorText = ac.Green(message)
	case ">":
		colorText = ac.Blue(message)
	case "!":
		colorText = ac.Red(message)
	case "#":
		colorText = ac.Yellow(message)
	case "?":
		colorText = ac.Magenta(message)
	}

	fmt.Println(ts + " " + prefix + " " + colorText)
}

func dropdead(message string) {
	doLog("!", message)
	os.Exit(3)
}

func checkLogdir() {
	if toolbox.FileExists("log") {
		return
	}

	err := os.Mkdir("log", 0755)
	if err != nil {
		dropdead("Unable to create the log directory")
	}
}

func createLogfile() *os.File {
	checkLogdir()

	f, _ := strftime.Format("%Y%m%d-%H%M", time.Now())

	logfile := "log/legion." + f + ".log"
	log, err := os.Create(logfile)

	if err != nil {
		dropdead("Unable to create the log file")
	}

	return log
}

func doConfig(k, v string) {
	k = strings.ToLower(k)
	val, ok := options[k]

	if ok {
		if v == val {
			doLog("#", fmt.Sprintf("Setting [%s] to [%s] (no change)", k, v))
		} else {
			doLog("?", fmt.Sprintf("Re-setting [%s] to [%s] from [%s]", k, v, val))
		}
	} else {
		doLog("#", fmt.Sprintf("Setting [%s] to [%s]", k, v))
	}

	options[k] = v
}

func doDebug() {
	doLog("#", "START CONFIG")
	for k, v := range options {
		doLog("#", "["+k+"] = ["+v+"]")
	}
	doLog("#", "END CONFIG")
}

func doEcho(message string) {
	doLog("#", interpolate(message))
}

func doInclude(filename string) {
	if toolbox.FileExists(filename) {
		processFile(filename)
	} else {
		dropdead(fmt.Sprintf("Include file %s not found", filename))
	}
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
			dropdead(fmt.Sprintf("Unbalanced template in %s at line %d: %s %s", currentLine.file, currentLine.line, currentLine.command, currentLine.args))
		}

		k := strings.ToLower(message[:i])
		val, ok := options[k]
		if ok {
			t = append(t, val)
			message = message[i+2:]
		} else {
			dropdead(fmt.Sprintf("Unable to find a substitute for %s in %s at line %d: %s %s", k, currentLine.file, currentLine.line, currentLine.command, currentLine.args))
		}
	}

	return strings.Join(t, "")
}

// https://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func processFile(filename string) {
	lines = append(lines, lineToExecute{filename, 0, "ECHO", "Reading commands from " + filename})

	lineNumber := 0

	readFile, err := os.Open(filename)

	if err != nil {
		dropdead("Unable to read " + filename)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		lineNumber++

		line := fileScanner.Text()
		line = standardizeSpaces(line)

		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			parts := strings.Fields(line)
			c := strings.ToUpper(parts[0])
			a := strings.Join(parts[1:], " ")

			if c == "INCLUDE" {
				// lines = append(lines, lineToExecute{filename, lineNumber, "ECHO", a})
				doInclude(a)
				lines = append(lines, lineToExecute{filename, lineNumber, "ECHO", "Resuming " + filename})
			} else {
				lines = append(lines, lineToExecute{filename, lineNumber, c, a})
			}
		}
	}

	readFile.Close()

	lines = append(lines, lineToExecute{filename, lineNumber, "ECHO", "Completed " + filename})
}

func process(scripts []string) {
	doLog("#", fmt.Sprintf("Legion command line %s", scripts))

	for _, script := range scripts {
		processFile(script)
	}

	for _, currentLine = range lines {
		switch currentLine.command {
		case "CMD":
			doCmd(currentLine.args)
		case "COPY":
			doCopy(currentLine.args)
		case "CONFIG":
			parts := strings.Fields(currentLine.args)
			doConfig(parts[0], parts[1])
		case "ECHO":
			doEcho(currentLine.args)
		case "DEBUG":
			doDebug()
		case "HALT":
			dropdead(fmt.Sprintf("%s commits suicide at line %d", currentLine.file, currentLine.line))
		case "INCLUDE":
			// Nothing to do in this pass
		default:
			dropdead(fmt.Sprintf("Unknown command [%s] at line %d of %s", currentLine.command, currentLine.line, currentLine.file))
		}
	}

	doLog("#", "Done")
}

func opts() []string {
	options["dry-run"] = "false"
	scripts := []string{}

	for i := 0; i < len(os.Args[1:]); i++ {
		k := os.Args[(1 + i)]

		switch k {
		case "--dry-run":
			options["dry-run"] = TRUE
		case "--config":
			i++
			v1 := os.Args[(1 + i)]
			var v2 string
			if strings.Contains(v1, "=") {
				parts := strings.Split(v1, "=")
				v1 = parts[0]
				v2 = parts[1]
			} else {
				i++
				v2 = os.Args[(1 + i)]
			}
			doConfig(v1, v2)
		default:
			if toolbox.FileExists(k) {
				scripts = append(scripts, k)
			} else {
				dropdead(fmt.Sprintf("[%s] is not a real file", k))
			}
		}
	}

	return scripts
}

func main() {
	logfile = createLogfile()
	defer logfile.Close()

	scripts := opts()

	process(scripts)
}
