# DocuScope Rules

Tools for preprocessing DocuScope dictionary files to JSON to be consumed by the DocuScope tagger.

## Administration and Support

For any questions regarding overall project or the language model used, please contact suguru@cmu.edu

The project code is supported and maintained by the [Eberly Center](https://www.cmu.edu/teaching/) at [Carnegie Mellon University](www.cmu.edu). For help with this fork, project or service please contact eberly-assist@andrew.cmu.edu.

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
There has to be a blank line between each CLASS.

See api/docuscope_rules_schema.json for the schema of the resulting JSON.

## Acknowledgements

This project was partially funded by the [A.W. Mellon Foundation](https://mellon.org/), [Carnegie Mello University](https://www.cmu.edu/)'s [Simon Initiative](https://www.cmu.edu/simon/) Seed Grant and [Berkman Faculty Development Fund](https://www.cmu.edu/proseed/proseed-seed-grants/berkman-faculty-development-fund.html).
