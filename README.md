# DocuScope Sidecar Tools

Tools for preprocessing DocuScope language model files to be consumed by
CMU_Sidecar/docuscope-tag>
and
CMU_SIDECAR/docuscope-classroom>.

## Administration and Support

For any questions regarding overall project or the language model used, please contact <suguru@cmu.edu>.

The project code is supported and maintained by the [Eberly Center](https://www.cmu.edu/teaching/) at [Carnegie Mellon University](www.cmu.edu). For help with this fork, project or service please contact <eberly-assist@andrew.cmu.edu>.

## Commands

- [docuscope-rules-neo4j](cmd/docuscope-rules-neo4j/README.md) imports a DocuScope dictionary to a graph database used by CMU_Sidecar/docuscope-tag>.
- [docuscope-rules](cmd/docuscope-rules/README.md) converts a DocuScope dictionary to JSON for easier consumption by CMU_Sidecar/docuscope-tag>. An alternative to using a graph database.
- [docuscope-wordclasses](cmd/docuscope-wordclasses/README.md) converts DocuScope dictionary _wordclasses.txt file to JSON for easier consumption by CMU_Sidecar/docuscope-tag> and CMU_Sidecar/docuscope-classroom>.
- [docuscope-tones](cmd/docuscope-tones/README.md) converts DocuScope dictionary _tones.txt file to JSON for consumption by CMU_Sidecar/docuscope-classroom>.
- [docuscope-rules-db](cmd/docuscope-rules-db/README.md) **DISCONTINUED** converts DocuScope dictionary to JSON to be consumed by a NoSQL database used by CMU_Sidecar/docuscope-tag>.

## Acknowledgments

This project was partially funded by the [A.W. Mellon Foundation](https://mellon.org/), [Carnegie Mellon University](https://www.cmu.edu/)'s [Simon Initiative](https://www.cmu.edu/simon/) Seed Grant, and the [Berkman Faculty Development Fund](https://www.cmu.edu/proseed/proseed-seed-grants/berkman-faculty-development-fund.html).
