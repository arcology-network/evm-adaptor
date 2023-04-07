from solcx import compile_files
import sys
import os

def read_contracts(dir):
    sources = []
    for root, _, files in os.walk(dir):
        for file in files:
            if file.endswith('.sol'):
                sources.append(os.path.join(root, file))
    return sources

def main():
    fileOrDir = sys.argv[1]
    contractName = sys.argv[2]
    
    compiled_sol = {}
    files = [fileOrDir]
    if os.path.isdir(fileOrDir):
        files = read_contracts(fileOrDir)   
 
    compiled_sol = compile_files(files, output_values = ['bin', 'abi'])
    contract = compiled_sol[fileOrDir + ':' + contractName]
    print(contract['bin'], end='')

if __name__ == '__main__':
    main()