
class T:
    '''T provides a set of asserts and keeps track of a failure state'''
    ok = True
    def Expected(self, expected, got, message):
        '''Compare two variables for type and value equality'''
        if type(expected) != type(got):
            print(f'expected type {type(expected)}, got type {type(got)}, {message}')
            self.ok = False
        elif expected != got:
            print(f'expected "{expected}", got "{got}", {message}')
            self.ok = False

    def ExpectedList(self, expected, got, message):
        '''Compare two Python Lists for type and value equality'''
        if not isinstance(expected, list):
            print(f'expected a list but found {type(expected)}, {message}')
            self.ok = False
        if not isinstance(got, list):
            print(f'test value should be list, found {type(got)}, {message}')
            self.ok = False
        if len(expected) != len(got):
            print(f'expected length {len(expected)}, got {len(got)}, {message}')
            self.ok = False
        if self.ok :
            for i, o in enumerate(expected):
                if o != got[i]:
                    print(f'expected ({i}) "{o}", got "{got[i]}", {message}')
                    self.ok = False

    def Results(self):
        return self.ok
