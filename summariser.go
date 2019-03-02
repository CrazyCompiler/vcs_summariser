package main

import (
	"flag"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"strings"
)

type orgDetails struct {
	storyID       string
	orgName       string
	directoryName string
	commits       []string
}

func getCommitMap(commitMaps []*orgDetails, storyID string) (*orgDetails, int) {
	for index, cm := range commitMaps {
		if cm.storyID == storyID {
			return cm, index
		}
	}
	return &orgDetails{}, -1
}

func getStoryID(commitMessage string) string {
	a := strings.Split(commitMessage, "PAM")[1]
	b := strings.Split(a, ">")[0]
	if strings.Contains(b,"-") {
		return strings.Split(b, "-")[1]
	}
	return strings.TrimSpace(b)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getAllFormatedCommits(basePath, orgName, directoryName string) []*orgDetails {
	var allCommits []*orgDetails
	path := basePath + orgName + "/" + directoryName
	fmt.Println("Getting git details of the path: ", path)
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		fmt.Println("error while opening the directory: ", err)
	}
	co, err := r.CommitObjects()
	if err != nil {
		fmt.Println("error while reading commits: ", err)
	}

	//ci := commitIterator(allCommits, orgName, directoryName)

	err = co.ForEach(func(commit *object.Commit) error {
		if strings.Contains(commit.Message, "PAM") {
			storyID := getStoryID(commit.Message)
			commitMap, index := getCommitMap(allCommits, storyID)
			commitMap.storyID = storyID
			commitMap.orgName = orgName
			commitMap.directoryName = directoryName

			if !stringInSlice(commit.Message, commitMap.commits) {
				commitMap.commits = append(commitMap.commits, commit.Message)
			}

			if index >= 0 {
				allCommits[index] = commitMap
			}else {
				allCommits = append(allCommits, commitMap)
			}
		}
		//fmt.Println("Doesn't contain PAM keyword : ", commit.Message)
		return nil
	})

	if err != nil {
		fmt.Println("error while iterating over commits: ", err)
	}

	return allCommits
}

func main() {
	goPath := os.Getenv("GOPATH")
	basePath := flag.String("basePath", goPath, "Base path from where the summariser will start the search")
	orgName := flag.String("orgName", "", "Organisation name or parent directory from where the summariser will start the search in the base path, if empty will search in all organisations")
	directoryName := flag.String("directoryName", "", "specific directory name under the base path and organisation, if kept empty will search for all directory in an organizations")

	flag.Usage()

	if *orgName != "" && *directoryName != "" {
		allCommits := getAllFormatedCommits(*basePath, *orgName, *directoryName)
		for _, commit := range allCommits {
			fmt.Println("story id:", commit.storyID)
			fmt.Println("org name:", commit.orgName)
			fmt.Println("directory name:", commit.directoryName)
			fmt.Println("commits:", strings.Join(commit.commits, "\n"))
			fmt.Println("___________________________________________________n\n\n\n\n\n\n\n\n\n")
		}
	}
}