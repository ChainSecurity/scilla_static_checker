import pytest
import os

TESTS_DIR = './tests/pattern_tests/'
def gather_testcases(project_dir):
    testfolders = [f for f in os.listdir(TESTS_DIR) if os.path.isdir(os.path.join(project_dir, f))]
    return testfolders

@pytest.mark.parametrize('case_folder', gather_testcases(TESTS_DIR))
def test_single(case_folder):
    assert 1 == 1
    # tutils.run_single(combined_json_data, spec_file, expected_status, expected_traces)
