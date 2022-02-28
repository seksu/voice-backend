import os
import sys

# This file can be used for scoring VAD.

# test_dir is a folder contains segment index files
# label_dir is a folder contains .lab files

# Sample of segment index files (input file)

# 03440c883b41446f8d2f58bb8a6d6ed6.flac
# 5.606 8.41
# 19.782 22.65
# 46.47 47.162
# END
# 23ee915f853545c380202d7807f0d026.flac
# 12.006 13.53
# 40.358 42.426
# END

# Example of usage 

# -- Compute stat for every index files in test folder
# compute_stat(test_dir = "test_dir", label_dir = "label_dir")

# -- Compute stat for only one index files
# compute_stat(test_file = "segment_index_file", label_dir = "label_dir")

# Results
# mean IOU, precision, recall


# Merge segments if their distance is less than interval
def merge_seg(seg, interval=0.5):
    if len(seg) < 2:
        return seg

    start, end = seg[0]
    new_seg = []

    idx = 1
    while idx < len(seg):
        if seg[idx][0] - end < interval:
            end = seg[idx][1]
        else:
            new_seg.append((start, end))
            start, end = seg[idx]

        idx += 1
        if idx == len(seg):
            new_seg.append((start, end))

    return new_seg


# Calculate segments' overlapping. If segments are not overlap, return 0.
def cal_overlap(seg1, seg2):
    min1, max1 = seg1
    min2, max2 = seg2
    return max(0, min(max1, max2) - max(min1, min2))


# Calculate union of segments. If segments are not overlap, return lab_start,lab_end
def cal_lab_union(label, vad):
    lab_start, lab_end = label
    vad_start, vad_end = vad
    if min(lab_end, vad_end) - max(lab_start, vad_start) < 0:
        return lab_start, lab_end

    return min(lab_start,vad_start), max(lab_end, vad_end)


# Calculate IOU (intersection over union) of vad and label.
# Input is a list of segments
# Return a list contains IOU of each overlaped-segments
def cal_iou(label_seg, vad_seg):
    iou = []
    lab_idx, vad_idx = 0, 0
    intersect_len = 0.0
    new_seg = True

    while lab_idx < len(label_seg) and vad_idx < len(vad_seg):
        if new_seg:
            # Set region of cuurent union
            union_start, union_end = label_seg[lab_idx]
        
        overlap = cal_overlap(label_seg[lab_idx], vad_seg[vad_idx])

        if overlap > 0:
            intersect_len += overlap
            # Extend current union's region
            new_union = cal_lab_union(label_seg[lab_idx], vad_seg[vad_idx])
            union_start, union_end = cal_lab_union((union_start, union_end), new_union)
            new_seg = False

        elif not new_seg:
            # End of overlapping, Calculate overall IOU of an area
            cur_iou = intersect_len / (union_end - union_start)
            cur_iou = round(cur_iou, 6)
            iou.append(cur_iou)

            new_seg = True
            intersect_len = 0.0

        if label_seg[lab_idx][1] < vad_seg[vad_idx][1]:
            lab_idx += 1
        else:
            vad_idx += 1

    if not new_seg:
        cur_iou = intersect_len / (union_end - union_start)
        cur_iou = round(cur_iou, 6)
        iou.append(cur_iou)

    return iou


# Calculate true positive, true negative, false positive, false negative
def cal_stat(lab_seg, vad_seg, interval=0.010):
    tp = 0 # true positive
    tn = 0 # true negative
    fp = 0 # false positive
    fn = 0 # false negative
    lab_idx, vad_idx = 0, 0
    lab_last, vad_last = 0.0, 0.0
    new_lab = True
    new_vad = True

    while lab_idx < len(lab_seg) and vad_idx < len(vad_seg):
        if new_lab:
            lab_last = lab_seg[lab_idx][0]
        if new_vad:
            vad_last = vad_seg[vad_idx][0]

        overlap = cal_overlap(lab_seg[lab_idx], vad_seg[vad_idx])

        if overlap > 0:
            tp += round(overlap / interval)

            if lab_seg[lab_idx][1] < vad_seg[vad_idx][1]:
                lab_length = lab_seg[lab_idx][1] - lab_last
                fn += round((lab_length - overlap) / interval)
                fp += round( max(0, lab_last - vad_last) / interval )
                
                vad_last = lab_seg[lab_idx][1]
                new_lab = True
                new_vad = False

            else:
                vad_length = vad_seg[vad_idx][1] - vad_last
                fp += round((vad_length - overlap) / interval)
                fn += round( max(0, vad_last - lab_last) / interval )

                lab_last = vad_seg[vad_idx][1]
                new_lab = False
                new_vad = True

        elif lab_seg[lab_idx][0] < vad_seg[vad_idx][0]:
            fn += round((lab_seg[lab_idx][1] - lab_last) / interval)
            new_lab = True
            new_vad = True
        else:
            fp += round((vad_seg[vad_idx][1] - vad_last) / interval)
            new_lab = True
            new_vad = True


        if lab_seg[lab_idx][1] < vad_seg[vad_idx][1]:
            lab_idx += 1
            if lab_idx == len(lab_seg) and not new_vad:
                fp += round((vad_seg[vad_idx][1] - vad_last) / interval)
                vad_idx += 1
        else:
            vad_idx += 1
            if vad_idx == len(vad_seg) and not new_lab:
                fn += round((lab_seg[lab_idx][1] - lab_last) / interval)
                lab_idx += 1

    # Clear the rest
    # if not new_lab:
    #       fn += round((lab_seg[-1][1] - lab_last) / interval)
    #       lab_idx += 1
    # if not new_vad:
    #       fp += round((vad_seg[-1][1] - vad_last) / interval)
    #       vad_idx += 1

    while lab_idx < len(lab_seg):
        fn += round((lab_seg[lab_idx][1] - lab_seg[lab_idx][0]) / interval)
        lab_idx += 1
    
    while vad_idx < len(vad_seg):
        fp += round((vad_seg[vad_idx][1] - vad_seg[vad_idx][0]) / interval)
        vad_idx += 1

    return tp, tn, fp, fn


