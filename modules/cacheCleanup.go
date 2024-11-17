package cache

import (
        "fmt"
        "os"
        "log"
        "time"

)

func CacheCleanup(cacheValidTime,cacheDir) {

        // Get the current time
        now := time.Now()

        // Read the contents of the directory
        files, err := os.ReadDir(cacheDir)
        if err != nil {
                log.Fatal(err)
        }

        // Print the list of files
        for _, file := range files {
                filePath := cacheDir + "/" + file.Name()

                // Get file information
                fileInfo, err := os.Stat(filePath)
                if err != nil {
                        log.Fatal(err)
                }

                // Check if the file was modifed more than one minute
                if now.Sub(fileInfo.ModTime()) > cacheValidTime * time.Minute {
                    os.Remove(filePath)
                    fmt.Println("file: ", file.Name(), " ***DELETED*** ")
                }
        }
}
