import series
import sys

aps = series.get_aps()

if not len(sys.argv) > 1:
    print()

print(f"Sum: {sum(aps)}" if not len(sys.argv) > 1 else sum(aps))
