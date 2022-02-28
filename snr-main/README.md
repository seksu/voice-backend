# snr
snr evaluation

## Dependencies
- Python 3
- numpy
- matplotlib
- Soundfile
- `segmenter` repository https://github.com/SLSCU/segmenter

You can use `pip install requirement.txt` to install all required packages


## Usage 

```
python snr.py <PATH> --method <METHOD> --snr <VAD_THR> --vad <VAD_THR> --clipping <CLIPPING_THR> --energy <ENERGY_THR>
```
where `<PATH>` can be either a single audio file or a directory containing multiple files. `<SNR_THR>`, `<VAD_THR>`, `<CLIPPING_THR>`, and `<ENERGY_THR>` are the threshold of SNR, VAD, Clipping and Loudness respectively. `<METHOD>` is a sound segmentation method.

The default of `<SNR_THR>`, `<VAD_THR>`, `<CLIPPING_THR>`, and `<ENERGY_THR>` are `5`, `0.6`, `0.98` and `0.01` respectively.

While running, the script will output a result on console as well as logging into a JSON file.

## Evaluation Result

**The audio is considered good if all OKs are printed**

```
asdf_1234567890
length  : 3.1573125
VAD             : 84.04%        : OK
SNR             : 31.77         : OK
Loudness        : 0.0943        : OK
Clipping        : 0 samples     : OK
```
Where
- Length is in seconds
- VAD stands for Voice Activity Detection (how well voice can be detected)
- SNR stands for Signal-to-Noise Ratio (how clean the audio is)
- Loudness is the overall loudness of the audio
- Clipping is when the audio signal is beyond maximum range, usually caused by speaking too loud or volume spikes

### Notes

- If VAD is too low, it will print `UNCLEAR`
- If SNR is too low, it will print `NOISY`
- If Loudness is too low, it will print `QUIET`
- If clipping occurs, it will print `CLIPPING`
