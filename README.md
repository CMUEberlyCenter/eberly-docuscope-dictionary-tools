# DocuScope Rules

Tools for preprocessing DocuScope language model files to JSON to be consumed by CMU_Sidecar/docuscope-tag>.

## Administration and Support

For any questions regarding overall project or the language model used, please contact <suguru@cmu.edu>.

The project code is supported and maintained by the [Eberly Center](https://www.cmu.edu/teaching/) at [Carnegie Mellon University](www.cmu.edu). For help with this fork, project or service please contact <eberly-assist@andrew.cmu.edu>.

## docuscope-rules

Generates JSON from a directory containing a DocuScope language model.
The directory should contain a collection of files where each file is named for a LAT with a .txt extension.
Each file containes sets of word or word classes that make up the patterns for that LAT, one pattern per line.
The directory should also contain the special file `_wordclasses.txt` which defines the word classes.
Each word class is like:

```
CLASS: <CLASSNAME>
<word>+

```

Where `<CLASSNAME>` is the all uppercase name of the class and
`words` are lowercase strings in that class, one per line.
There has to be a blank line between each CLASS.

See api/docuscope_rules_schema.json for the schema of the resulting JSON.

### Usage
1. `docuscope-rules <path>`
<path> is the path to the top level directory of a DocuScope language model (eg) `dictionaries/default`.

## docuscope-rules-neo4j
Expects the same information and format as for docuscope-rules with additional environment variables specifying the target [Neo4J](https://neo4j.com) database.  This tool can get these values from a .env file.

| Variable | Description |
| --- | --- |
| **NEO4J_DATABASE** | Database identifier |
| **NEO4J_URI** | URI for the database host |
| **NEO4J_USER** | Username to access the database |
| **NEO4J_PASSWORD** | Passwort to access the database |

## Acknowledgements

This project was partially funded by the [A.W. Mellon Foundation](https://mellon.org/), [Carnegie Mello University](https://www.cmu.edu/)'s [Simon Initiative](https://www.cmu.edu/simon/) Seed Grant, and the [Berkman Faculty Development Fund](https://www.cmu.edu/proseed/proseed-seed-grants/berkman-faculty-development-fund.html).
