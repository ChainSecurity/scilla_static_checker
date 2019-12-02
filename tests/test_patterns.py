import os
import pytest
import subprocess

TESTS_DIR = './tests/pattern_tests/'
def gather_testcases(project_dir):
    testfolders = [os.path.join(TESTS_DIR, f) for f in os.listdir(TESTS_DIR) if os.path.isdir(os.path.join(project_dir, f))]
    return testfolders

@pytest.mark.parametrize('folder', gather_testcases(TESTS_DIR))
def test_single(folder):
    ast_json = os.path.join(folder, "ast.json")
    expected_output_file  = os.path.join(folder, "facts_out/patternMatch.csv")
    output_folder = f"/tmp/souffle_test_{os.path.basename(folder)}"
    if not os.path.isdir(output_folder):
        os.mkdir(output_folder)

    cmd = ["go", "run", "cmd/scilla_static/main.go", "-analysis_dir", output_folder, ast_json]
    cmd = " ".join(cmd)
    process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE)
    stdout, err = process.communicate()
    stdout = stdout.decode('utf-8') if stdout is not None else ''

    output_file = os.path.join(output_folder, "facts_out/patternMatch.csv")
    print(output_file)
    print(stdout)
    assert os.path.isfile(output_file)



    assert err is None



    # tutils.run_single(combined_json_data, spec_file, expected_status, expected_traces)
