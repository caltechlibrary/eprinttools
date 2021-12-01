
class TestSet:
    '''TestSet provides a lightweight test manager for a collection
       of simple boolean tests'''
    def __init__(self, name):
        self.name = name
        self.success, self.errors = 0, 0
        self.tests = []

    def add(self, fn):
        '''add expects a boolean function to run as a test'''
        self.tests.append(fn)
    
    def run(self):
        for _, fn in enumerate(self.tests):
            if fn():
                self.success += 1
            else:
                self.errors += 1
        if self.errors == 0:
            print(f'OK, {self.name}')
            return True
        else:
            tot = self.success + self.errors
            print(f'''
{self.name}, failed.
=================================

Success: {self.success}
 Errors: {self.errors}
---------------------------------
  Total: {tot}

''')
            return False

