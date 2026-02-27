# 🍔 Burger System — Sistema de Gestión de Restaurante

## Datos del Proyecto

- **Materia:** Programación orientada a objetos con GO
- **Universidad:** Universidad Internacional del Ecuador (UIDE)
- **Fecha:** Febrero 2026
- **Integrantes:** Luis Agapito 

## Objetivo

Desarrollar un sistema de gestión para un restaurante de hamburguesas que permita administrar el menú de productos y las órdenes de clientes mediante una API REST con servicios web JSON, aplicando los conceptos de las 4 unidades de la materia: programación funcional, POO, encapsulación, interfaces, manejo de errores y servicios web.

## Tecnologías Utilizadas

| Tecnología | Uso |
|---|---|
| **Go (Golang)** | Backend, API REST, lógica de negocio |
| **MySQL** | Base de datos relacional (XAMPP) |
| **HTML + Bootstrap 5** | Frontend web responsive |
| **JavaScript (vanilla)** | Comunicación con la API mediante fetch() |
| **Gorilla Mux** | Enrutador HTTP para Go |
| **go-sql-driver/mysql** | Driver MySQL para Go |

## Estructura del Proyecto

```
burger-system/
├── main.go                    # Punto de entrada, rutas y servidor HTTP
├── go.mod                     # Módulo Go con dependencias
├── go.sum                     # Checksums de dependencias
├── api/
│   └── handlers.go            # Handlers HTTP para la API REST
├── storage/
│   ├── database.go            # Capa de persistencia MySQL (CRUD)
│   └── store.go               # Almacenamiento en memoria (referencia)
├── menu/
│   └── hamburguesa.go         # Modelo Hamburguesa con encapsulación
├── orders/
│   └── orden.go               # Modelo Orden con cálculo de IVA 15%
├── utils/
│   ├── interfaces.go          # Interfaces: Identificable, Informable, Validable
│   └── helpers.go             # Funciones de polimorfismo con interfaces
├── templates/
│   └── index.html             # Frontend completo (SPA con Bootstrap)
├── static/
│   └── img/                   # Imágenes de hamburguesas
│       ├── clasica.png
│       ├── premium.png
│       ├── vegetariana.png
│       ├── mexicana.png
│       ├── hawaiana.png
│       ├── chili.png
│       └── portobelo.png
└── README.md                  # Este archivo
```

## Funcionalidades Principales

### Gestión de Hamburguesas
- Crear hamburguesas con nombre, descripción, precio, categoría e ingredientes
- Listar todas las hamburguesas con imágenes, precios y badges de categoría
- Editar hamburguesas existentes (nombre, descripción, precio)
- Eliminar hamburguesas con protección de integridad referencial (no permite eliminar si tiene órdenes asociadas)

### Gestión de Órdenes
- Crear órdenes seleccionando hamburguesas del menú con cantidades y notas
- Soporte para 3 tipos de orden: mesa, para llevar y delivery
- Calcular subtotal, IVA 15% y total automáticamente
- Listar todas las órdenes con estadísticas de resumen (total órdenes, ingresos, pendientes)
- Ver detalle completo de cada orden
- Cambiar estado de órdenes: PENDIENTE → EN_COCINA → LISTA → ENTREGADA / CANCELADA
- Eliminar órdenes (elimina también sus items asociados)

## Servicios Web (Endpoints API REST)

La API expone **10 servicios web** con serialización JSON. Todos responden en formato `application/json`.

### Hamburguesas (5 endpoints)

| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/api/hamburguesas` | Listar todas las hamburguesas |
| `GET` | `/api/hamburguesas/{id}` | Obtener hamburguesa por ID |
| `POST` | `/api/hamburguesas` | Crear nueva hamburguesa |
| `PUT` | `/api/hamburguesas/{id}` | Actualizar hamburguesa existente |
| `DELETE` | `/api/hamburguesas/{id}` | Eliminar hamburguesa |

### Órdenes (5 endpoints)

| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/api/ordenes` | Listar todas las órdenes |
| `GET` | `/api/ordenes/{id}` | Obtener orden con sus items |
| `POST` | `/api/ordenes` | Crear nueva orden con items |
| `PUT` | `/api/ordenes/{id}/estado` | Actualizar estado de una orden |
| `DELETE` | `/api/ordenes/{id}` | Eliminar orden y sus items |

### Ejemplos de Uso

**Crear hamburguesa:**
```bash
curl -X POST http://localhost:8080/api/hamburguesas \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Hamburguesa BBQ","descripcion":"Con salsa BBQ","precio":12.99,"categoria":"Premium","ingredientes":["carne","queso","bbq"]}'
```

**Crear orden:**
```bash
curl -X POST http://localhost:8080/api/ordenes \
  -H "Content-Type: application/json" \
  -d '{"tipoOrden":"mesa","numeroMesa":5,"items":[{"hamburguesaID":"BURG-001","cantidad":2,"notas":"sin cebolla"}]}'
```

**Cambiar estado:**
```bash
curl -X PUT http://localhost:8080/api/ordenes/ORD-001/estado \
  -H "Content-Type: application/json" \
  -d '{"estado":"EN_COCINA"}'
```

## Conceptos de la Materia Aplicados

