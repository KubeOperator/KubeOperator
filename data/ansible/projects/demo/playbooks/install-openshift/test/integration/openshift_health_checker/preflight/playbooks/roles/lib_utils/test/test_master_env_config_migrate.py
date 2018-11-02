import os
import sys
import io

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'library'))
sys.path.insert(1, MODULE_PATH)

import master_env_config_migrate  # noqa

INFILE = u"""
t1=Yes
t2=A Space
t3=An\ escaped
t4="A quoted space"
t5="a quoted \\ # in-line comment
  escaped line"
t6 = an unquoted multiline \\
  string
# comment line
"""


def test_read_write_ini():
    infile = io.StringIO(INFILE)

    outfile = io.StringIO()

    config = master_env_config_migrate.SectionlessParser()

    config.readfp(infile)
    config.write(outfile, False)
    print(outfile.getvalue())
    # TODO(michaelgugino): Come up with some clever way to assert the file is
    # correct.

# Contents for t.in:
############################
# t1=This is a spaced string
# t2="Quoted spaced string"
# t3=Escaped\ spaced\ string
# t4 = String
# t5="Quoted multiline string\
#  with escaped newline"
# t6=escaped\ spaced\\
#  multiline\ unquoted.


if __name__ == '__main__':
    test_read_write_ini()

    with open('t.in') as f:
        config = master_env_config_migrate.SectionlessParser()
        config.readfp(f)
    with open('t2.out', 'w') as f:
        config.write(f, False)
