---
title: "Modelo de datos del MVP"
labels: wayfinder:research
blocking: ["002-setup-proyecto-go", "003-layout-navegacion-ui", "005-dashboard-metricas"]
---

## Question

Definir el schema completo de SQLite para el MVP: tablas, columnas, tipos, relaciones, constraints e índices.

Debe cubrir:

- **Productos:** nombre, precio, categoría, stock actual, SKU/código
- **Ventas:** fecha, total, cajero que registró, items de la venta (producto, cantidad, precio_unitario)
- **Inventario:** movimientos (entrada/salida/ajuste), producto, cantidad, fecha, motivo, usuario
- **Clientes:** nombre, teléfono, dirección, fecha de registro
- **Usuarios (cajeros):** nombre, PIN (hasheado), rol

Incluir:
- DDL completo con tipos SQLite
- Comentarios en castellano
- Índices sugeridos para las consultas típicas del POS
- Archivo `.sql` listo para sqlc