### Unidad 1 — Estructura y Programación Funcional
- Organización modular en paquetes (`menu`, `orders`, `api`, `storage`, `utils`)
- Funciones puras para validación y cálculo de totales
- Funciones de orden superior en helpers.go

### Unidad 2 — Programación Orientada a Objetos
- Structs como clases: `Hamburguesa`, `Orden`, `ItemOrden`, `Server`, `Database`
- Métodos asociados a structs (receptores de puntero)
- Composición de tipos para extender funcionalidad

### Unidad 3 — Encapsulación, Interfaces y Manejo de Errores
- **Encapsulación:** campos privados `id` y `fechaCreado` en Hamburguesa, accesibles solo por getters/setters
- **Interfaces:** `Identificable`, `Informable`, `Validable`, `Actualizable` definen contratos que implementan Hamburguesa y Orden
- **Polimorfismo:** funciones `MostrarInformacion()`, `ValidarEntidad()` aceptan cualquier tipo que implemente la interfaz
- **Manejo de errores:** validación en todos los niveles (modelos, handlers, base de datos), mensajes descriptivos al usuario

### Unidad 4 — Servicios Web
- API REST con 10 endpoints usando Gorilla Mux
- Serialización/deserialización JSON con `encoding/json`
- `MarshalJSON` personalizado para incluir campos privados en la respuesta
- CORS middleware para comunicación frontend-backend
- Códigos HTTP apropiados: 200 OK, 201 Created, 400 Bad Request, 404 Not Found, 409 Conflict

## Paquetes de Terceros

| Paquete | Versión | Uso | Documentación |
|---|---|---|---|
| `github.com/gorilla/mux` | v1.8.1 | Router HTTP con variables de ruta `{id}` | https://github.com/gorilla/mux |
| `github.com/go-sql-driver/mysql` | v1.7.1 | Driver MySQL para `database/sql` | https://github.com/go-sql-driver/mysql |

## Instalación y Ejecución

### Prerrequisitos
- Go 1.21 o superior
- XAMPP con MySQL activo
- Git

### Pasos

1. **Clonar el repositorio:**
```bash
git clone https://github.com/perluis/burger-system.git
cd burger-system
```

2. **Crear la base de datos en phpMyAdmin:**
```sql
CREATE DATABASE burger_system;
USE burger_system;

CREATE TABLE hamburguesas (
    id VARCHAR(10) PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    precio DECIMAL(10,2) NOT NULL,
    categoria VARCHAR(20) NOT NULL,
    ingredientes TEXT,
    disponible BOOLEAN DEFAULT TRUE
);

CREATE TABLE ordenes (
    id VARCHAR(10) PRIMARY KEY,
    numero_mesa INT DEFAULT 0,
    tipo_orden VARCHAR(20) NOT NULL,
    estado VARCHAR(20) DEFAULT 'PENDIENTE',
    subtotal DECIMAL(10,2) DEFAULT 0,
    total DECIMAL(10,2) DEFAULT 0,
    fecha_hora TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items_orden (
    id INT AUTO_INCREMENT PRIMARY KEY,
    orden_id VARCHAR(10) NOT NULL,
    hamburguesa_id VARCHAR(10) NOT NULL,
    nombre VARCHAR(100),
    cantidad INT DEFAULT 1,
    precio_unit DECIMAL(10,2),
    notas TEXT,
    FOREIGN KEY (orden_id) REFERENCES ordenes(id),
    FOREIGN KEY (hamburguesa_id) REFERENCES hamburguesas(id)
);

-- Datos iniciales
INSERT INTO hamburguesas VALUES
('BURG-001', 'Hamburguesa Clásica', 'Deliciosa hamburguesa tradicional', 8.99, 'Clasica', 'carne,queso,lechuga,tomate,pan', TRUE),
('BURG-002', 'BBQ Premium', 'Con tocino y salsa BBQ casera', 12.99, 'Premium', 'carne,queso cheddar,tocino,bbq,cebolla,pan', TRUE),
('BURG-003', 'Veggie Deluxe', 'Hamburguesa vegetariana gourmet', 10.99, 'Vegetariana', 'portobello,queso,lechuga,tomate,aguacate,pan integral', TRUE),
('BURG-004', 'Hamburguesa Mexicana', 'Con jalapeños y guacamole', 11.99, 'Premium', 'carne,queso,jalapeño,guacamole,pico de gallo,pan', TRUE);
```

3. **Instalar dependencias:**
```bash
go mod tidy
```

4. **Ejecutar el servidor:**
```bash
go run main.go
```

5. **Abrir en el navegador:**
```
http://localhost:8080
```

## Base de Datos — Diagrama Relacional

```
hamburguesas (1) ──────── (N) items_orden (N) ──────── (1) ordenes
     id ◄─── FK ──── hamburguesa_id    orden_id ──── FK ──► id
```

- Una hamburguesa puede estar en múltiples items de orden
- Una orden puede tener múltiples items
- Eliminar una hamburguesa con órdenes asociadas está protegido por foreign key

## Zona Horaria

El sistema usa la zona horaria de Ecuador (`America/Guayaquil` UTC-5) configurada en el DSN de MySQL para que las fechas de las órdenes se registren correctamente.