import random
import os
import string

import numpy as np

file = open("data.test", "w")
file.write("load 1s\n")
for idx in range(20):
    file.write("    http_requests{job=\"app-server\", instance=\"" + str(idx) + "\", group=\"canary\", severity=\"overwrite-me\"} ")
    for i in range(10000000):
        vector = np.random.random((1, )) * 100
        file.write(str(vector[0]) + " ")
    file.write("\n")
file.close()