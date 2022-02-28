import sys
import numpy as np
import soundfile as sf
import os
import glob
import json
import subprocess
import argparse
from segmenter.n_segmenter import Segmenter, add_segmenter_arguments

def segments_parser(segments_path, sel_spk='spk1'):
    signal_segments = dict()
    segments_names = dict()
    with open(segments_path) as f:
        for line in f:
            seg_id, wav_id, seg_st, seg_end = line.strip().split(' ')
            # if(not sel_spk in seg_id):
            #     continue
            seg_st = float(seg_st)
            seg_end = float(seg_end)
            if(wav_id in signal_segments):
                signal_segments[wav_id] += [ (seg_st, seg_end) ] 
                segments_names[wav_id] += [seg_id]
            else:
                signal_segments[wav_id] =  [ (seg_st, seg_end) ]
                segments_names[wav_id] = [seg_id]
    for wav_id in signal_segments:
        signal_segments[wav_id] = sorted(signal_segments[wav_id])
    return signal_segments, segments_names

def wav_parser(wavefile_path):
    wav = dict()
    for f in os.listdir(wavefile_path):
        wav_id = f.split('.')[0]
        wav[wav_id] = wavefile_path + f
    return wav

def load_wav(wav_path):
    raw_waveform, sample_rate = sf.read(wav_path)
    np.random.seed(99)
    waveform = np.clip(raw_waveform, -0.98, 0.98) + (np.random.rand(*raw_waveform.shape) - 0.5) * 0.00001
    return waveform, sample_rate, raw_waveform

def get_idx_from_ctm(timestamp, sample_rate=8000):
    return int(timestamp * sample_rate)

def get_file_snr(waveform, sample_rate, signal_segments):
    """
        Given signal_segments, this function calculate SNR
            on the entire wav file
        The signal power is an average of power of each signal segment
        The noise power is calculate from the rest of the file
            that is not marked by signal_segments
    """    
    is_signal = np.zeros(waveform.shape, dtype=bool)
    for sig in signal_segments:
        start_idx = get_idx_from_ctm(sig[0], sample_rate)
        end_idx = get_idx_from_ctm(sig[1], sample_rate) + 1
        is_signal[start_idx:end_idx] = 1
    S_p= (waveform[is_signal]**2).mean()
    N_p = (waveform[~is_signal]**2).mean()
    # print((waveform[~is_signal]**2).max())
    if N_p == 0 or np.isnan(N_p):
        N_p = S_p/1e3
    SNR_dB = 10 * np.log10(S_p / N_p)
    return SNR_dB, S_p, N_p

def snr(path, signal_segments, snr_thr=5, vad_thr=0.6, energy_thr=0.01, clipping_thr=0.98):
    waveform, sample_rate, raw_waveform = load_wav(path)  
    time = 0
    snr = 0

    # VAD
    for st,en in signal_segments:
        time += en - st
    # SNR
    if method == 'heuristics':
        snr, S_p, N_p = get_file_snr(waveform, sample_rate, signal_segments)
    elif method == 'wada':
        from wada_snr import wada_snr
        snr, S_p, N_p = wada_snr(waveform)
    else:
        print('Unknown method')
        exit(1)

    length = signal_segments[-1][-1] - signal_segments[0][0]
    vad = time/length
    energy = np.round(np.sqrt(S_p),4)
    vad_msg = 'OK'
    snr_msg = 'OK'
    energy_msg = 'OK'
    clipping_msg = 'OK'
    clipped_samples = 0
    if vad < vad_thr:
        vad_msg = 'UNCLEAR'
    if snr < snr_thr:
        snr_msg = 'NOISY'
    if energy < energy_thr:
        energy_msg = 'QUIET'
    if np.any(np.abs(raw_waveform) > clipping_thr):
        clipping_msg = 'CLIPPING'
        clipped_samples = np.sum(np.abs(raw_waveform) > 0.98)

    # print(wav_id)
    # print('length\t:',length)
    # print('VAD\t\t: {:.2f}% \t: {}'.format(100*vad,vad_msg))
    # print("SNR\t\t: {:.2f}  \t: {}".format(snr,snr_msg))
    # print('Loudness\t: {} \t: {}'.format(energy,energy_msg))
    # print('Clipping\t: {} samples\t: {}'.format(clipped_samples,clipping_msg))  
    # print()

    result = { 'id':wav_id,
            'length':length,
            'VAD':{
                'value':vad,
                'status':vad_msg
            },
            'SNR':{
                'value':snr,
                'status':snr_msg
            },
            'energy':{
                'value':energy,
                'status':energy_msg
            },
            'clipping':{
                'value':int(clipped_samples), # fix numpy.int64 JSON not serializable
                'status':clipping_msg
            }}
    print(result)
    return result

