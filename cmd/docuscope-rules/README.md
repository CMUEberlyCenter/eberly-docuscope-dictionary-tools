# docuscope-rules

Generates JSON from a directory containing a DocuScope language model.
The directory should contain a collection of files where each file is named for a LAT with a .txt extension.
Each file contains sets of word or word classes that make up the patterns for that LAT, one pattern per line.
The directory should also contain the special file `_wordclasses.txt` which defines the word classes.
Each word class is like:

```
CLASS: <CLASSNAME>
<word>+

```

Where `<CLASSNAME>` is the all uppercase name of the class and
`words` are lowercase strings in that class, one per line.
There has to be a blank line between each CLASS.

See (../../api/docuscope_rules_schema.json) for the schema of the resulting JSON.

## Usage
1. `docuscope-rules <path> | gzip > default.json.gz`
<path> is the path to the top level directory of a DocuScope language model (eg) `dictionaries/default`.
Using gzip compression is optional but strongly recommended as the non-compressed result can be several gigabytes however it is highly regular and thus compresses down to under 100 megabytes.

Execute `docuscope-rules-neo4j -h` for available command line arguments.
