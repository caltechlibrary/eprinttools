
import sys
from datetime import tzinfo, timedelta, datetime

class Logger:
    def __init__(self, pid, time_format = '%Y/%m/%d %H:%M:%S'):
        self.pid = pid
        self.time_format = time_format

    def sprint(self, msg):
        dt = datetime.now().strftime(self.time_format)
        pid = self.pid
        return (f'{dt} (pid: {pid}) {msg}')

    def print(self, msg):
        print(self.sprint(msg), flush = True)

    def fatal(self, msg):
        print(self.sprint(msg), flush = True)
        sys.exit(1)

