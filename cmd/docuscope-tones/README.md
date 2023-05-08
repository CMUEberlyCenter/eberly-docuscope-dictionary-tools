# DocuScope Tones

Tools for preprocessing DocuScope dictionary tone files to JSON to be consumed by CMU_Sidecar/docuscope-classroom>.

## Input

Generates JSON from a DocuScope dictionary `_tones.txt` file.
The `_tones.txt` file should have a format similar to the following:

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

## Output
See (../../api/docuscope_tones_schema.json) for the schema of the resulting JSON.

## Usage
1. `docuscope_tones < _tones.txt > tones.json`

Execute `docuscope_tones -h` for command line help.
