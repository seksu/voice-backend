## segmenter

Voice Activity Detection (VAD) script

Requirement
- Python 3
- Numpy
- Matplotlib
- Soundfile

### Get slices (Single file)
```
python3 n_segment.py <audio_path> [OPTIONS]
```
Type `python3 n_segment.py -h` for available options.

### Get kaldi segments (File/Directory) 
```
get_segment.sh <path> <out_dir> [OPTIONS]
```
where
- `path` is a path to either a single audio file or a directory containing multiple files
- `out_dir` is a directory to store `segments` output file.

You can pass n_segmenter options directly, for example:
```
get_segment.sh a.wav out --front 0.5 --back 0.5
```