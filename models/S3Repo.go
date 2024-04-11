package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type S3Repo struct {
	Project  string
	CreateAt string
	FileType string
	FilePath string
	Size     string
}

func GetS3Repos(awsRegion string, s3repo string) (msg string, error error) {
	if awsRegion == "" {
		return "awsRegion is required.", errors.New("GetS3Repos")
	}
	if s3repo == "" {
		return "s3repo is required.", errors.New("GetS3Repos")
	}
	out, error := GetRepos(s3repo)
	if error != nil {
		log.Printf("error: %s\n", error)
		return "bash error", error
	}
	out2 := strings.Split(out, "\n")
	var S3Arry = make([]map[string]string, 0)
	for i := range out2 {
		var S3 = make(map[string]string)
		outArry := strings.Split(out2[i], "\t")
		if len(outArry) == 3 {
			S3["updatedDt"] = outArry[0]
			S3["size"] = outArry[1]
			S3["fileName"] = outArry[2]
			idx := strings.Index(s3repo, "/")
			if idx > -1 {
				S3["filePath"] = "http://" + s3repo[0:idx] + ".s3-website." + awsRegion + ".amazonaws.com/" + S3["fileName"]
			} else {
				S3["filePath"] = "http://" + s3repo + ".s3-website." + awsRegion + ".amazonaws.com/" + S3["fileName"]
			}
			if S3["size"] != "0" {
				S3Arry = append(S3Arry, S3)
			}
		}
	}
	sort.Slice(S3Arry, func(i, j int) bool {
		return S3Arry[i]["updatedDt"] > S3Arry[j]["updatedDt"]
	})
	b, err := json.Marshal(S3Arry)
	if err != nil {
		log.Printf("error: %s\n", err)
		return "bash error", err
	}
	outStr := string(b)
	log.Println(outStr)
	return outStr, nil
}

func DeleteS3Repo(s3repo string, object string) (msg string, error error) {
	if s3repo == "" {
		return "s3repo is required.", errors.New("s3repo is required.")
	}
	var timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	idx := strings.Index(s3repo, "/")
	cmdName := "aws s3 rm s3://"
	if idx > -1 {
		cmdName += s3repo[0:idx] + "/" + object
	} else {
		cmdName += s3repo + "/" + object
	}
	shFile := fmt.Sprintf("/tmp/essh_%s.sh", timestamp)
	esFile, err := os.OpenFile(shFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Printf("esFile error: %s %s\n", esFile, err)
		return "esFile error", err
	}
	log.Printf("cmdName %s\n", cmdName)
	esFile.WriteString(cmdName)
	esFile.Close()
	doc, err := exeCmd("bash "+shFile, 10)
	if err != nil {
		log.Printf("bash error: %s\n", err)
		return "bash error", err
	}
	log.Println(doc)
	os.Remove(shFile)
	return doc, nil
}

func GetS3Repo(s3repo string, object string) (msg string, error error) {
	if s3repo == "" {
		return "s3repo is required.", errors.New("s3repo is required.")
	}
	_, err := os.Stat(object)
	if err == nil {
		return "", nil
	}
	idx := strings.Index(s3repo, "/")
	cmdName := "aws s3 cp s3://"
	if idx > -1 {
		cmdName += s3repo[0:idx] + "/" + object + " ./" + object
	} else {
		cmdName += s3repo + "/" + object + " ./" + object
	}
	doc, err := exeCmd(cmdName, 100)
	if err != nil {
		log.Printf("bash error: %s\n", err)
		return "bash error", err
	}
	log.Println(doc)
	return doc, nil
}

func UploadS3Repo(s3repo string, sourceFileName string, targetFileName string) (msg string, error error) {
	if s3repo == "" {
		return "s3repo is required.", errors.New("s3repo is required.")
	}
	_, err := os.Stat(sourceFileName)
	if err != nil {
		return "", nil
	}
	cmdName := "aws s3 cp " + sourceFileName + " s3://" + s3repo + targetFileName
	doc, err := exeCmd(cmdName, 100)
	if err != nil {
		log.Printf("bash error: %s\n", err)
		return "bash error", err
	}
	log.Println(doc)
	log.Println("========sourceFileName: " + sourceFileName)
	os.Remove(sourceFileName)
	return doc, nil
}

type ResultDoc struct {
	Raw   string `json:"raw"`
	Error string `json:"error"`
}