# Read segment labels from a label file
# Input: label file
# Output: list of speech segments tuple
def read_label(file):
    # label
    lab_segments = []
    # Open label file to get lab_segments
    if(os.path.isfile(file)):
        with open(file) as label_file:
            for line in label_file:
                line = line.strip()
                if line == '':
                    break

                seg = line.split(" ")

                if seg[2] == "lab":
                    seg = [float(seg[0]), float(seg[1])]
                    lab_segments.append(tuple(seg))

    return merge_seg(lab_segments)


# Read segments in a segment index file
# Input: segment index file
# Output: yield (start_time, end_time) of every segment in every audio in an index file
def read_segments(file):
    with open(file, "r") as vad_files:
        while True :
            audio_name = vad_files.readline().strip()
            if audio_name == '':
                break

            audio_name = audio_name.split(".")[0]
            vad_segments = []

            # For every file, read until END to get vad_segments
            while True:
                seg = vad_files.readline().strip()
                if seg == "END":
                    break

                seg = seg.split(" ")
                seg = [float(i) for i in seg]
                vad_segments.append(tuple(seg))

            vad_segments = merge_seg(vad_segments)
            yield(vad_segments, audio_name)


# Calculate overall stat (mean of files)
def vad_stat(test_file, label_dir):
    overall_iou = []
    overall_tp, overall_tn, overall_fp, overall_fn = 0, 0, 0, 0
    for vad_segments, audio_name in read_segments(test_file):     
        lab_segments = read_label(os.path.join(label_dir ,audio_name + ".lab"))
        # Compute IOU
        cur_iou = cal_iou(lab_segments, vad_segments)
        overall_iou.extend(cur_iou)

        # Compute frame-wise stat
        cur_tp, cur_tn, cur_fp, cur_fn = cal_stat(lab_segments, vad_segments)
        overall_tp += cur_tp
        overall_tn += cur_tn
        overall_fp += cur_fp
        overall_fn += cur_fn

    try:
        mean_IOU = round(sum(overall_iou)/len(overall_iou), 4)
    except:
        mean_IOU = None
    try:
        precision = round(overall_tp / (overall_tp + overall_fp),4)
    except:
        precision = None
    try:
        recall = round(overall_tp / (overall_tp + overall_fn),4)
    except:
        recall = None
    return mean_IOU, precision, recall


# Print results
def compute_stat(test_dir = None, label_dir = None, test_file = None):
    if(label_dir == None or (test_dir == None and test_file == None)):
        return

    if(test_file != None):
        try: 
            mean_IOU, precision, recall = vad_stat(test_file, label_dir) 
            print "{}, mean_IOU: {}, Precision: {}, Recall: {}".format(test_file, mean_IOU, precision, recall)
        except:
            print test_file,"has invalid format, skipped"
        return

    for vad_file_name in sorted(os.listdir(test_dir)):
        if(os.path.isdir(os.path.join(test_dir,vad_file_name))):
            print vad_file_name,"is a directory. skipped"
            continue
        try:
            mean_IOU, precision, recall = vad_stat(os.path.join(test_dir, vad_file_name), label_dir)
            print "{}, mean_IOU: {}, Precision: {}, Recall: {}".format(vad_file_name, mean_IOU, precision, recall)
        except:
            print vad_file_name,"has invalid format, skipped"


if __name__ == '__main__':
    if(len(sys.argv) == 4):
        if(sys.argv[1] == 'file'):
            test_file = os.path.join(os.getcwd(), sys.argv[2])
            label_dir = os.path.join(os.getcwd(), sys.argv[3])
            compute_stat(test_file=test_file, label_dir=label_dir)
        else:
            test_dir = os.path.join(os.getcwd(), sys.argv[2])
            label_dir = os.path.join(os.getcwd(), sys.argv[3])
            compute_stat(test_dir=test_dir, label_dir=label_dir)
