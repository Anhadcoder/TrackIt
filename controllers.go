package controllers

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/sergi/go-diff/diffmatchpatch"
)

const (
    gitgoDir   = ".gitGo"
    objectsDir = ".gitGo/objects"
    refsDir    = ".gitGo/refs"
    headFile   = ".gitGo/HEAD.gotem"
    masterRef  = ".gitGo/refs/master"
)

// GitGoInit creates .gitGo folder and initial files
func GitGoInit(repoPath string, filename string) {
    gitPath := filepath.Join(repoPath, gitgoDir)

    if _, err := os.Stat(gitPath); !os.IsNotExist(err) {
        fmt.Println("GitGo repository already initialized")
        return
    }

    err := os.MkdirAll(gitPath, 0755)
    if err != nil {
        fmt.Println("Error creating .gitGo directory:", err)
        return
    }

    // Create objects and refs folders
    os.MkdirAll(filepath.Join(gitPath, "objects"), 0755)
    os.MkdirAll(filepath.Join(gitPath, "refs"), 0755)

    // Create initial HEAD file pointing to master branch
    err = ioutil.WriteFile(headFile, []byte("refs/master"), 0644)
    if err != nil {
        fmt.Println("Error writing HEAD file:", err)
        return
    }

    // Create empty master ref file
    masterRefPath := filepath.Join(repoPath, masterRef)
    ioutil.WriteFile(masterRefPath, []byte(""), 0644)

    fmt.Println("GitGo repository initialized, tracking file:", filename)
}

// hashObject reads a file and returns its SHA1 hash and writes it into objects dir
func hashObject(filePath string) (string, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    h := sha1.New()
    h.Write(data)
    sha := hex.EncodeToString(h.Sum(nil))

    // Write object content under objects/sha prefix (first 2 chars as folder)
    objDir := filepath.Join(objectsDir, sha[:2])
    objFile := filepath.Join(objDir, sha[2:])
    if _, err := os.Stat(objFile); os.IsNotExist(err) {
        os.MkdirAll(objDir, 0755)
        err = ioutil.WriteFile(objFile, data, 0644)
        if err != nil {
            return "", err
        }
    }
    return sha, nil
}

// readRef returns the commit SHA that HEAD points to
func readRef() (string, error) {
    headRef, err := ioutil.ReadFile(headFile)
    if err != nil {
        return "", err
    }
    refPath := strings.TrimSpace(string(headRef))
    refContent, err := ioutil.ReadFile(refPath)
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(refContent)), nil
}

// updateRef updates the ref file with the new commit SHA
func updateRef(newSHA string) error {
    headRef, err := ioutil.ReadFile(headFile)
    if err != nil {
        return err
    }
    refPath := strings.TrimSpace(string(headRef))
    return ioutil.WriteFile(refPath, []byte(newSHA), 0644)
}

// GitGoCommit creates a commit object and updates refs
func GitGoCommit(message string) {
    // For simplicity, commit the file specified in init (assume single file 'file.txt')
    fileToCommit := "file.txt"
    sha, err := hashObject(fileToCommit)
    if err != nil {
        fmt.Println("Error hashing object:", err)
        return
    }

    parent, _ := readRef() // if no previous commit, parent is empty string

    // Compose commit content
    commitContent := fmt.Sprintf("tree %s\n", sha)
    if parent != "" {
        commitContent += fmt.Sprintf("parent %s\n", parent)
    }
    commitContent += fmt.Sprintf("date %s\n", time.Now().Format(time.RFC3339))
    commitContent += fmt.Sprintf("message %s\n", message)

    // Hash commit object
    h := sha1.New()
    h.Write([]byte(commitContent))
    commitSHA := hex.EncodeToString(h.Sum(nil))

    // Write commit object to objects dir
    objDir := filepath.Join(objectsDir, commitSHA[:2])
    objFile := filepath.Join(objDir, commitSHA[2:])
    if _, err := os.Stat(objFile); os.IsNotExist(err) {
        os.MkdirAll(objDir, 0755)
        ioutil.WriteFile(objFile, []byte(commitContent), 0644)
    }

    // Update refs/master with new commit SHA
    err = updateRef(commitSHA)
    if err != nil {
        fmt.Println("Error updating ref:", err)
        return
    }

    fmt.Println("Committed with SHA:", commitSHA)
}

// GitGoLog_noarg prints commit history starting from HEAD
func GitGoLog_noarg() {
    currentSHA, err := readRef()
    if err != nil || currentSHA == "" {
        fmt.Println("No commits yet")
        return
    }

    for currentSHA != "" {
        objDir := filepath.Join(objectsDir, currentSHA[:2])
        objFile := filepath.Join(objDir, currentSHA[2:])
        content, err := ioutil.ReadFile(objFile)
        if err != nil {
            fmt.Println("Error reading commit object:", err)
            return
        }
        fmt.Println("Commit:", currentSHA)
        fmt.Println(string(content))
        fmt.Println("------")

        // parse parent commit SHA
        lines := strings.Split(string(content), "\n")
        parentSHA := ""
        for _, line := range lines {
            if strings.HasPrefix(line, "parent ") {
                parentSHA = strings.TrimSpace(line[len("parent "):])
                break
            }
        }
        currentSHA = parentSHA
    }
}

// GitGoStatus shows current commit and tracked file status
func GitGoStatus() {
    sha, err := readRef()
    if err != nil {
        fmt.Println("Error reading HEAD:", err)
        return
    }
    fmt.Println("On commit:", sha)
}

// GitGoDiff shows the diff between two files using sergi/go-diff
func GitGoDiff(file1, file2 string) {
    dmp := diffmatchpatch.New()
    f1, err1 := ioutil.ReadFile(file1)
    f2, err2 := ioutil.ReadFile(file2)
    if err1 != nil || err2 != nil {
        fmt.Println("Error reading files")
        return
    }
    diffs := dmp.DiffMain(string(f1), string(f2), false)
    fmt.Println(dmp.DiffPrettyText(diffs))
}
