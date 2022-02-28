#!/bin/bash
if [ $# -lt 3 ]; then
    echo "usage: $0 segments_file wav_path out_dir"
    exit
fi
wav_path=$2
out_dir=$3

while IFS= read -r line; do
    t=( $line )
    segment=${t[0]}
    rec=${t[1]}
    st=${t[2]}
    en=${t[3]}    
    echo $rec $st $en
    mkdir -p "$out_dir/$rec"
    sox "${wav_path}/${rec}.wav" "$out_dir/${rec}/${segment}.wav" trim $st =$en

done < "$1" # segments file