func GetRepos(s3repo string) (msg string, error error) {
	var timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	cmdName := "aws s3 ls " + s3repo + " --recursive | awk '{print $1\" \"$2\"\\\t\"$3\"\\\t\"$4}'"
	shFile := fmt.Sprintf("/tmp/essh_%s.sh", timestamp)
	esFile, err := os.OpenFile(shFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Printf("esFile error: %s %s\n", esFile, err)
		return "esFile error", err
	}
	log.Printf("cmdName %s\n", cmdName)
	esFile.WriteString(cmdName)
	esFile.Close()
	doc, err := exeCmd("bash "+shFile, 30)
	if err != nil {
		log.Printf("bash error: %s\n", err)
		return "bash error", err
	}
	os.Remove(shFile)
	log.Println(doc)
	return doc, nil
}

func WriteUploadFile(ifile io.Reader) (string, error) {
	data, err := ioutil.ReadAll(ifile)
	if err != nil {
		return "", err
	}
	fmt.Println("File size:", len(data))
	f, err := os.CreateTemp("", "s3file")
	_, err2 := f.Write(data)
	if err2 != nil {
		fmt.Println(err2)
		f.Close()
		return "bash error", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func RefreshCloudFront(cloudfront string, awsRegion string) (msg string, error error) {
	if cloudfront == "" {
		return "cloudfront is required.", errors.New("cloudfront is required.")
	}
	if awsRegion == "" {
		return "awsRegion is required.", errors.New("awsRegion is required.")
	}
	var timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	shFile := fmt.Sprintf("/tmp/essh_%s.sh", timestamp)
	esFile, err := os.OpenFile(shFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Printf("esFile error: %s %s\n", esFile, err)
		return "esFile error", err
	}
	cmdName := "aws cloudfront create-invalidation --distribution-id " + cloudfront + " --paths '/*' --region " + awsRegion
	esFile.WriteString(cmdName)
	esFile.Close()
	doc, err := exeCmd("bash "+shFile, 30)
	log.Println(doc)
	return doc, nil
}

func exeCmd(str string, timeout int) (string, error) {
	res := ResultDoc{}
	resultChan := make(chan string)
	errChan := make(chan error, 10)
	cmdName := ""
	var args []string
	if strings.HasPrefix(str, "bash -c ") {
		cmdName = "bash"
		args = append(args, "-c")
		args = append(args, str[len("bash -c ")+1:len(str)-1])
	} else {
		parts := strings.Fields(str)
		cmdName = parts[0]
		args = parts[1:len(parts)]
		for n := range args {
			if args[n] == "'Content-Type_application/json'" {
				args[n] = "'Content-Type: application/json'"
			} else {
				args[n] = strings.Replace(args[n], "`", " ", -1)
			}
		}
	}
	log.Println("= cmdName: ", cmdName)
	log.Println("= args: ", args)
	cmd := exec.Command(cmdName, args...)
	stdout, err := cmd.StdoutPipe()
	if wErr, ok := err.(*exec.ExitError); ok {
		if s := wErr.Error(); s != "0" {
			errChan <- err
		}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	go func() {
		stdo, err := ioutil.ReadAll(stdout)
		if err != nil {
			errChan <- err
		}
		resultChan <- string(stdo[:])
		stde, f := ioutil.ReadAll(stderr)
		if f != nil {
			log.Println(f)
			res.Error = string(stde)
		}
	}()
	err = cmd.Start()
	if err != nil {
		errChan <- err
		res.Error = fmt.Sprintf("Runner: %s", err.Error())
		if res.Error != "" {
			res.Raw = res.Error
		}
		log.Println("= res.Error2: ", res.Error)
		return res.Raw, errors.New(res.Error)
	}

loop:
	for {
		var aTimeout time.Duration = 5
		if timeout > 0 {
			aTimeout = time.Duration(timeout) * time.Second
		} else {
			aTimeout = time.Duration(1000000) * time.Hour
		}
		select {
		case <-time.After(aTimeout):
			cmd.Process.Kill()
			res.Error = "Runner: timedout"
			res.Raw = res.Error
			log.Println("------------------------------- timeout: " + string(timeout))
			log.Println("= res.Error1: ", res.Error)
			if timeout < 0 {
				os.Exit(1)
			}
			break loop
		case err := <-errChan:
			res.Error = fmt.Sprintf("Runner: %s", err.Error())
			if res.Error != "" {
				res.Raw = res.Error
			}
			log.Println("= res.Error2: ", res.Error)
			if timeout < 0 {
				os.Exit(1)
			}
			break loop
		case cmdResult := <-resultChan:
			if cmdResult != "" {
				res.Raw = cmdResult
				log.Println("= cmdResult: ", cmdResult)
				break loop
			}
		}
	}
	cmd.Wait()
	if res.Error == "" {
		res.Error = "Runner: OK"
		return res.Raw, nil
	}
	return res.Raw, errors.New(res.Error)
}