class ArgumentParser(argparse.ArgumentParser):
    def error(self, message):
        self.print_help(sys.stderr)
        self.exit(2, '%s: error: %s\n' % (self.prog, message))

if __name__ == '__main__':
    
    parser = ArgumentParser()
    parser.add_argument('path', type=str, help='path of wave files')
    parser.add_argument('--method', type=str, default='heuristics', help='segmentation method')
    parser.add_argument('--snr', type=float, default=5, help='SNR threshold')
    parser.add_argument('--vad', type=float, default=0.6, help='VAD threshold')
    parser.add_argument('--energy', type=float, default=0.01, help='loudness thershold')
    parser.add_argument('--clipping', type=float, default=0.98, help='clipping threshold')
    
    parser_segmenter = parser.add_argument_group('segmenter')
    add_segmenter_arguments(parser_segmenter)

    args = parser.parse_args()
    path = args.path
    method = args.method
    snr_thr = args.snr
    vad_thr = args.vad
    energy_thr = args.energy
    clipping_thr = args.clipping

    # print(args)

    if path.endswith('.wav') or path.endswith('.flac'): # if path is a file not folder
        files = [path]
    else:
        # files = [os.path.join(path,x) for x in os.listdir(path)]
        files = sorted(glob.glob(path + '/*.wav')) + sorted(glob.glob(path + '/*.flac'))

    # print(files)
    log = []
    for wav_path in files:
        tmp_path = 'tmp.wav'
        convertCmd = "sox -R {} -c 1 -b 16 -r 16000 -t wav {}".format(wav_path,tmp_path)
        process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
        process.communicate()
        try:
            data, samplerate = sf.read(tmp_path)
        except:
            print('sox failed for {}'.format(wav_path.split('/')[-1].split('.')[0]))
            result = { 'id':wav_path.split('/')[-1].split('.')[0],
                       'sox failed':True}
            log.append(result)
            convertCmd = "rm {}".format(tmp_path)
            process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
            process.communicate()
            continue

        # no sound check
        if len(data)/samplerate <= 0.5 or sum(abs(data))/len(data) <= 5e-5:
            print(wav_path.split('/')[-1].split('.')[0], 'no sound')
            result = { 'id':wav_path.split('/')[-1].split('.')[0],
                       'no sound':True}
            log.append(result)
            convertCmd = "rm {}".format(tmp_path)
            process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
            process.communicate()
            continue

        # segmenter
        sg = Segmenter(
        min_formant_duration=args.min_formant_duration, 
        max_duration=args.max_duration,
        max_gap=args.gap, 
        formant_band=args.formant_band,
        front_porch=args.front,
        back_porch=args.back,
        hnr_max_gap=args.gap, 
        hnr_front_porch=args.front,
        hnr_back_porch=args.back,
        log_var_formant_low_threshold=args.log_var_formant_low_threshold,
        log_var_formant_high_threshold=args.log_var_formant_high_threshold,
        frequency_buckets=1024, 
        win=None,
        hnr_low_threshold=args.hnr_low_threshold,
        hnr_active=args.hnr,
        liftering_active=args.liftering,
        energy_raw_threshold=args.energy_raw_threshold,
        energy_threshold=args.energy_threshold,
        energy_band=args.energy_band)
        
        segments = sg(tmp_path)

        if segments == None:
            print(wav_path.split('/')[-1].split('.')[0], 'no segments')
            
            result = { 'id':wav_path.split('/')[-1].split('.')[0],
                       'no segments':True}
            # exit(0)
            log.append(result)
            convertCmd = "rm {}".format(tmp_path)
            process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
            process.communicate()
            continue
        
        signal_segments = segments.signal_slices
        wav_id = os.path.splitext(os.path.split(wav_path)[-1])[0]
        result = snr(tmp_path, signal_segments, snr_thr=snr_thr, vad_thr=vad_thr, energy_thr=energy_thr, clipping_thr=clipping_thr)
        
        log.append(result)
        convertCmd = "rm {}".format(tmp_path)
        process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
        process.communicate()

    # print(len(log))
    if not os.path.exists('log'):
        os.mkdir('log')
    with open('log/log_{}.json'.format(os.path.splitext(os.path.split(path)[-1])[0]),'w') as f:
        json.dump(log, f)
