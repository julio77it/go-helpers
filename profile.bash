#!/bin/bash

echo ""
echo "to measure Heap consuptions in db.Rows.GetFields"
echo "(pprof) list Scan"
echo "to measure Heap consuptions in SQLRows.GetFields"
echo "(pprof) list GetFieds"
echo "to measure Heap consuptions in SQLRows.GetByIndex"
echo "(pprof) list GetByIndex"
echo "to measure Heap consuptions in SQLRows.GetByName"
echo "(pprof) list GetByName"
echo ""

go tool pprof -alloc_space mem.out
