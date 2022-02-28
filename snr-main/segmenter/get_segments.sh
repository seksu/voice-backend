#/bin/bash
# set -x # echo on
if [ $# -lt 2 ]; then
    echo "usage: $0 path out_dir [OPTIONS]"
    exit
fi

mkdir -p $2 || exit 1

if [ -d $1 ]; then
    out=$2/segments_`basename $1`
    rm -f $out
    for f in $1/*.wav; do
        echo $f
        python3 n_segmenter.py $f ${@:3} | python3 format_nsegment_output.py $f $2 >> $out || exit 1
    done
else
    echo $1
    python3 n_segmenter.py $1 ${@:3} | python3 format_nsegment_output.py $1 $2 > $2/segments_`basename ${1%.*}` || exit 1
fi