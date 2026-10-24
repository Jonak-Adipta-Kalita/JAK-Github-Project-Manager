from decimal import Decimal
import sys


def count_decimals(f):
    s = str(f)
    if 'e' in s or 'E' in s:
        return 0
    return len(s.split('.')[-1])


def max_decimals(a, d):
    return max(count_decimals(Decimal(str(a))), count_decimals(Decimal(str(d))))


def get_aps():
    a = float(sys.argv[1]) if len(
        sys.argv) > 1 else float(input("First Term: "))
    d = float(sys.argv[2]) if len(sys.argv) > 1 else float(
        input("Common Difference: "))
    n = int(sys.argv[3]) if len(
        sys.argv) > 1 else int(input("No. of Terms: "))

    decimals = max_decimals(a, d)

    aps = []
    curr_term = a

    for _ in range(1, n + 1):
        aps.append(curr_term)
        curr_term += d

    return [round(ap, decimals) for ap in aps]


if __name__ == "__main__":
    aps = get_aps()
    aps_list = ", ".join([str(ap) for ap in aps])

    if not len(sys.argv) > 1:
        print()
    print(f"AP: {aps_list}" if not len(sys.argv) > 1 else aps_list)
