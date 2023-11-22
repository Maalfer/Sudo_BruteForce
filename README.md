# Sudo_BruteForce
**Script hecho en bash para realizar un ataque de fuerza bruta a un usuario de un sistema Linux**

Se ejecuta de la siguiente forma:

```bash
bash Linux-Su-Force.sh mario /usr/share/wordlists/rockyou.txt
```
![image](https://github.com/Maalfer/Sudo_BruteForce/assets/96432001/8fb151eb-4e87-4521-9dc2-db8ba9a5e41a)

En caso de encontrarse la contraseña dentro del diccionario, se limpiará la pantalla y se mostrará dicha contraseña:

![image](https://github.com/Maalfer/Sudo_BruteForce/assets/96432001/5cd106b8-cdd3-4d96-866a-8e435ed30f50)

Comprobamos su correcto funcionamiento:

![image](https://github.com/Maalfer/Sudo_BruteForce/assets/96432001/33deb32e-54a6-4a18-9ac5-6d646f3e6e22)


Comprobamos también el correcto funcionamiento de las excepciones, donde en caso de haber un error al proporcionar los parámetros, se imprimen unas instrucciones:

![image](https://github.com/Maalfer/Sudo_BruteForce/assets/96432001/8b3e11b9-314d-45b5-8fdc-04f1e2867f95)
