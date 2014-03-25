filewalker
==========

I decided to learn/play with Google's [Go](http://golang.org) programming language. I also wanted to have a nice status board for details about files on my storage. To this end I made this simple utility to walk a given directory and generate a JSON report of it's contents.

There are plenty of utilities that can report on usage but they all have their own custom GUIs and data repos. I wanted something that would allow more flexibility on how I can store the data, report on it, and display it. The output is to stdout which can then be collected or fed into something like [Splunk](www.splunk.com) or [Elasticsearch](www.elasticsearch.org) with [Logstash](logstash.net) and generate a dashboard reporting whatever details you like.

Usage
=====

     Usage: ./filewalker-amd64 [ OPTION ]...  DIRECTORY

     REQUIRED
             DIRECTORY Target directory to scan

     OPTIONS
                    -h Show this usage help
                    -v Show version and license
               -pretty Pretty print JSON

     Exclusions        Can be colon delimited for multiple excludes
           -excludeDir Directory patterns to exclude
          -excludeFile File patterns to exclude

                       EXAMPLE Single:       -excludeDir "pattern1"
                       EXAMPLE Multiple:     -excludeFile "pattern1:pattern2:pattern3"

     Excluding fields from report
                    -t Excludes date/time file was scanned
                    -s Excludes file size
                    -e Excludes file extension
                    -u Excludes uid
                    -g Excludes gid
                    -i Excludes inode
                    -m Excludes mtime
                    -a Excludes atime
                    -c Excludes ctime

     Adding extra data to report
                    -k Space delimited key values.

                       EXAMPLE Single:       -k "key1:value1"
                       EXAMPLE Multiple:     -k "key1:value1 key2:value2 key3:value3"

     Include errors in report
               -errors Include errors in report

     Change time formating
                 -time Valid time formats: ANSIC, UnixDate, RubyDate, RFC822, RFC822Z, RFC850, RFC1123,
                       RFC1123Z, RFC3339, RFC3339Nano, Kitchen, Stamp, StampMilli, StampMicro, StampNano

                       EXAMPLE:             -time RFC3339

Sample Output
=============

Normal output to stdout is:

     {"ScanDate":"2014-03-24T20:45:35-04:00","Type":"f","Path":"/home/test","Name":"main.go","Extension":".go","Size":9240,"Uid":1000,"Gid":1000,"Inode":8522834,"Mtime":"2014-03-24T20:04:37-04:00","Atime":"2014-03-24T20:25:09-04:00","Ctime":"2014-03-24T20:04:37-04:00"}

It can also be "pretty"

     {
         "ScanDate": "2014-03-24T20:45:43-04:00",
         "Type": "f",
         "Path": "/home/test",
         "Name": "main.go",
         "Extension": ".go",
         "Size": 9240,
         "Uid": 1000,
         "Gid": 1000,
         "Inode": 8522834,
         "Mtime": "2014-03-24T20:04:37-04:00",
         "Atime": "2014-03-24T20:25:09-04:00",
         "Ctime": "2014-03-24T20:04:37-04:00"
     }
     
     
