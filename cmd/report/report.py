# Take results data, and transform it into the format needed for the report output
import os
import shutil
import csv

# just to make less typing
join = os.path.join

report_path = join(os.getcwd(), 'output', 'reportResults')
input_path = join(os.getcwd(), 'output', 'results')

if "cmd" in os.getcwd():
    print("Please run from top level of the ballot validator repo")

if not os.path.exists(report_path):
    os.mkdir(report_path)


# now copy some of the files as is
shutil.copy2(join(input_path, "president.txt"), report_path)
shutil.copy2(join(input_path, "sfc-chair.txt"), report_path)

# format some of the files so they are ready to be used in the report
for file in ["undergradSenate.csv", "graduateSenate.csv", "sfc-at-large.csv"]:
    with open(join(input_path, file)) as input_file, \
            open(join(report_path, file), "w") as output_csv_file, \
            open(join(report_path, f"{file}.txt"), "w") as output_txt_file:
        input_csv = csv.reader(input_file)
        output_csv = csv.writer(output_csv_file)

        candidate_list = []
        for row_num, row in enumerate(input_csv):
            if row_num == 0:
                row.append("Status")
            else:
                row.append("")
                candidate_list.append(row[0])
            output_csv.writerow(row)
        
        # now write the candidate list in reverse vote order, for graph insertion
        candidate_list.reverse()
        output_txt_file.write(",".join(candidate_list) + ",")
