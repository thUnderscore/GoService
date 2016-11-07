start /w /MIN cmd /c sh exec.sh %* ^>Build.Log.txt ^2^>Build.Error.Log.txt
rem ^&^1
type Build.Log.txt
type Build.Error.Log.txt > ^&^2
rm Build.Log.txt
rm Build.Error.Log.txt