import numpy as np
import matplotlib.mlab as mlab
import soundfile as sf
import argparse

class Segmenter():
    def __init__(self, 
            min_formant_duration=0, max_gap=0.5,
            formant_band=(200,4800),
            front_porch=0.10, back_porch=0.15, 
            hnr_max_gap=0.5, hnr_front_porch=0.10, hnr_back_porch=0.15,
            log_var_formant_low_threshold=None,
            log_var_formant_high_threshold=1.0,
            frequency_buckets=1024, win=None,
            hnr_low_threshold=None,
            hnr_active=False,
            liftering_active=True,
            energy_raw_threshold=0.08,
            energy_threshold=1e-4,
            energy_band=(100,1000),
            max_duration=None
            ):

        self.min_formant_duration = min_formant_duration
        self.max_gap = max_gap
        self.front_porch = front_porch
        self.back_porch = back_porch
        self.hnr_max_gap = hnr_max_gap
        self.hnr_front_porch = hnr_front_porch
        self.hnr_back_porch = hnr_back_porch
        self.formant_band = formant_band
        self.energy_threshold = energy_threshold
        self.energy_raw_threshold = energy_raw_threshold
        self.energy_band = energy_band
        if max_duration is None:
            self.max_duration = 100000
        else:
            self.max_duration = max_duration
        self.total_duration = None
     
        # Set default values
        # liftering, logE threshold = (0.4,1.0), hnr threshold = 0
        # only HNR, logE threshold = (0.6,1.0), hnr threshold = 2
        # logEnergy, logE threshold = (0.7,1.0)

        if(hnr_active):
            if(liftering_active):
                if(hnr_low_threshold == None):
                    hnr_low_threshold = 0
                if(log_var_formant_low_threshold == None):
                    log_var_formant_low_threshold = 0.4
            else:
                if(hnr_low_threshold == None):
                    hnr_low_threshold = 0.5 # 2
                if(log_var_formant_low_threshold == None):
                    log_var_formant_low_threshold = 0.6
        else:
            if(log_var_formant_low_threshold == None):
                log_var_formant_low_threshold = 0.4

        self.logvar_formant_low_threshold = log_var_formant_low_threshold
        self.logvar_formant_high_threshold = log_var_formant_high_threshold
        self.frequency_buckets = frequency_buckets
        self.win = win
        self.threshold = None
        self.hnr_low_threshold = hnr_low_threshold # in dB
        self.hnr_active = hnr_active
        self.liftering_active = liftering_active


    def load_audio(self, audio):
        data, samplerate = sf.read(audio)
        total_duration = len(data) / float(samplerate)
        return data, samplerate, total_duration


    def compute_spectrum(self, audio_data, samplerate):
        # Add a little noise to avoid 0 variance
        # audio_data = audio_data[:] + np.random.random(len(audio_data)) * .00001
        np.random.seed(99)
        audio_data = np.clip(audio_data, -0.98, 0.98) + (np.random.rand(*audio_data.shape) - 0.5) * 0.00001
    
        spectrum, freqs, t = mlab.specgram(
            x=audio_data, 
            NFFT=self.frequency_buckets, 
            Fs=samplerate,
            window=mlab.window_hanning, 
            noverlap=self.frequency_buckets/2,
            sides='default')

        lower = int(self.formant_band[0] / (samplerate / float(self.frequency_buckets)))
        higher = int(self.formant_band[1] / (samplerate / float(self.frequency_buckets)))

        # Size of a spectral frame in msec
        frame_duration = t[1] - t[0]
        # Compute the normalized log variance for power in the formant band
        # Think of it as "wiggle" in the band
        bandpass = spectrum[lower:higher, :]
        var = bandpass.var(axis=0)
        logvar = np.log(var)
        logvar -= logvar.min()
        logvar_max = logvar.max()
        logvar /= logvar.max()
        logvar = logvar
        return spectrum, logvar, frame_duration, logvar_max


    def selectN_aboveR(self, x, N = 5, R = 3):
        # x in list of boolean
        # N elements will been set to True if there are True not less than R elements, 
        # otherwise they will been to False.
        # return list of boolean
        cumulative_true = 0 # dynamic programming sum
        modified_x = [False] * len(x)
        for i in range(len(x)):
            cumulative_true += int(x[i])
            if(i<N-1):      
                continue
            
            if(cumulative_true >= R):
                modified_x[i-N+1:i+1] = [True] * N

            cumulative_true -= int(x[i-N+1])

        return modified_x


    def get_slices(self, where):
        # where is a list of booleans
        # return list of edges of regions of true
        runs = []
        left = None
        for i, w in enumerate(list(where)):
            # Looking for a left edge
            if left is None:
                if w:
                    left = i
                continue
            if not w:
                runs.append((left,i))
                left = None
        if left is not None:
            runs.append((left, i))
        return runs


    def energy_speech(self, spectrum, samplerate, spectrumc):
        # compute energy signal slices by summing energy of all frequency bands
        # speech = np.zeros_like(audio_data)
        #speech = np.zeros_like(spectrum)

        lower = int(self.energy_band[0] / (samplerate / float(self.frequency_buckets)))
        higher = int(self.energy_band[1] / (samplerate / float(self.frequency_buckets)))

        # bandpass = spectrum[lower:higher, :].sum(axis=0)
        bandpass = spectrum.sum(axis=0)
        above_threshold = self.energy_threshold < bandpass
        #speech[np.where(above_threshold)] = 1.0
        signal_slices = self.get_slices(above_threshold)
        signal_slices = [(spectrumc.seconds(l), spectrumc.seconds(r)) for (l, r) in signal_slices]
        return signal_slices

    # NOT USE 
    def energy_speech_raw(self, audio_data, samplerate):
        # compute energy signal slices from raw signal (PCM)
        speech = np.zeros_like(audio_data)
        above_threshold = self.energy_raw_threshold < abs(audio_data)
        speech[np.where(above_threshold)] = 1.0
        signal_slices = self.get_slices(above_threshold)
        signal_slices = [(l/samplerate, r/samplerate) for (l, r) in signal_slices]
        return signal_slices


    def find_speech(self, logvar, spectrumc, hnr=[]):
        hnr_signal_slices = None
        above_hnr_threshold = np.ones_like(logvar, dtype=bool)
        
        # compute hnr signal slsices
        if(self.hnr_active):
            # get speech slices by considering hnr threshold
            above_hnr_threshold = hnr > self.hnr_low_threshold
            hnr_speech = np.zeros_like(hnr) 
            above_hnr_threshold = self.selectN_aboveR(above_hnr_threshold)
            hnr_speech[np.where(above_hnr_threshold)] = 1.0

            hnr_signal_slices = self.get_slices(above_hnr_threshold)
            hnr_signal_slices = [(spectrumc.seconds(l), spectrumc.seconds(r)) for (l, r) in hnr_signal_slices]
        
        # compute logvar signal slices
        # Get a threshold which we think a formant is present
        v = list(logvar)
        low_threshold_index = int(len(v) * self.logvar_formant_low_threshold)
        high_threshold_index = int(len(v) * self.logvar_formant_high_threshold)-1
        sorted_v = sorted(v)
        low_threshold = sorted_v[low_threshold_index]
        high_threshold = sorted_v[high_threshold_index]
        if self.win is not None:
            filtered = sig.convolve(logvar, self.win, mode='same') / sum(self.win)
        else:
            filtered = logvar

        speech = np.zeros_like(logvar)

        between_threshold = (low_threshold < filtered) & (filtered <= high_threshold)
        speech[np.where(between_threshold)] = 1.0
        signal_slices = self.get_slices(between_threshold)

        # Convert signal slices to seconds
        signal_slices = [(spectrumc.seconds(l), spectrumc.seconds(r)) for (l, r) in signal_slices]
        return signal_slices, hnr_signal_slices


    def merge_segment(self, signal_slices, 
                        max_gap = None, min_formant_duration = None,
                        front_porch = None, back_porch = None):
        # Default value
        max_gap = self.max_gap if max_gap == None else max_gap
        min_formant_duration = self.min_formant_duration if min_formant_duration == None else min_formant_duration
        front_porch = self.front_porch if front_porch == None else front_porch
        back_porch = self.back_porch if back_porch == None else back_porch

        # Remove formant peaks before merge segment
        # signal_slices = [[l, r] for (l, r) in signal_slices if r - l > min_formant_duration]
        # if len(signal_slices) == 0:
        #     # No signal found
        #     return signal_slices
        
        #Now bridge signal across short gaps
        bridged = [list(signal_slices[0])]
        for l, r in signal_slices[1:]:
            if l - bridged[-1][1] <= max_gap and r -  bridged[-1][0] <= self.max_duration:
                bridged[-1][1] = r
            else:
                bridged.append([l, r])

        # Remove very short formant peaks
        signal_slices = [[l, r] for (l, r) in bridged if r - l > min_formant_duration]
        if len(signal_slices) == 0:
            # No signal found
            return signal_slices

        # Now put a front and back porch on the segment
        signal_slices = [[max(0.0, l - front_porch), min(r + back_porch, self.total_duration)] for (l, r) in signal_slices]
        # Now concatenate potentially overlapping segments
        bridged = [signal_slices[0]]
        for l, r in signal_slices[1:]:
            if l - bridged[-1][1] <= max_gap:
                if r -  bridged[-1][0] <= self.max_duration:
                    bridged[-1][1] = r
                elif l - bridged[-1][1] < 0:
                    bridged[-1][1] = l
                    bridged.append([l, r])
                else:
                    bridged.append([l, r])
            else:
                bridged.append([l, r])

        
        signal_slices = [(round(l, 4), round(r, 4)) for (l, r) in bridged]
        return signal_slices


    def extract_noise(self, signal_slices, total_duration):
        # Extract the chunks of non_speech
        noise_slices = []
        lastr = 0.0
        for start, stop in signal_slices:
            noise_slices.append((lastr, start))
            lastr = stop
        if lastr < total_duration:
            noise_slices.append((lastr, total_duration))
        return noise_slices


    def find_intersect(self, seg1, seg2):
        min1, max1 = seg1
        min2, max2 = seg2
        if(max(0, min(max1, max2) - max(min1, min2)) == 0):
            return None
        else:
            return (max(min1, min2), min(max1,max2))


    def intersect_slices(self, slices1, slices2):
        s1_i = 0
        s2_i = 0
        slices_result = list()
        while(s1_i < len(slices1) and s2_i < len(slices2)):
            intersect = self.find_intersect(slices1[s1_i], slices2[s2_i])
            if(intersect != None and intersect[1] - intersect[0] > self.min_formant_duration):
                slices_result.append(intersect)
            if(slices1[s1_i][1] < slices2[s2_i][1]):
                s1_i += 1
            else:
                s2_i += 1
        
        return slices_result


    def __call__(self, audio_path):
        audio_data, samplerate, total_duration = self.load_audio(audio_path)
        self.total_duration = total_duration
        # Energy-based
        spectrum, logvar, frame_duration, logvar_max = self.compute_spectrum(audio_data, samplerate)
        hnr = []
        if(self.hnr_active):
            # Harmonic-based
            hnr = harmonic_to_noise(spectrum, samplerate, self.frequency_buckets)

        # magic number for too low amplitude
        # if logvar_max < 5:
        #     return None

        spectrumc = SecToIndex(duration=frame_duration)
        energy_slices = self.energy_speech(spectrum, samplerate, spectrumc)
        # energy_slices = self.energy_speech_raw(audio_data, samplerate)
        if(len(energy_slices) == 0):
            return None

        energy_slices = self.merge_segment(energy_slices, max_gap=self.max_gap)

        signal_slices, hnr_signal_slices = self.find_speech(logvar, spectrumc, hnr)

        signal_slices = self.merge_segment(signal_slices, max_gap=self.max_gap)
        
        signal_slices = self.intersect_slices(signal_slices, energy_slices)

        if(hnr_signal_slices):
            hnr_signal_slices = self.merge_segment(
                hnr_signal_slices,
                max_gap = self.hnr_max_gap, 
                front_porch = self.hnr_front_porch, 
                back_porch = self.hnr_back_porch)
            signal_slices = self.intersect_slices(signal_slices, hnr_signal_slices)
        noise_slices = self.extract_noise(signal_slices, total_duration)
        if(len(signal_slices) == 0):
            return None

        return Segmentation(audio_path, audio_data, samplerate, signal_slices, noise_slices)

