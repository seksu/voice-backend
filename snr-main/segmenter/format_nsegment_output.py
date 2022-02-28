import os
import sys
import ast 

slices = sys.stdin.read()
segments = ast.literal_eval(slices)


rec_id = os.path.splitext(os.path.split(sys.argv[1])[-1])[0]
out_path = sys.argv[2] if len(sys.argv) > 2 else '.'


duration, gap = [],[]
last_end = 0
for start,end in segments:
    utt_id = "{}_{:07d}_{:07d}".format(rec_id,int(round(start*1000)),int(round(end*1000)))
    print(u"{} {} {} {}".format(utt_id,rec_id,start,end))

    duration.append(end - start)
    gap.append(start - last_end)
    last_end = end
