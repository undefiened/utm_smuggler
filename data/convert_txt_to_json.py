import csv
import simplejson

filename = '50x50.txt'
data = {}

with open(filename) as csv_file:
    csv_reader = csv.reader(csv_file, delimiter=' ')
    heights = []

    for row in csv_reader:
        heights.append([float(x) for x in row])

    data = {'heights': heights}
    print(data)

    with open(filename.replace('.txt', '.json'), 'w') as json_file:
        simplejson.dump(data, json_file)