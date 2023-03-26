import sys 
import string
import random
from typing import Dict, List

def generate_instruments(n_instr):
    instruments = set() 

    while len(instruments) < n_instr:
        builder = [] 
        for j in range(4):
            builder.append(random.choice(string.ascii_uppercase))
        instrument = "".join(builder) 
        instruments.add(instrument)

    return list(instruments)
 
ID_COUNTER = 0

def create_order(n, side, thread_history: Dict[int, List[int]], instruments):
    global ID_COUNTER
    ID_COUNTER += 1

    thread = random.randint(0, n - 1)  
    instrument = random.choice(instruments) 
    price = random.randint(1, 100) 
    count = random.randint(1, 100) 
    s_builder = [str(thread), side, str(ID_COUNTER), instrument, str(price), str(count)] 
    if thread not in thread_history:
        thread_history[thread] = [] 
    thread_history[thread].append(ID_COUNTER)
    return " ".join(s_builder) 

def create_cancel(n, thread_history: Dict[int, List[int]]):
    thread = random.choice(list(thread_history.keys()))
    cancel_id = random.choice(thread_history[thread]) 

    s_builder = [str(thread), 'C', str(cancel_id)] 
    return " ".join(s_builder)


def main():
    args = sys.argv[1:] 
    if (len(args) < 5):
        print("Usage: python3 sanic_v2.py [n] [b] [s] [c] [n_instr]")
        sys.exit(1) 

    n, b, s, c, n_instr = list(map(int, args)) 
    if n <= 0 or b <= 0 or s <= 0 or n_instr <= 0:
        print("n, b, s, c must be greater than 0")
        sys.exit(1) 

    available = ['B'] * b + ['S'] * s + ['C'] * c 

    commands = []
    thread_history = {}
    instruments = generate_instruments(n_instr)

    while available: 
        index = random.randint(0, len(available) - 1)
        available[index], available[-1] =  available[-1], available[index]

        a = available.pop() 
        
        if a != 'C':
            commands.append(create_order(n, a, thread_history, instruments))
        elif not thread_history:
            continue  
        else: 
            commands.append(create_cancel(n, thread_history))
    
    f = open("input.in", 'w') 
    f.write(str(n) + '\n') 
    f.write('o' + '\n')
    for command in commands:
        f.write(command + '\n') 
    f.write('x' + '\n') 
    f.close() 


if __name__ == "__main__":
    main() 
    