def add_segmenter_arguments(parser):
    parser.add_argument('--front',type=float,default=0.1,help='front padding duration')
    parser.add_argument('--back',type=float,default=0.15,help='back padding duration')
    parser.add_argument('--gap',type=float,default=0.5,help='segments with gap less than this value will be merged')
    parser.add_argument('--max_duration',type=float,default=None,help='max duration segments can be, set to None to disable')
    parser.add_argument('--min_formant_duration',type=float,default=0,help='segments shorter than min_formant_duration will be ignored')
    parser.add_argument('--hnr',action='store_true',help='enable hnr')
    parser.add_argument('--liftering',action='store_true',help='enable hnr')
    parser.add_argument('--hnr_low_threshold',type=float,default=None,help='enable liftering') 
    parser.add_argument('--log_var_formant_low_threshold',type=float,default=0.35,help='log var formant low threshold') 
    parser.add_argument('--log_var_formant_high_threshold',type=float,default=1.0,help='log var formant high threshold') 
    parser.add_argument('--energy_raw_threshold',type=float,default=None,help='energy raw threshold') 
    parser.add_argument('--energy_threshold',type=float,default=1e-8,help='energy threshold') 
    parser.add_argument('--energy_band',nargs=2,type=int,default=(100,1000),help='energy band') 
    parser.add_argument('--formant_band',nargs=2,type=int,default=(200,4800),help='formant band') 
    parser.add_argument('--raw',action='store_true',help='do not preprocess audio file (may cause problem)')


if __name__ == '__main__':
    from segmentation import Segmentation, SecToIndex
    from hnr import harmonic_to_noise
    import sys
    import subprocess
    import os

    parser = argparse.ArgumentParser()
    parser.add_argument('audio_path')
    add_segmenter_arguments(parser)
    args = parser.parse_args()
    # print(args)

    if not args.raw:
        tmp_path = 'tmp.wav'
        convertCmd = "sox -R {} -c 1 -b 16 -r 16000 -t wav {}".format(args.audio_path,tmp_path)
        process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
        process.communicate()

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

    result = sg(tmp_path)
    if result == None:
        print([])
    else:
        print(result.signal_slices)

    if not args.raw:
        convertCmd = "rm {}".format(tmp_path)
        process = subprocess.Popen(convertCmd.split(), stdout=subprocess.PIPE)
        process.communicate()
else:
    from .segmentation import Segmentation, SecToIndex
    from .hnr import harmonic_to_noise
