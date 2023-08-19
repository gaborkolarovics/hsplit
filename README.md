# hsplit
Split a file into pieces with rolling hash

Output pieces of FILE to PREFIXaaaa, PREFIXaaab, ...; default trailing zero bits in the rolling checksum is 16, and default PREFIX is 'x'.

```
  Usage:
    hsplit [FILE] [PREFIX]

  Positional Variables: 
    FILE     filename_desc (Required)
    PREFIX   prefix_desc (default: x)
  Flags: 
       --version          Displays the program version string.
    -h --help             Displays help with available flag, subcommand, and positional value parameters.
    -b --bits             SplitBits is the number of trailing zero bits in the rolling checksum required to produce a chunk. (default: 16)
    -m --minSize          MinSize is the minimum chunk size. (default: 1024)
    -f --fanout           Fan-out of the nodes in the tree produced. (default: 8)
    -p --minBytesOfPart   MinBytesOfPart is the minimum file part size (default: 1024)
```

# Inspiration 
* https://github.com/bobg/hashsplit
* https://ostechnix.com/split-combine-files-command-line-linux/
