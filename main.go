package main

import (
	"flag"
	"fmt"
	"os"
	v1 "patpat/api/v1"
	"patpat/config"
	"patpat/initialize"
	"patpat/judge"
	"patpat/util"
	"strconv"
	"strings"
)

func main() {
	initialize.InitMySQL()

	stuCmd := flag.NewFlagSet("stu", flag.ExitOnError)
	judgePtr := stuCmd.String("judge", "0-12345-hanhan", "Please specify the name of the folder containing Java files to judge.")
	taCmd := flag.NewFlagSet("ta", flag.ExitOnError)
	taJudgePtr := taCmd.String("judge", "0-12345-hanhan", "Please specify the name of the folder containing Java files to judge.")
	tagPtr := taCmd.String("tag", "test", "Tag for this judge.")
	regCmd := flag.NewFlagSet("reg", flag.ExitOnError)
	regSidPtr := regCmd.Int("sid", 123456, "Please specify your SID.")
	regPwdPtr := regCmd.String("pwd", "888888", "Please enter your password.")

	switch os.Args[1] {
	case "stu":
		stuCmd.Parse(os.Args[2:])
		folderName := *judgePtr
		paramList := strings.Split(folderName, "-")
		num, err := strconv.Atoi(paramList[0])
		if err != nil {
			panic("Cannot parse num!")
		}
		sid, err := strconv.Atoi(paramList[1])
		if err != nil {
			panic("Cannot parse sid!")
		}
		name := paramList[2]
		fmt.Println("Lab:", num, "SID:", sid, "Name:", name)
		tests := config.FetchJudgeConfig("test/judge.yaml")
		fmt.Println("Test cases:", tests)
		for _, t := range tests {
			testName, testData := util.FetchTestCase("test/" + t)
			fmt.Println(testName)
			testInputList, testInput, testOutputLines, testOutput, mapTable := util.ParseData(testData)
			initialize.CompileJava("javac", folderName+"/src/*.java")
			runStatus, actualOutput, actualOutputLines := initialize.RunJava(2, testInput, "java", "-classpath", folderName+"/src", "Test")
			compareResult, smallerLen, wrongOutputPos := judge.Compare(testOutputLines, actualOutputLines, mapTable)
			judge.ReportGen(t[0:len(tests[0])-5], runStatus, compareResult, smallerLen, wrongOutputPos, testInputList, testOutputLines, actualOutputLines, testOutput, actualOutput)
			judge.GradeUpload(num, sid, name, t, judge.CalcGrade(runStatus, compareResult))
		}
	case "ta":
		taCmd.Parse(os.Args[2:])
		folderName := *taJudgePtr
		paramList := strings.Split(folderName, "-")
		num, err := strconv.Atoi(paramList[0])
		if err != nil {
			panic("Cannot parse num!")
		}
		sid, err := strconv.Atoi(paramList[1])
		if err != nil {
			panic("Cannot parse sid!")
		}
		name := paramList[2]
		fmt.Println("Lab:", num, "SID:", sid, "Name:", name)
		tests := config.FetchJudgeConfig("test/judge.yaml")
		fmt.Println("Test cases:", tests)
		for _, t := range tests {
			testName, testData := util.FetchTestCase("test/" + t)
			fmt.Println(testName)
			_, testInput, testOutputLines, _, mapTable := util.ParseData(testData)
			initialize.CompileJava("javac", folderName+"/src/*.java")
			runStatus, _, actualOutputLines := initialize.RunJava(2, testInput, "java", "-classpath", folderName+"/src", "Test")
			compareResult, _, _ := judge.Compare(testOutputLines, actualOutputLines, mapTable)
			// judge.ReportGen(t[0:len(tests[0])-5], runStatus, compareResult, smallerLen, wrongOutputPos, testInputList, testOutputLines, actualOutputLines, testOutput, actualOutput)
			judge.GradeUploadFormal(num, sid, name, t, judge.CalcGrade(runStatus, compareResult), *tagPtr)
		}
	case "reg":
		regCmd.Parse(os.Args[2:])
		sid := *regSidPtr
		pwd := *regPwdPtr
		result := v1.Register(sid, pwd)
		fmt.Println(result)

	default:
		fmt.Println("Expected 'stu' subcommands!")
	}
}
