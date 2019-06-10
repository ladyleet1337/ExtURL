# ExtURL
ExtURL is a tiny tool that extracts a lot of URLs from one page.  
**Warning:** ExtURL scans the website recursively, so a large number of requests may be sent.
# Installation

Download ExtURL from [releases](https://github.com/Ry0taK/ExtURL/releases) page.

# Usage
```
$ ./ExtURL
  -o string
    	Output file path (default "output.txt")
  -s boolean
        Whether to target only the same domain (default true)
  -t int
    	Number of threads (default 10)
  -u string
    	Base URL to find URLs (required)
```
