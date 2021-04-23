# DocuScope Rules

Tools for preprocessing DocuScope dictionary files to JSON to be consumed by the DocuScope tagger.

## docuscope-rules

Generates JSON from a directory containing a DocuScope dictionary.
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

See api/docuscope_rules_schema.json for the schema of the resulting JSON.
