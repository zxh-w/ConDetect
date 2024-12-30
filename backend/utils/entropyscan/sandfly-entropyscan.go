package entropyscan

/*
This utility will help find packed or encrypted files or processes on a Linux system by calculating the entropy
to see how random they are. Packed or encrypted malware often appears to be a very random executable file and this
utility can help identify potential intrusions.

You can calculate entropy on all files, or limit the search just to Linux ELF executables that have an entropy of
your threshold. Linux processes can be scanned as well automatically.

Sandfly Security produces an agentless endpoint detection and incident response platform (EDR) for Linux. You can
find out more about how it works at: https://www.sandflysecurity.com

MIT License

Copyright (c) 2019-2022 Sandfly Security Ltd.
https://www.sandflysecurity.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of
the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Version: 1.1.1
Author: @SandflySecurity
*/

import (
	"ConDetect/backend/utils/entropyscan/fileutils"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	// constVersion Version
	constVersion = "1.1.1"
	// constProcDir default /proc dir for processes.
	constProcDir = "/proc"
	// constDelimeterDefault default delimiter for CSV output.
	constDelimeterDefault = ","
	// constMinPID minimum PID value allowed for process checks.
	constMinPID = 1
	// constMaxPID maximum PID value allowed for process checks. 64bit linux is 2^22. This value is a limiter.
	constMaxPID = 4194304
)

type FileData struct {
	Path    string  `json:"path"`
	Name    string  `json:"name"`
	Entropy float64 `json:"entropy"`
	Elf     bool    `json:"elf"`
	Hash    hashes  `json:"hash"`
}

type hashes struct {
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
	Sha512 string `json:"sha512"`
}

func AnalyzeEntropy(filePath, dirPath string, entropyMaxVal float64, elfOnly, procOnly bool) ([]FileData, error) {
	var results []FileData
	if entropyMaxVal > 8 {
		return results, fmt.Errorf("max entropy value is 8.0")
	}
	if entropyMaxVal < 0 {
		return results, fmt.Errorf("min entropy value is 0.0")
	}

	if procOnly {
		if os.Geteuid() != 0 {
			return results, fmt.Errorf("process checking option requires UID/EUID 0 (root) to run")
		}
		pidPaths, err := genPIDExePaths()
		if err != nil {
			return results, fmt.Errorf("error generating PID list: %v", err)
		}
		for _, path := range pidPaths {
			fileInfo, err := checkFilePath(path, true, entropyMaxVal)
			if err == nil && fileInfo.Entropy >= entropyMaxVal {
				results = append(results, fileInfo)
			}
		}
		return results, nil
	}

	if filePath != "" {
		fileInfo, err := checkFilePath(filePath, elfOnly, entropyMaxVal)
		if err != nil {
			return results, fmt.Errorf("error processing file (%s): %v", filePath, err)
		}
		if fileInfo.Entropy >= entropyMaxVal {
			results = append(results, fileInfo)
		}
		return results, nil
	}

	if dirPath != "" {
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking directory (%s) inside search function: %v", path, err)
			}
			if info != nil && !info.IsDir() && info.Mode().IsRegular() {
				fileInfo, err := checkFilePath(path, elfOnly, entropyMaxVal)
				if err != nil {
					return fmt.Errorf("error processing file (%s): %v", path, err)
				}
				if fileInfo.Entropy >= entropyMaxVal {
					// printResults(fileInfo, csvOutput, delimChar)
					results = append(results, fileInfo)
				}
			}
			return nil
		})
		if err != nil {
			return results, fmt.Errorf("error walking directory (%s): %v", dirPath, err)
		}
		return results, nil
	}
	return results, nil
}

// Prints results
func printResults(fileInfo FileData, csvFormat bool, delimChar string) {

	if !csvFormat {
		fmt.Printf("filename: %s\npath: %s\nentropy: %.2f\nelf: %v\nmd5: %s\nsha1: %s\nsha256: %s\nsha512: %s\n\n",
			fileInfo.Name,
			fileInfo.Path,
			fileInfo.Entropy,
			fileInfo.Elf,
			fileInfo.Hash.Md5,
			fileInfo.Hash.Sha256,
			fileInfo.Hash.Sha1,
			fileInfo.Hash.Sha512)
	} else {
		fmt.Printf("%s%s%s%s%.2f%s%v%s%s%s%s%s%s%s%s\n",
			fileInfo.Name,
			delimChar,
			fileInfo.Path,
			delimChar,
			fileInfo.Entropy,
			delimChar,
			fileInfo.Elf,
			delimChar,
			fileInfo.Hash.Md5,
			delimChar,
			fileInfo.Hash.Sha1,
			delimChar,
			fileInfo.Hash.Sha256,
			delimChar,
			fileInfo.Hash.Sha512)
	}
}

func checkFilePath(filePath string, elfOnly bool, entropyMaxVal float64) (fileInfo FileData, err error) {
	isElfType, err := fileutils.IsElfType(filePath)
	if err != nil {
		return fileInfo, err
	}
	_, fileName := filepath.Split(filePath)

	fileInfo.Path = filePath
	fileInfo.Name = fileName
	fileInfo.Elf = isElfType
	fileInfo.Entropy = -1

	// If they only want Linux ELFs.
	if elfOnly && isElfType {
		entropy, err := fileutils.Entropy(filePath)
		if err != nil {
			log.Fatalf("error calculating entropy for file (%s): %v\n", filePath, err)
		}
		fileInfo.Entropy = entropy
	}
	// They want entropy on all files.
	if !elfOnly {
		entropy, err := fileutils.Entropy(filePath)
		if err != nil {
			log.Fatalf("error calculating entropy for file (%s): %v\n", filePath, err)
		}
		fileInfo.Entropy = entropy
	}

	if fileInfo.Entropy >= entropyMaxVal {
		md5, err := fileutils.HashMD5(filePath)
		if err != nil {
			log.Fatalf("error calculating MD5 hash for file (%s): %v\n", filePath, err)
		}
		sha1, err := fileutils.HashSHA1(filePath)
		if err != nil {
			log.Fatalf("error calculating SHA1 hash for file (%s): %v\n", filePath, err)
		}
		sha256, err := fileutils.HashSHA256(filePath)
		if err != nil {
			log.Fatalf("error calculating SHA256 hash for file (%s): %v\n", filePath, err)
		}
		sha512, err := fileutils.HashSHA512(filePath)
		if err != nil {
			log.Fatalf("error calculating SHA512 hash for file (%s): %v\n", filePath, err)
		}
		fileInfo.Hash.Md5 = md5
		fileInfo.Hash.Sha1 = sha1
		fileInfo.Hash.Sha256 = sha256
		fileInfo.Hash.Sha512 = sha512
	}

	return fileInfo, nil
}

func genPIDExePaths() (pidPaths []string, err error) {

	for pid := constMinPID; pid < constMaxPID; pid++ {
		pidPaths = append(pidPaths, path.Join(constProcDir, strconv.Itoa(pid), "/exe"))
	}

	return pidPaths, nil
}
