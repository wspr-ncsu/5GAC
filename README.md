# 5GAC
5GAC is an analysis tool for 5G core implementations. It's main purpose is to identify which Network Functions (NFs) call APIs in other NFs for the purpose of least-privilege analysis. 5GAC is designed to analyze 4 open-source 5G cores and 2 versions of one 5G core:

- free5GC
- Open5GS (Based on 3GPP Release 16)
- Open5GS (Based on 3GPP Release 17)
- OpenAirInterface (OAI)

This repository contains the CodeQL files to analyze the cores. 5GAC is extensible in that there is a higher level model defined in the [x] folder and each implementation is defined in subfolders. This allows much of the analysis abstractions to be re-used across code bases.

## Running 5GAC

The best way to run 5GAC is by using the included Dockerfile. Specifically, run 
```bash
docker build  ./ -t 5gac:latest && docker run -v $(pwd):/repo 5gac:latest
```

The analysis results will be added to the json folder in this directory. Each result file name is the 5G core that was analyzed followed by the date and time the analysis was run.

If you want to run 5GAC without docker, we recommend running the commands one by one as they are listed in the Dockerfile. This handles cloning the repos and any packages needed to compile and run all of the 5G cores.

## Output


## Citation
If you use 5GAC in your research, please cite our work using this bibtex entry:

```tex
TODO
```
