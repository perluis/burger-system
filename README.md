# 🍔 Sistema de Gestión - Restaurante de Hamburguesas

Sistema completo de gestión para restaurante desarrollado en Go con API REST, MySQL y Frontend web.

**Autor:** Luis Agapito Pérez  
**Universidad:** UIDE 2026  
**Curso:** Programación Orientada a Objetos  
**Proyecto:** Autónomo 3 - Servicios Web

---

## 📋 Descripción

Sistema integral para la gestión de un restaurante de hamburguesas que incluye:

- **API REST** con 8 endpoints funcionales
- **Base de datos MySQL** para persistencia de datos
- **Frontend web** con Bootstrap para interfaz visual
- **Serialización JSON** para comunicación cliente-servidor
- **Cálculo automático** de totales con IVA 12%

---

## 🛠️ Tecnologías Utilizadas

| Tecnología | Uso |
|------------|-----|
| **Go 1.24** | Backend y API REST |
| **Gorilla Mux** | Router HTTP |
| **MySQL/MariaDB** | Base de datos |
| **HTML5** | Estructura web |
| **Bootstrap 5** | Estilos y diseño responsivo |
| **JavaScript** | Interactividad frontend |
| **XAMPP** | Servidor MySQL local |

---

## 📁 Estructura del Proyecto
```
burger-system/
├── main.go                 # Servidor HTTP principal
├── api/
│   └── handlers.go         # Controladores de la API
├── storage/
│   ├── memory.go           # Storage en memoria (legacy)
│   └── database.go         # Conexión y operaciones MySQL
├── menu/
│   └── hamburguesa.go      # Modelo de hamburguesas
├── orders/
│   └── orden.go            # Modelo de órdenes
├── templates/
│   └── index.html          # Página web principal
├── static/
│   └── img/                # Imágenes de hamburguesas
├── utils/
│   ├── interfaces.go       # Interfaces del sistema
│   └── helpers.go          # Funciones auxiliares
├── go.mod                  # Dependencias Go
└── README.md               # Este archivo
```

---

## 🚀 Instalación y Ejecución

### Prerrequisitos

- Go 1.21 o superior
- XAMPP (MySQL/MariaDB)
- Git

### Paso 1: Clonar repositorio
```bash
git clone https://github.com/perluis/burger-system.git
cd burger-system
```

### Paso 2: Instalar dependencias
```bash
go mod download
```

### Paso 3: Configurar Base de Datos

1. Iniciar XAMPP (Apache + MySQL)
2. Abrir phpMyAdmin: http://localhost/phpmyadmin
3. Crear base de datos: `burger_system`
4. Ejecutar el siguiente SQL:
```sql
-- Tabla de hamburguesas
CREATE TABLE hamburguesas (
    id VARCHAR(20) PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    precio DECIMAL(10,2) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    ingredientes TEXT,
    disponible BOOLEAN DEFAULT TRUE,
    fecha_creado DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de órdenes
CREATE TABLE ordenes (
    id VARCHAR(20) PRIMARY KEY,
    numero_mesa INT,
    tipo_orden VARCHAR(50) NOT NULL,
    estado VARCHAR(50) DEFAULT 'PENDIENTE',
    subtotal DECIMAL(10,2),
    total DECIMAL(10,2),
    fecha_hora DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de items de orden
CREATE TABLE items_orden (
    id INT AUTO_INCREMENT PRIMARY KEY,
    orden_id VARCHAR(20) NOT NULL,
    hamburguesa_id VARCHAR(20) NOT NULL,
    nombre VARCHAR(100),
    cantidad INT NOT NULL,
    precio_unit DECIMAL(10,2),
    notas TEXT,
    FOREIGN KEY (orden_id) REFERENCES ordenes(id),
    FOREIGN KEY (hamburguesa_id) REFERENCES hamburguesas(id)
);

-- Datos iniciales
INSERT INTO hamburguesas (id, nombre, descripcion, precio, categoria, ingredientes, disponible) VALUES
('BURG-001', 'Hamburguesa Clásica', 'Deliciosa hamburguesa tradicional', 8.99, 'Clasica', 'carne,queso,lechuga,tomate,pan', TRUE),
('BURG-002', 'BBQ Premium', 'Con tocino y salsa BBQ casera', 12.99, 'Premium', 'carne,queso cheddar,tocino,bbq,cebolla,pan', TRUE),
('BURG-003', 'Veggie Deluxe', 'Hamburguesa vegetariana gourmet', 10.99, 'Vegetariana', 'portobello,queso,lechuga,tomate,aguacate,pan integral', TRUE);
```

### Paso 4: Ejecutar el servidor
```bash
go run main.go
```

### Paso 5: Abrir en navegador

