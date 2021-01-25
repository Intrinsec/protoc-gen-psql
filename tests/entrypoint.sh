#!/bin/sh


for f in `ls /sql/*.pb.psql`; do
    echo "Executing $f"
    psql -v "ON_ERROR_STOP=1" -U postgres -h db -p 5432 -a -f $f
done
