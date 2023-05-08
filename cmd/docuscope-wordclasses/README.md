# DocuScope Word Classes Converter
Generates the JSON representation of a DocuScope dictionary word classes for
consumption by
CMU_Sidecar/docuscope-tag>
and
CMU_SIDECAR/docuscope-classroom>.

## Usage
1. `docuscope-wordclasses <path> | gzip > wordclasses.json.gz`
Compression is optional but strongly recommended as the result is highly regular and thus has a very high compression ratio.
`<path>` is the directory path to the DocuScope dictionary that contains LAT files and the _wordclasses.txt file.

Execute `docuscope-worclasses -h` for command line options.
