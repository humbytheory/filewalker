package main
/*
     Copyright (C) 2014  Humberto Castro

     This program is free software: you can redistribute it and/or modify
     it under the terms of the GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the GNU General Public License
     along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
    "strings"
    "log"
    "syscall"
    "encoding/json"
    "fmt"
    "os"
    "time"
    "flag"
    "path/filepath"
)

type FileDetailsJSON struct {
    ScanDate    string            `json:"ScanDate,omitempty"`
    Error       error             `json:"Error,omitempty"`
    Type        string            `json:"Type,omitempty"`
    Path        string            `json:"Path,omitempty"`
    Name        string            `json:"Name,omitempty"`
    Extension   string            `json:"Extension,omitempty"`
    Size        int64             `json:"Size,omitempty"`
    Uid         uint32            `json:"Uid,omitempty"`
    Gid         uint32            `json:"Gid,omitempty"`
    Inode       uint64            `json:"Inode,omitempty"`
    Mtime       string            `json:"Mtime,omitempty"`
    Atime       string            `json:"Atime,omitempty"`
    Ctime       string            `json:"Ctime,omitempty"`
    KV          map[string]string `json:"KV,omitempty"`
}

var (
    extraKeyValues  = flag.String("k", "", "Optional: Space delimited key values")
    pretty          = flag.Bool("pretty", false, "Optional: Pretty print JSON")
    excludeScanDate = flag.Bool("t", false, "Optional: Exclude date/time file was scanned")
    excludeDir      = flag.String("excludeDir",  "", "Optional: Directory patterns to exclude. Can be colon delimited for multiple excludes")
    excludeFile     = flag.String("excludeFile", "", "Optional: Filename patterns to exclude. Can be colon delimited for multiple excludes")
    excludeSize     = flag.Bool("s", false, "Optional: Exclude file size")
    excludeExt      = flag.Bool("e", false, "Optional: Exclude file extension")
    excludeUid      = flag.Bool("u", false, "Optional: Exclude uid")
    excludeGid      = flag.Bool("g", false, "Optional: Exclude gid")
    excludeInode    = flag.Bool("i", false, "Optional: Exclude inode")
    excludeMtime    = flag.Bool("m", false, "Optional: Exclude mtime")
    excludeAtime    = flag.Bool("a", false, "Optional: Exclude atime")
    excludeCtime    = flag.Bool("c", false, "Optional: Exclude ctime")
    showVersion     = flag.Bool("v", false, "Optional: Display version detail")
    includeErrors   = flag.Bool("errors", false, "Optional: Include errors")
    timeFormat      = flag.String("time", "RFC3339", "Optional: time format")
    eKV  = make(map[string]string)
    ExcludePatternsDir []string
    ExcludePatternsFile []string
)

var Usage = func() {
    fmt.Fprintf(os.Stderr, "\nUsage: %s [ OPTION ]...  DIRECTORY\n", os.Args[0])
    fmt.Printf("\n%-17s\n%17s %s\n",
        "REQUIRED",
        "DIRECTORY",        "Target directory to scan")
    fmt.Printf("\n%-17s\n%17s %s\n%17s %s\n%17s %s\n",
        "OPTIONS",
        "-h",               "Show this usage help",
        "-v",               "Show version and license",
        "-pretty",          "Pretty print JSON")
    fmt.Printf("\n%-17s %s\n%17s %s\n%17s %s\n\n%17s %s\n%17s %s\n",
        "Exclusions",       "Can be colon delimited for multiple excludes",
        "-excludeDir",      "Directory patterns to exclude",
        "-excludeFile",     "File patterns to exclude",
        " ",                "EXAMPLE Single:       -excludeDir \"pattern1\"",
        " ",                "EXAMPLE Multiple:     -excludeFile \"pattern1:pattern2:pattern3\"")
    fmt.Printf("\n%-17s\n%17s %s\n%17s %s\n%17s %s\n%17s %s\n%17s %s\n%17s %s\n%17s %s\n%17s %s\n%17s %s\n",
        "Excluding fields from report",
        "-t",               "Excludes date/time file was scanned",
        "-s",               "Excludes file size",
        "-e",               "Excludes file extension",
        "-u",               "Excludes uid",
        "-g",               "Excludes gid",
        "-i",               "Excludes inode",
        "-m",               "Excludes mtime",
        "-a",               "Excludes atime",
        "-c",               "Excludes ctime")
    fmt.Printf("\n%-17s\n%17s %s\n\n%17s %s\n%17s %s\n",
        "Adding extra data to report",
        "-k",               "Space delimited key values.",
        " ",                "EXAMPLE Single:       -k \"key1:value1\"",
        " ",                "EXAMPLE Multiple:     -k \"key1:value1 key2:value2 key3:value3\"")
    fmt.Printf("\n%-17s\n%17s %s\n",
        "Include errors in report",
        "-errors",          "Include errors in report")
    fmt.Printf("\n%-17s\n%17s %s\n%17s %s\n\n%17s %s\n",
        "Change time formating",
        "-time",            "Valid time formats: ANSIC, UnixDate, RubyDate, RFC822, RFC822Z, RFC850, RFC1123,",
        " ",                "RFC1123Z, RFC3339, RFC3339Nano, Kitchen, Stamp, StampMilli, StampMicro, StampNano",
        " ",                "EXAMPLE:             -time RFC3339")
    os.Exit(2)
}

func main() {
    flag.Usage = Usage
    flag.Parse()
    args := flag.Args()

    if *showVersion {
        fmt.Printf("Version %s\n\n%s\n%s\n%s\n%s\n%s\n",
                   "1",
                   "filewalker  Copyright (C) 2014  Humberto Castro",
                   "This program comes with ABSOLUTELY NO WARRANTY.",
                   "This is free software, and you are welcome to redistribute it",
                   "under certain conditions; see the following URL for details:",
                   "<https://github.com/humbytheory/filewalker>")

        os.Exit(0)
    }
    if len(args) > 1 {
        log.Print("Too many arguements given.\n\n")
        Usage()
    }
    if len(args) == 0 {
        log.Print("Missing directory to scan.\n\n")
        Usage()
    }
    targetPath := args[0]

    keyValues := strings.Fields(*extraKeyValues)
    for _,element := range keyValues{
        keyValues := strings.Split(element, ":")
        if len(keyValues) != 2 || len(keyValues[0]) == 0 || len(keyValues[1]) == 0{
            log.Fatal("Error in key/pair:   \""+element+"\"")
        }
        k := keyValues[0]
        v := keyValues[1]
        eKV[ k ] = v
    }

    selectTimeFormat(*timeFormat)
    _ = filepath.Walk(targetPath, workTargetPathsJSON)
}

func workTargetPathsJSON(targetPath string, fi os.FileInfo, err error) error {
    var m *FileDetailsJSON
    m = &FileDetailsJSON{}
    if !*excludeScanDate{
        tmp := time.Now()
        m.ScanDate = tmp.Format(*timeFormat)
    }
    if err != nil{
        log.Println(err)
        m.Error = err
        if *includeErrors {
            printReport(m)
        }
        return nil
    }
    m.Path, m.Name = filepath.Split(targetPath)
    if fi.IsDir() {
        if checkExclude(m.Name,*excludeDir){
            log.Print("skipping excluded directory: ",targetPath)
            return filepath.SkipDir
        }
        m.Type = "d"
    } else {
        if checkExclude(m.Name,*excludeFile){
            log.Print("skipping excluded file: ",targetPath)
            return nil
        }
        m.Type = "f"
    }
    fetchFileDetailsJSON(m, fi)
    if !*excludeExt && m.Type == "f" {
        m.Extension = filepath.Ext( m.Name )
    }
    if len(eKV) > 0{
        m.KV = eKV
    }
    printReport(m)
    return nil
}
func printReport( m *FileDetailsJSON ){
    if *pretty{
        b, _ := json.MarshalIndent(m, "", "    ")
        fmt.Printf("%s\n",b)
    } else {
        b, _ := json.Marshal(m)
        fmt.Printf("%s\n",b)
    }
}
func fetchFileDetailsJSON(m *FileDetailsJSON, fi os.FileInfo) {
    if !*excludeSize{
        m.Size = fi.Size()
    }
    if !*excludeMtime{
        tmp := fi.ModTime()
        m.Mtime = tmp.Format(*timeFormat)
    }
    stat := fi.Sys().(*syscall.Stat_t)
    if !*excludeUid{
        m.Uid = stat.Uid
    }
    if !*excludeGid{
        m.Gid = stat.Gid
    }
    if !*excludeInode{
        m.Inode = stat.Ino
    }
    if !*excludeAtime{
        tmp := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
        m.Atime = tmp.Format(*timeFormat)
    }
    if !*excludeCtime{
        tmp := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
        m.Ctime = tmp.Format(*timeFormat)
    }
    /*
    if 1==0 {
        log.Print("Number of links: ",stat.Nlink)
    }
    */
}
func checkExclude(target, excludePatterns string)(retval bool){
    patterns  := filepath.SplitList(excludePatterns)
    for _, excludePattern := range patterns {
        //log.Printf("%10s %-20s %10s %-20s\n","target:",target,"excludePattern:",excludePattern)
        match, err := filepath.Match(excludePattern, target)
        if err == nil && match {
            return true
        }
    }
    return false
}
func selectTimeFormat ( requestedFormat string){
        switch requestedFormat {
        case "ANSIC":
                *timeFormat = time.ANSIC
        case "UnixDate":
                *timeFormat = time.UnixDate
        case "RubyDate":
                *timeFormat = time.RubyDate
        case "RFC822":
                *timeFormat = time.RFC822
        case "RFC822Z":
                *timeFormat = time.RFC822Z
        case "RFC850":
                *timeFormat = time.RFC850
        case "RFC1123":
                *timeFormat = time.RFC1123
        case "RFC1123Z":
                *timeFormat = time.RFC1123Z
        case "RFC3339":
                *timeFormat = time.RFC3339
        case "RFC3339Nano":
                *timeFormat = time.RFC3339Nano
        case "Kitchen":
                *timeFormat = time.Kitchen
        case "Stamp":
                *timeFormat = time.Stamp
        case "StampMilli":
                *timeFormat = time.StampMilli
        case "StampMicro":
                *timeFormat = time.StampMicro
        case "StampNano":
                *timeFormat = time.StampNano
        default:
            log.Printf("Invalid time format: %s\n\n",requestedFormat)
            Usage()
        }
}

