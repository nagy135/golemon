import sys

try:
    curr=int(sys.argv[1])*60+int(sys.argv[2])
    total=int(sys.argv[3])*60+int(sys.argv[4])
    print(int(curr/total*100))
except:
    print(0)
