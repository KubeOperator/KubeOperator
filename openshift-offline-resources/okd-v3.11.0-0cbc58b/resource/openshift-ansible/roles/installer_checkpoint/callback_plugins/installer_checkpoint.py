"""Ansible callback plugin to print a summary completion status of installation
phases.
"""
from datetime import datetime
from ansible.plugins.callback import CallbackBase
from ansible import constants as C


class CallbackModule(CallbackBase):
    """This callback summarizes installation phase status."""

    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'installer_checkpoint'
    CALLBACK_NEEDS_WHITELIST = False

    def __init__(self):
        super(CallbackModule, self).__init__()

    def v2_playbook_on_stats(self, stats):

        # Return if there are no custom stats to process
        if stats.custom == {}:
            return

        phases = stats.custom['_run']

        # Find the longest phase title
        max_column = 0
        for phase in phases:
            max_column = max(max_column, len(phases[phase].get('title', '')))

        # Sort the phases by start time
        ordered_phases = sorted(phases, key=lambda x: (phases[x].get('start', 0)))

        self._display.banner('INSTALLER STATUS')
        # Display status information for each phase
        for phase in ordered_phases:
            phase_title = phases[phase].get('title', '')
            padding = max_column - len(phase_title) + 2
            phase_status = phases[phase]['status']
            phase_time = phase_time_delta(phases[phase])
            if phase_title:
                self._display.display(
                    '{}{}: {} ({})'.format(phase_title, ' ' * padding, phase_status, phase_time),
                    color=self.phase_color(phase_status))
            # If the phase is not complete, tell the user what playbook to rerun
            if phase_status == 'In Progress' and phase != 'installer_phase_initialize':
                self._display.display(
                    '\tThis phase can be restarted by running: {}'.format(
                        phases[phase]['playbook']))
            # Display any extra messages stored during the phase
            if 'message' in phases[phase]:
                self._display.display(
                    '\t{}'.format(
                        phases[phase]['message']))

    def phase_color(self, status):
        """ Return color code for installer phase"""
        valid_status = [
            'In Progress',
            'Complete',
        ]

        if status not in valid_status:
            self._display.warning('Invalid phase status defined: {}'.format(status))

        if status == 'Complete':
            phase_color = C.COLOR_OK
        elif status == 'In Progress':
            phase_color = C.COLOR_ERROR
        else:
            phase_color = C.COLOR_WARN

        return phase_color


def phase_time_delta(phase):
    """ Calculate the difference between phase start and end times """
    if not phase.get('start'):
        return ''
    time_format = '%Y%m%d%H%M%SZ'
    phase_start = datetime.strptime(phase['start'], time_format)
    if 'end' not in phase:
        # The phase failed so set the end time to now
        phase_end = datetime.now()
    else:
        phase_end = datetime.strptime(phase['end'], time_format)
    delta = str(phase_end - phase_start).split(".")[0]  # Trim microseconds

    return delta