- **Web:** http://localhost:8080
- **API:** http://localhost:8080/api/hamburguesas

---

## 📡 API Endpoints

### Hamburguesas

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/hamburguesas` | Listar todas las hamburguesas |
| `GET` | `/api/hamburguesas/{id}` | Obtener hamburguesa por ID |
| `POST` | `/api/hamburguesas` | Crear nueva hamburguesa |
| `PUT` | `/api/hamburguesas/{id}` | Actualizar hamburguesa |
| `DELETE` | `/api/hamburguesas/{id}` | Eliminar hamburguesa |

### Órdenes

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/ordenes` | Crear nueva orden |
| `GET` | `/api/ordenes/{id}` | Obtener orden por ID |
| `PUT` | `/api/ordenes/{id}/estado` | Actualizar estado de orden |

---

## 📖 Ejemplos de Uso (CURL)

### Listar hamburguesas
```bash
curl http://localhost:8080/api/hamburguesas
```

### Crear hamburguesa
```bash
curl -X POST http://localhost:8080/api/hamburguesas \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Hamburguesa Especial",
    "descripcion": "Con ingredientes secretos",
    "precio": 13.99,
    "categoria": "Premium",
    "ingredientes": ["carne", "queso", "bacon", "huevo"]
  }'
```

### Crear orden
```bash
curl -X POST http://localhost:8080/api/ordenes \
  -H "Content-Type: application/json" \
  -d '{
    "tipoOrden": "mesa",
    "numeroMesa": 5,
    "items": [
      {"hamburguesaID": "BURG-001", "cantidad": 2, "notas": "sin cebolla"},
      {"hamburguesaID": "BURG-002", "cantidad": 1, "notas": ""}
    ]
  }'
```

### Actualizar estado de orden
```bash
curl -X PUT http://localhost:8080/api/ordenes/ORD-001/estado \
  -H "Content-Type: application/json" \
  -d '{"estado": "EN_COCINA"}'
```

---

## 🏗️ Arquitectura del Sistema

### Principios Aplicados

- **Encapsulación:** Campos privados con getters/setters
- **Interfaces:** Polimorfismo para operaciones comunes
- **Manejo de Errores:** Validaciones en todas las operaciones
- **Separación de Responsabilidades:** Paquetes independientes
- **MVC:** Modelo-Vista-Controlador adaptado

### Flujo de Datos
```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Cliente │ ──► │  Router  │ ──► │ Handler  │ ──► │  MySQL   │
│ (Browser)│ ◄── │  (Mux)   │ ◄── │  (API)   │ ◄── │   (BD)   │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
     │                                                   │
     │              JSON Request/Response                │
     └───────────────────────────────────────────────────┘
```

---

## 🔧 Funcionalidades

### Módulo de Hamburguesas
- ✅ CRUD completo (Crear, Leer, Actualizar, Eliminar)
- ✅ Validación de datos (precio positivo, categoría válida)
- ✅ Generación automática de IDs
- ✅ Control de disponibilidad
- ✅ Persistencia en MySQL

### Módulo de Órdenes
- ✅ Creación de órdenes multi-item
- ✅ Cálculo automático de subtotales
- ✅ Aplicación de IVA (12% Ecuador)
- ✅ Gestión de estados (PENDIENTE → EN_COCINA → LISTA → ENTREGADA)
- ✅ Notas personalizadas por item

### Frontend Web
- ✅ Diseño responsivo con Bootstrap
- ✅ Visualización de menú con imágenes
- ✅ Formulario para crear hamburguesas
- ✅ Formulario para crear órdenes
- ✅ Actualización dinámica sin recargar página

---

## 🎯 Estados de Orden

| Estado | Descripción |
|--------|-------------|
| `PENDIENTE` | Orden recién creada |
| `EN_COCINA` | En preparación |
| `LISTA` | Lista para entregar |
| `ENTREGADA` | Entregada al cliente |
| `CANCELADA` | Orden cancelada |

---

## 🔮 Proyección Futura

- Autenticación de usuarios (login/registro)
- Panel de administración
- Reportes de ventas
- Sistema de inventario
- Integración con sistemas de pago
- Aplicación móvil con Flutter
- Notificaciones en tiempo real

---

## 📝 Notas Técnicas

- El servidor corre en el puerto 8080
- La conexión MySQL usa: `root@localhost:3306/burger_system`
- El IVA aplicado es del 12% (Ecuador)
- Las imágenes se sirven desde `/static/img/`

---

## 👨‍💻 Autor

**Luis Agapito Pérez**  
Universidad Internacional del Ecuador (UIDE)  
Programación Orientada a Objetos - 2026

---

## 📄 Licencia

Proyecto académico - UIDE 2026