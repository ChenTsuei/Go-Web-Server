package main

import (
	"flag"
	"os"
	"net"
	"strconv"
	"bufio"
	"strings"
	"io"
	"mime"
	"path"
	"os/exec"
	"log"
	"fmt"
	"io/ioutil"
)

const SERVER_STRING = "Server: gows 1.0\n"

var cgi = false
var rootPath string

func main() {
	defaultPath, _ := os.Getwd()
	rootPathFlag := flag.String("path", defaultPath, "root directory path")
	portFlag := flag.Int("port", 8080, "binding port")
	flag.Parse()
	rootPath = *rootPathFlag
	port := *portFlag
	service := ":" + strconv.Itoa(port)
	fmt.Println("Root directory is " + rootPath)
	fmt.Println("Binding port is " + service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go acceptRequest(conn)
	}
}

func acceptRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	buf, err := reader.ReadString('\n')
	checkErr(err)
	fields := strings.Fields(buf)
	if len(fields) < 3 {
		return
	}
	method := fields[0]
	fmt.Println(fields)
	if method == "POST" {
		cgi = true
	}
	url := fields[1]
	var queryString string
	if method == "GET" {
		fields := strings.Split(url, "?")
		if len(fields) > 1 {
			cgi = true
			url = fields[0]
			queryString = fields[1]
		}
	}
	filePath := rootPath + url
	if filePath[len(filePath)-1] == '/' {
		filePath += "index.html"
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			notFound(conn)
			return
		}
		checkErr(err)
	}
	if fileInfo.IsDir() {
		filePath += "/index.html"
	}
	flag := fileInfo.Mode().Perm() & os.FileMode(73)
	if uint32(flag) == uint32(73) {
		cgi = true;
	}
	if !cgi {
		serveFile(conn, filePath, method)
	} else {
		executeCgi(conn, *reader, filePath, method, queryString)
	}
}

func serveFile(conn net.Conn, path string, method string) {
	resource, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			notFound(conn)
			return
		}
		checkErr(err)
	}
	defer resource.Close()
	headers(conn, path)
	if method == "GET" {
		cat(conn, resource)
	}
}

func executeCgi(conn net.Conn, reader bufio.Reader, filepath string, method string, queryString string) {
	contentLength := -1
	if method == "POST" {
		for {
			buf, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				checkErr(err)
			}
			fields := strings.Fields(buf)
			if len(fields) > 0 && fields[0] == "Content-Length:" {
				contentLength, err = strconv.Atoi(fields[1])
				checkErr(err)
				break
			}
		}
		if contentLength == -1 {
			badRequest(conn)
		}
	}
	for {
		buf, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		checkErr(err)
		if buf == "\r\n" {
			break
		}
	}
	fmt.Println(filepath)
	proc := exec.Command(filepath)
	stdin, err := proc.StdinPipe()
	checkErr(err)
	stdout, err := proc.StdoutPipe()
	checkErr(err)
	defer stdin.Close()
	defer stdout.Close()
	conn.Write([]byte("HTTP/1.0 200 OK\n"))
	os.Setenv("REQUEST_METHOD", method)
	if method == "GET" {
		fmt.Print(queryString)
		os.Setenv("QUERY_STRING", queryString)
	} else if method == "POST" {
		os.Setenv("CONTENT_LENGTH", strconv.Itoa(contentLength))
	}
	proc.Start()
	if method == "POST" {
		input := make([]byte, contentLength)
		reader.Read(input)
		stdin.Write(input)
		stdin.Close()
	}
	buf, _ := ioutil.ReadAll(stdout)
	conn.Write(buf)
	stdout.Close()
	proc.Wait()
}

func response(conn net.Conn, status string, contentType string, content string) {
	conn.Write([]byte("HTTP/1.0 " + status + "\n"))
	conn.Write([]byte(SERVER_STRING))
	conn.Write([]byte("Content-Type: " + contentType + "\n"))
	conn.Write([]byte("\n"))
	conn.Write([]byte(content))
}

func cannotExecute(conn net.Conn) {
	response(conn, "500 Internal Server Error", "text/html",
		"<html>\n" + "<head>\n<title>500 Internal Server Error</title>\n</head>\n"+
			"<body>\n<h1>500 Internal Server Error</h1>\n<p>500 Internal Server Error.</p>\n</body>\n"+ "</html>")
}

func cat(conn net.Conn, resource *os.File) {
	buf, _ := ioutil.ReadAll(resource)
	conn.Write(buf)
}

func badRequest(conn net.Conn) {
	response(conn, "400 Bad Request", "text/html",
		"<html>\n" + "<head>\n<title>400 Bad Request</title>\n</head>\n"+
			"<body>\n"+ "<h1>400 Bad Request</h1>\n"+
			"<p>Your browser sent a bad request, such as a POST without a Content-Length.</p>\n</body>\n</html>")
}

func headers(conn net.Conn, filepath string) {
	response(conn, "200 OK", mime.TypeByExtension(path.Ext(filepath)), "")
}

func notFound(conn net.Conn) {
	response(conn, "404 Not Found", "text/html",
		"<html>\n<head>\n<title></title>\n</head>\n<body>\n<h1>404 Not Found</h1>\n"+
			"<p>Sorry, but the page you were trying to view does not exist.</p>\n</body>\n</html>")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}
