from django.dispatch import Signal, receiver
from django.utils import timezone

pre_script_execution_start = Signal(providing_args=('execution',))
post_script_execution_start = Signal(providing_args=('execution', 'result'))
django_ready = Signal()


@receiver(pre_script_execution_start)
def on_execution_start(sender, execution, **kwargs):
    execution.date_start = timezone.now()
    execution.state = execution.STATE_STARTED
    execution.save()


@receiver(post_script_execution_start)
def on_execution_end(sender, execution, result, ignore_errors, **kwargs):
    date_finished = timezone.now()
    timedelta = (timezone.now() - execution.date_start).seconds
    if result.get('summary', {}).get("success", False):
        state = execution.STATE_SUCCESS
    else:
        state = execution.STATE_FAILURE
    execution.result_summary = result.get('summary', {})
    execution.result_raw = result.get('raw', {})
    execution.state = state
    execution.date_end = date_finished
    execution.timedelta = timedelta
    execution.save()
