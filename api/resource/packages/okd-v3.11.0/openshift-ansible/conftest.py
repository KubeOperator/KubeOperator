"""pytest configuration"""


def pytest_ignore_collect(path):
    """Hook to ignore symlink files and directories."""
    return path.islink()
