class SecToIndex(object):
    def __init__(self, frequency=None, duration=None):
        assert frequency or duration
        assert not frequency or not duration
        if frequency:
            self.frequency = frequency
            self.duration = 1.0 / frequency
        else:
            self.frequency = 1.0 / duration
            self.duration = duration

    def index(self, seconds):
        return int(seconds * self.frequency)

    def seconds(self, index):
        return float(index) * self.duration

class Segmentation(object):

    def __init__(self, name, audio_data, samplerate, signal_slices, noise_slices):
        self.name = name
        self.audio_data = audio_data
        self.samplerate = samplerate
        self.signal_slices = signal_slices
        self.noise_slices = noise_slices
        self.total_seconds = self.audio_data.shape[0] / float(self.samplerate)
        self.nsamples = self.audio_data.shape[0]
        self.sample_convert = SecToIndex(frequency=samplerate)
        self.seci = self.sample_convert.index
        self.isec = self.sample_convert.seconds


    def drop_audio(self):
        # Drop reference to the large data buffer
        self.audio_data = None


    def aslice(self, l, r):
        assert self.audio_data is not None
        return self.audio_data[self.seci(l):self.seci(r)]


    def signal_slice(self, i):
        l, r = self.signal_slices[i]
        return self.aslice(l, r)


    def noise_slice(self, i):
        l, r = self.noise_slices[i]
        return self.aslice(l, r)


    def aseconds(self, slices):
        seconds = 0
        for (l, r) in slices:
            seconds += (r - l)
        return seconds


    @property
    def signal_seconds(self):
        return self.aseconds(self.signal_slices)


    @property
    def noise_seconds(self):
        return self.aseconds(self.noise_slices)


    def concatenate_slices(self, slices, spacer_seconds=0.25):
        assert self.audio_data is not None
        if spacer_seconds > 0.0:
            spacer = np.zeros(int(.25 * self.samplerate), dtype=self.audio_data.dtype)
        else:
            spacer = None
        signals = []
        for (l, r) in slices:
            signals.append(self.aslice(l, r))
            if spacer is not None:
                signals.append(spacer)
        if signals:
            return np.hstack(signals)
        else:
            return None


    def concatenate_signal(self, spacer_seconds):
        return self.concatenate_slices(self.signal_slices, spacer_seconds)


    def concatenate_noise(self, spacer_seconds):
        return self.concatenate_slices(self.noise_slices, spacer_seconds)

  
    def __unicode__(self):
        return "seconds:{total_seconds:.03f} noise:{noise_seconds:.03f} slices:{nsignal}".format(
            total_seconds=self.total_seconds,
            noise_seconds=self.noise_seconds,
            nsignal=len(self.signal_slices)
        )

    __str__ = __unicode__