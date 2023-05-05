# DocuScope Tones

Tools for preprocessing DocuScope dictionary tone files to JSON to be consumed by CMU_Sidecar/docuscope-classroom>.

## Administration and Support

For any questions regarding overall project or the language model used, please contact suguru@cmu.edu

The project code is supported and maintained by the [Eberly Center](https://www.cmu.edu/teaching/) at [Carnegie Mellon University](www.cmu.edu). For help with this fork, project or service please contact eberly-assist@andrew.cmu.edu.

## docuscope-tones

Generates JSON from a DocuScope dictionary `_tones.txt` file.
The `_tones.txt` file should have a format simmilar to the following:

```
CLUSTER: <ClusterName>
DIMENSION: <DimensionName>
LAT|LAT*|CLASS: <LatName>
```

This is essentially a flattened hierarchical structure where each CLUSTER
has one or more DIMENSION entries and each DIMENSION has one or more LAT listings
prefixed with `LAT:`, `LAT*:`, or `CLASS:`.
The `<ClusterName>` should correspond to the `name` field for clusters in the `common-dictionary.json`
used with CMU_Sidecar/docuscope-classroom> and `<LatName>` should refer to the LAT ids used
in the dictionary used with CMU_Sidecar/docuscope-tag>.  `<DimensionName>`s are not currently used
in the related projects (though they are in other DocuScope projects) and must be unique.

See api/docuscope_tones_schema.json for the schema of the resulting JSON.

## Acknowledgements

This project was partially funded by the [A.W. Mellon Foundation](https://mellon.org/), [Carnegie Mello University](https://www.cmu.edu/)'s [Simon Initiative](https://www.cmu.edu/simon/) Seed Grant and [Berkman Faculty Development Fund](https://www.cmu.edu/proseed/proseed-seed-grants/berkman-faculty-development-fund.html).
