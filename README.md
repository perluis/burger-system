# 🍔 Sistema de Gestión - Restaurante de Hamburguesas

Sistema de gestión para restaurante desarrollado en Go con API REST y serialización JSON.

**Curso:** Programación Orientada a Objetos  
**Proyecto:** El impacto de las nuevas tecnologías en la sociedad: visualización del futuro.

---

## 📋 Descripción

API REST para gestión completa de restaurante de hamburguesas que permite:
- Administrar menú de hamburguesas (CRUD completo)
- Gestionar órdenes de clientes
- Cálculo automático de totales con IVA
- Control de estados de órdenes
- Serialización completa en JSON

---

## 🛠️ Tecnologías Utilizadas

- **Lenguaje:** Go 1.21+
- **Framework Web:** Gorilla Mux
- **Formato de Datos:** JSON
- **Arquitectura:** REST API
- **Almacenamiento:** En memoria (RAM)

---

## 📁 Estructura del Proyecto
```
burger-system/
├── main.go                 # Servidor HTTP principal
├── api/
│   └── handlers.go         # Controladores HTTP
├── storage/
│   └── memory.go           # Almacenamiento en memoria
├── menu/
│   └── hamburguesa.go      # Modelo y lógica de hamburguesas
├── orders/
│   └── orden.go            # Modelo y lógica de órdenes
├── utils/
│   ├── interfaces.go       # Interfaces del sistema
│   └── helpers.go          # Funciones auxiliares
└── README.md
```

---

## 🚀 Instalación y Ejecución

### Prerrequisitos
- Go 1.21 o superior instalado
- Git

### Pasos

1. **Clonar el repositorio:**
```bash
git clone https://github.com/perluis/burger-system.git
cd burger-system
```

2. **Instalar dependencias:**
```bash
go mod download
```

3. **Ejecutar el servidor:**
```bash
go run main.go
```

4. **Verificar:**
Abrir navegador en: http://localhost:8080

---

## 📡 API Endpoints

### Hamburguesas

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/hamburguesas` | Listar todas las hamburguesas |
| GET | `/api/hamburguesas/{id}` | Obtener hamburguesa por ID |
| POST | `/api/hamburguesas` | Crear nueva hamburguesa |
| PUT | `/api/hamburguesas/{id}` | Actualizar hamburguesa |
| DELETE | `/api/hamburguesas/{id}` | Eliminar hamburguesa |

### Órdenes

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/ordenes` | Crear nueva orden |
| GET | `/api/ordenes/{id}` | Obtener orden por ID |
| PUT | `/api/ordenes/{id}/estado` | Actualizar estado de orden |

---

## 📖 Ejemplos de Uso

### Listar todas las hamburguesas
```bash
curl http://localhost:8080/api/hamburguesas
```

### Crear una hamburguesa
```bash
curl -X POST http://localhost:8080/api/hamburguesas \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Hamburguesa Especial",
    "descripcion": "Con ingredientes premium",
    "precio": 11.99,
    "categoria": "Premium",
    "ingredientes": ["carne", "queso", "bacon", "aguacate"]
  }'
```

### Crear una orden
```bash
curl -X POST http://localhost:8080/api/ordenes \
  -H "Content-Type: application/json" \
  -d '{
    "tipoOrden": "mesa",
    "numeroMesa": 5,
    "items": [
      {
        "hamburguesaID": "BURG-001",
        "cantidad": 2,
        "notas": "sin cebolla"
      }
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

## 🏗️ Arquitectura

### Principios Aplicados
- **Encapsulación:** Campos privados con getters
- **Interfaces:** Polimorfismo para operaciones comunes
- **Manejo de Errores:** Validaciones en todas las operaciones
- **Separación de Responsabilidades:** Paquetes independientes

### Flujo de una Petición
```
Cliente → HTTP Request → Router (Mux) → Handler → Storage → Response (JSON)
```

---

## 🔧 Funcionalidades Implementadas

### Módulo de Hamburguesas
- ✅ CRUD completo
- ✅ Validación de datos (precio, categoría)
- ✅ Generación automática de IDs
- ✅ Control de disponibilidad

### Módulo de Órdenes
- ✅ Creación de órdenes multi-item
- ✅ Cálculo automático de totales
- ✅ Aplicación de IVA (12% Ecuador)
- ✅ Gestión de estados (PENDIENTE → EN_COCINA → LISTA → ENTREGADA)
- ✅ Notas personalizadas por item

---

## 🎯 Características Técnicas

- **Serialización JSON:** Todas las respuestas en formato JSON
- **Validaciones:** Entrada de datos validada
- **Concurrencia:** Mutex para acceso seguro al storage
- **Estados:** Máquina de estados para órdenes
- **Cálculos:** Subtotales e IVA automáticos

---

## 📝 Notas Importantes

- El sistema usa almacenamiento en memoria (los datos se pierden al reiniciar)
- El puerto por defecto es 8080
- Todas las respuestas son en formato JSON
- El IVA aplicado es del 12% (Ecuador)

---

## 🔮 Evolución Futura

- Persistencia en base de datos (PostgreSQL/MySQL)
- Autenticación y autorización
- Interfaz web con frontend
- Sistema de reportes avanzados
- Integración con sistemas de pago
- Notificaciones en tiempo real

---

## 👨‍💻 Autor

**Luis Agapito**  
Universidad Internacional del Ecuador (UIDE)  

---

## 📄 Licencia

Proyecto académico - UIDE 2026
