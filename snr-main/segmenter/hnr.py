import numpy as np

def harmonic_to_noise(spectrum, samplerate, frequency_buckets, liftering_active=True):
    # reference http://www.cs.northwestern.edu/~pardo/courses/eecs352/papers/boersma-pitchtracking.pdf
    auto_correlation = compute_harmonicity(spectrum, samplerate, frequency_buckets)
    HNR = compute_HNR(auto_correlation)
    return HNR


def compute_harmonicity(spectrum, samplerate, frequency_buckets, liftering_active=True):  
    # Note that,
    # onesided FFT has dimension of 513 and contains frequency range [0,8000]
    # twosided FFT has dimension of 1024 and contains frequency range [-8000,8000). (8000 isn't included)
    if(liftering_active):
        # reference https://groups.csail.mit.edu/sls/publications/2011/Ekapol_Interspeech11.pdf
        # Discard non-speech periodic signals by doing
        # bandpass cepstral liftering (80-400Hz).                   
        cepstrum = np.fft.ifft(np.log(spectrum), axis=0)
        # liftering
        # index_80Hz = 200
        # index_400Hz = 40
        index_80Hz = samplerate // 80
        index_400Hz = samplerate // 400
        
        zero_out = np.ones_like(cepstrum)
        zero_out[index_400Hz:index_80Hz+1,:] = 0.
        zero_out[-index_80Hz:-index_400Hz+1,:] = 0.
        
        cepstrum[np.where(zero_out == 1)] = 0.
        spectral = np.fft.fft(cepstrum, axis=0)
        spectral = np.exp(spectral)
        twosided_spectrum = np.vstack((spectral[-1], np.conjugate(np.flip(spectral[1:-1], axis=0)/2.), spectral[0], spectral[1:-1,:]/2.))
    else:                
        # From observation of results from specgram function, 
        # I found that value of onesided FFT are equal to 2*value of their correlation frequency in twosided FFT.
        # Excluding, 0 and the last frequency which is 8000 in this case.
        # vstack of frequency [-8000], (-8000,0), [0], (0,8000)
        # value is divided by    1   ,    2    ,  1,   2            
        twosided_spectrum = np.vstack((spectrum[-1], np.flip(spectrum[1:-1], axis=0)/2., spectrum[0], spectrum[1:-1,:]/2.))
    
    lag_t = np.arange(0, frequency_buckets, 1, dtype=np.float) / samplerate
    # r_a: wave autocorrelation
    # r_w: norm window autocorrelation
    # r_x: r_a / r_w, used autocorrelation
    r_a = np.fft.ifft(twosided_spectrum, axis=0).real # Equation (16)
    r_w = hanning_window(lag_t, float(frequency_buckets)/float(samplerate)) # Equation (8)
    r_a_0 = np.copy(r_a[0,:]) # lag 0 of every time step
    r_a_0[np.where(r_a_0 == 0.)] = 1 # prevent division by zero
    r_a_norm = r_a / r_a_0 # Normalized, Equation (2) 
    # Three periods of a wave were contained in a window was an assumption.
    # So, only 1/3 of window length was used due to sampling resolution (Fig. 1)
    r_x = r_a_norm[:frequency_buckets//3, :] / r_w[:frequency_buckets//3, None] # Equation (9)
    return r_x


def compute_HNR(auto_correlation):
    # this function do Equation (4)
    # low_dB_clipped is used to handle a case of which a frame has only zeros.
    # Since log of zero is not defined, handling of the case is required.
    low_dB_clipped = -150.
    # high_dB_clipped is used to handle a case of which a window is very small. 
    # Even if only the first half of autocorrelation is used, 
    # the possibility of having a result that is equal to or greater than one still exists.
    high_dB_clipped = 150. 

    a_max = np.max(auto_correlation[1:,:],axis=0)
    # reversed of decibel
    high_clipped = 10.**(high_dB_clipped / 10.) / (1. + 10.**(high_dB_clipped / 10.))
    low_clipped = 10.**(low_dB_clipped / 10.) / (1. + 10.**(low_dB_clipped / 10.))
    a_clipped = np.clip(a_max, low_clipped, high_clipped)
    # in decibel
    return 10. * np.log10(a_clipped / (1. - a_clipped))


def hanning_window(lags, window_length):
    window_length = float(window_length)
    a = 1. - np.abs(lags)/window_length
    b = 2./3. + 1./3. * np.cos(2. * np.pi * lags / window_length)
    c = 1./(2.*np.pi) * np.sin(2. * np.pi * np.abs(lags) / window_length)
    return a * b + c