import os
import pytest
import subprocess
import csv
import collections

TESTS_DIR = "./tests/pattern_tests/"
def gather_testcases(project_dir):
    testfolders = [os.path.join(TESTS_DIR, f) for f in os.listdir(TESTS_DIR) if os.path.isdir(os.path.join(project_dir, f))]
    return testfolders

@pytest.mark.parametrize("folder", gather_testcases(TESTS_DIR))
def test_single(folder):
    ast_json = os.path.join(folder, "ast.json")
    expected_output_file  = os.path.join(folder, "facts_out/patternMatch.csv")


    expected_out_lines = []
    with open(expected_output_file) as csv_file:
        csv_reader = csv.reader(csv_file, delimiter="\t")
        expected_out_lines = [row for row in csv_reader]


    output_folder = f"/tmp/souffle_test_{os.path.basename(folder)}"
    if not os.path.isdir(output_folder):
        os.mkdir(output_folder)

    cmd = ["go", "run", "cmd/scilla_static/main.go", "-analysis_dir", output_folder, ast_json]
    cmd = " ".join(cmd)
    process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE)
    stdout, err = process.communicate()
    stdout = stdout.decode("utf-8") if stdout is not None else ""

    assert err is None
    output_file = os.path.join(output_folder, "facts_out/patternMatch.csv")

    out_lines = []
    with open(output_file) as csv_file:
        csv_reader = csv.reader(csv_file, delimiter="\t")
        out_lines = [row for row in csv_reader]


    expected_violations = collections.Counter([r[0] for r in expected_out_lines if r[2] == 'violation'])
    expected_warnings = collections.Counter([r[0] for r in expected_out_lines if r[2] == 'warning'])

    violations = collections.Counter([r[0] for r in out_lines if r[2] == 'violation'])
    warnings = collections.Counter([r[0] for r in out_lines if r[2] == 'warning'])

    for key in violations.keys():
        assert key in expected_violations
        assert violations[key] >= expected_violations[key]

    for key in warnings.keys():
        assert key in expected_warnings
        assert warnings[key] >= expected_warnings[key]
