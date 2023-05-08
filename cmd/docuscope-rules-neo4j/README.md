# DocuScope Dictionary LAT Rules to Graph Converter

Connects to a [Neo4J](https://neo4j.com) graph database and inserts LAT rules
to be consumed by CMU_Sidecar/docuscope-tag>.

## Setup
This tool expects the following environment variables to be set to configure the target [Neo4J](https://neo4j.com) database.  This tool can get these values from a .env file.

| Variable | Description |
| --- | --- |
| **NEO4J_DATABASE** | Database identifier |
| **NEO4J_URI** | URI for the database host |
| **NEO4J_USER** | Username to access the database |
| **NEO4J_PASSWORD** | Password to access the database |

## Input
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

## Usage
1. `docuscope-rules-neo4j <path>`
<path> is the path to the top level directory of a DocuScope language model (eg) `dictionaries/default`.

Execute `docuscope-rules-neo4j -h` for available command line arguments.
