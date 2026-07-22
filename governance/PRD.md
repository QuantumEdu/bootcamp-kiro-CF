# PRD: POS SaaS para Suplementos, Vitaminas y Nutracéuticos

**Versión:** 1.0
**Estado:** MVP avanzado — backend completo, frontend en desarrollo
**Stack:** Next.js 16 + Supabase (PostgreSQL, RLS, Edge Functions Deno)

---

## Problem Statement

Las tiendas de suplementos, vitaminas y nutracéuticos en LATAM operan con sistemas POS genéricos que no entienden su negocio. El inventario con lotes, caducidades, múltiples presentaciones por producto y control FEFO (First Expired First Out) son requisitos fundamentales que ningún sistema genérico resuelve bien.

El resultado: inventario poco confiable, faltantes frecuentes, productos vencidos no detectados, créditos mal controlados y cero trazabilidad. El dueño de la tienda no puede responder rápido a preguntas como "¿cuánto vendí hoy?" o "¿qué debo comprar?".

Además, el mercado de suplementos crece aceleradamente en la región, y los dueños necesitan un sistema que no solo registre operaciones, sino que les ayude a tomar mejores decisiones de compra, inventario y crédito.

---

## Solution

Un SaaS multiempresa especializado en suplementos, construido 100% sobre Supabase (sin backend externo). El frontend es Next.js 16 con React 19, Tailwind CSS 4 y Supabase SSR.

La arquitectura separa claramente:

- **Lecturas**: Frontend → SDK Supabase → RLS (cada empresa ve solo sus datos)
- **Operaciones críticas**: Frontend → Edge Function (valida usuario/empresa/rol) → RPC SQL transaccional (atómico)

Esto garantiza que ninguna operación de dinero, inventario o cobranza pueda dejar el sistema en estado inconsistente.

El sistema está preparado para evolucionar hacia AI-first, con datos estructurados, APIs claras y un modelo de dominio rico que permite integrar predicción de demanda, recomendaciones inteligentes y automatización de compras.

---

## User Stories

### Core — Multiempresa y Seguridad

1. Como dueño de una tienda, quiero crear mi empresa en el sistema, para que mis datos estén aislados de otras empresas.
2. Como dueño, quiero crear sucursales para mi empresa, para gestionar inventarios y cajas separados por local.
3. Como dueño, quiero crear usuarios con roles (admin/cajero), para que cada persona acceda solo a lo que necesita.
4. Como cajero, quiero acceder solo a mi sucursal asignada, para no ver información de otras sucursales.
5. Como administrador, quiero que ningún usuario de otra empresa pueda ver mis datos, para garantizar la privacidad de mi negocio.

### Catálogo de Productos

6. Como administrador, quiero crear productos con nombre, marca, categoría y descripción, para organizar mi catálogo.
7. Como administrador, quiero crear variantes de un producto (SKU, código de barras, presentación, tamaño, precio), para vender el mismo producto en diferentes formatos.
8. Como administrador, quiero gestionar marcas y categorías, para clasificar mis productos correctamente.
9. Como administrador, quiero activar/desactivar productos sin borrarlos, para mantener el historial de ventas.

### Compras y Recepción

10. Como administrador, quiero registrar proveedores, para saber a quién le compro.
11. Como administrador, quiero crear pedidos de compra en estado borrador, para preparar compras antes de enviarlas.
12. Como administrador, quiero enviar un pedido a mi proveedor, para iniciar el proceso de compra.
13. Como administrador, quiero recibir mercancía parcial o totalmente, para actualizar el inventario solo cuando la mercancía llega físicamente.
14. Como administrador, quiero registrar lotes y fechas de caducidad al recibir mercancía, para controlar la frescura de mis productos.
15. Como administrador, quiero cancelar pedidos, para manejar compras que no se concretan.

### Inventario

16. Como administrador, quiero ver la existencia física y disponible de cada producto, para saber qué tengo realmente.
17. Como administrador, quiero saber qué productos están comprometidos (apartados/pre-vendidos), para no sobre-vender.
18. Como administrador, quiero ajustar inventario con motivo registrado, para corregir diferencias sin perder trazabilidad.
19. Como administrador, quiero registrar mermas, para dar de baja productos dañados o vencidos.
20. Como administrador, quiero ver alertas de productos próximos a caducar, para tomar acción antes de que se venzan.
21. Como administrador, quiero que el sistema consuma primero los lotes más próximos a vencer (FEFO), para reducir pérdidas por caducidad.

### Clientes y Demanda

22. Como administrador, quiero registrar clientes con nombre y teléfono, para gestionar créditos y preventas.
23. Como administrador, quiero registrar solicitudes de clientes, para saber qué productos quiere la gente que aún no tengo.
24. Como administrador, quiero registrar preventas, para asegurar stock antes de que llegue.
25. Como administrador, quiero ver sugerencias de compra basadas en solicitudes y ventas, para decidir qué comprar.

### POS y Ventas

26. Como cajero, quiero capturar productos por código de barras, para vender rápido.
27. Como cajero, quiero buscar productos por nombre o SKU, para cuando no hay código de barras.
28. Como cajero, quiero ver un catálogo visual de productos, para seleccionar rápidamente.
29. Como cajero, quiero aplicar descuentos (con autorización), para manejar promociones y cortesías.
30. Como cajero, quiero cobrar en efectivo, tarjeta, transferencia o mixto, para aceptar cualquier método de pago.
31. Como cajero, quiero registrar ventas a crédito (con cliente y autorización), para clientes que pagan después.
32. Como administrador, quiero que cada venta descuente inventario automáticamente por FEFO, para mantener el stock siempre actualizado.

### Caja

33. Como cajero, quiero abrir caja con un monto inicial, para comenzar mi turno.
34. Como cajero, quiero cerrar caja al finalizar, para cuadrar mis ventas con el efectivo.
35. Como administrador, quiero ver las diferencias de caja, para detectar errores o inconsistencias.

### Crédito y Abonos

36. Como administrador, quiero ver los saldos pendientes de cada cliente, para saber quién me debe.
37. Como administrador, quiero registrar abonos de clientes, para actualizar sus saldos.
38. Como administrador, quiero ver el estado de cuenta de un cliente, para saber su historial de crédito.

### Devoluciones

39. Como administrador, quiero procesar devoluciones totales o parciales, para manejar cambios y reclamos.
40. Como administrador, quiero elegir el destino del producto devuelto (inventario, merma, garantía, desecho), para decidir según el estado del producto.
41. Como administrador, quiero que las devoluciones nunca borren la venta original, para mantener el historial completo.

### Dashboard y Reportes

42. Como administrador, quiero ver las ventas de hoy, esta semana y este mes en el dashboard, para tener visibilidad inmediata.
43. Como administrador, quiero ver productos con stock bajo y agotados, para decidir compras.
44. Como administrador, quiero ver productos próximos a caducar, para evitar pérdidas.
45. Como administrador, quiero ver créditos pendientes, para gestionar cobranza.
46. Como administrador, quiero exportar reportes a CSV y Excel, para tener control sobre mis datos.

### Auditoría

47. Como administrador, quiero que toda operación crítica quede registrada con usuario y fecha, para saber quién hizo qué.
48. Como administrador, quiero consultar el historial de cambios, para investigar inconsistencias.

### AI-First (Futuro)

49. Como administrador, quiero que el sistema me sugiera cantidades de compra basadas en histórico de ventas, para optimizar mi inventario.
50. Como administrador, quiero recomendaciones de productos complementarios en el POS, para aumentar el ticket promedio.
51. Como administrador, quiero preguntar en lenguaje natural "¿cuánto vendí ayer?" y obtener la respuesta, para acceder a la información más rápido.
52. Como administrador, quiero alertas automáticas cuando haya anomalías en caja o inventario, para detectar problemas temprano.

---

## Implementation Decisions

### Arquitectura General

- **Sin backend externo**: 100% Supabase en V1. No hay servidores propios, solo Supabase Auth + PostgreSQL + RLS + Edge Functions.
- **Frontend**: Next.js 16 con App Router y Supabase SSR. Migrado desde Vue 3 (el init original de SDD referencia Vue, pero el código actual es Next.js).
- **Edge Functions**: En Deno, validan usuario/empresa/sucursal/rol antes de ejecutar lógica.
- **RPC SQL**: Toda operación crítica es una transacción PostgreSQL atómica.

### Modelo de Datos

- **Multiempresa**: Toda tabla operativa tiene `company_id`. RLS policies filtran por el usuario autenticado.
- **Soft delete**: `is_active`, `deleted_at`, `deleted_by` en todas las tablas críticas.
- **UUIDs**: Para todas las entidades principales.
- **Movimientos de inventario**: No se edita stock directamente — todo cambio genera un registro en `inventory_movements`.
- **Lotes**: `inventory_batches` con cantidad, fecha de recepción, fecha de caducidad y costo.

### Flujo de Operaciones Críticas

```
Frontend (Next.js)
    → Edge Function (Deno) — valida auth, company, branch, role, input
        → RPC SQL (PostgreSQL) — transacción atómica
            → INSERT/UPDATE tablas
            → INSERT movimiento inventario
            → INSERT auditoría
            → COMMIT / ROLLBACK
        ← Resultado consistente
    ← Response
```

### Migración de Stack

- El proyecto inició con **Vue 3** como frontend y migró a **Next.js 16 + React 19**.
- La documentación SDD original (sdd-init) aún referencia Vue 3 — necesita actualización.
- Esta migración no afectó el backend (100% Supabase, sin cambios).

### Preparación AI-First

- **Datos estructurados**: El modelo de dominio (14 migraciones, ~40 tablas) tiene datos históricos de ventas, inventario, clientes, créditos — ideales para entrenar modelos.
- **Edge Functions como API**: Cada operación crítica tiene un endpoint claro con validación y response consistente.
- **Puntos de integración detectados**: Predicción de demanda (historial de ventas), recomendaciones (combinaciones de productos), detección de anomalías (caja e inventario), asistente conversacional (datos agregados).
- **Carencias actuales**: No hay vector embeddings, no hay histórico de navegación/usuario, no hay pipelines de datos para ML.

---

## Testing Decisions

### Qué hace un buen test

- Prueba el comportamiento externo, no la implementación interna.
- Una operación de venta debe probar: auth → validación → descuento de inventario → creación de pago → auditoría → resultado.
- No probar queries individuales — probar transacciones completas.

### Tests existentes

- **Edge Functions**: Tests con Deno (`deno test supabase/functions/_test/`).
- **Base de datos**: `supabase test db` para migraciones.

### Prioridades de testing

1. **Transacciones críticas**: create_sale_transaction, receive_purchase_transaction, register_customer_payment_transaction, close_cash_session_transaction.
2. **RLS**: Pruebas de aislamiento por empresa y rol.
3. **FEFO**: Que el algoritmo de consumo de lotes respete orden de caducidad.
4. **Caja**: Que apertura → venta → cierre calcule diferencias correctamente.

### Tests faltantes

- Tests E2E del frontend (Next.js) — no implementados.
- Tests de integración frontend + Edge Functions — no implementados.
- Tests de regresión para AI features — futuros.

---

## Out of Scope (MVP)

Funcionalidades explicitamente excluidas del MVP actual:

| Funcionalidad | Prevista para |
|---|---|
| CFDI / Facturación electrónica | V2 |
| Pasarelas de pago en línea | V2 |
| Catálogo global de productos | V1.5 |
| Transferencias entre sucursales | V1.5 |
| Tickets PDF / térmicos | V1.5 |
| Notificaciones WhatsApp / correo | V2 |
| Promociones avanzadas | V2 |
| App móvil nativa | V3 |
| Modo offline | V3 |
| Predicción de demanda (AI) | V3 |
| Recomendaciones inteligentes (AI) | V3 |
| Automatización de compras (AI) | V3 |
| Asistente conversacional (AI) | V3+ |
| Integraciones con marketplaces | V3 |

---

## Further Notes

### Estado Actual del Desarrollo

El MVP tiene el backend **completo** (14 migrations, 36 Edge Functions) y el frontend **en desarrollo** (Next.js con login, dashboard, POS y productos funcionales). Es un proyecto real y operativo, no un prototipo.

### Stack Relevante para AI-First

- **PostgreSQL**: Datos relacionales perfectos para features de embeddings (pgvector).
- **Edge Functions (Deno)**: Pueden ejecutar inferencia ligera o llamar a APIs externas (OpenAI, Anthropic, etc.).
- **Next.js 16**: Puede servir UI generativa (chat, recomendaciones, dashboards inteligentes).

### Deuda Técnica Identificada

1. **Documentación desactualizada**: sdd-init.md referencia Vue 3, el código real es Next.js 16.
2. **Tests E2E faltantes**: Sin tests E2E, agregar features AI aumenta riesgo de regresión.
3. **Sin observabilidad**: No hay logging estructurado, métricas ni tracing — difícil debuggear en prod.
4. **API sin documentación**: Los contratos de Edge Functions no están documentados formalmente.

### Arquitectura de Base de Datos

14 migraciones que cubren:

| Migración | Dominio | Tablas |
|---|---|---|
| 00001 | Companies, Branches, Profiles | companies, branches, profiles |
| 00002 | RLS Helpers | Funciones helper para RLS |
| 00003 | RLS Policies | Políticas de seguridad |
| 00004 | Catalog Domain | brands, categories, units, products, product_variants |
| 00005 | Inventory Domain | inventory_batches, inventory_movements, inventory_adjustments, inventory_reservations |
| 00006 | Purchasing Domain | suppliers, purchase_orders, purchase_order_items, purchase_receipts, purchase_receipt_items |
| 00007 | Customers & Demand | customers, customer_requests, preorders, preorder_items |
| 00008 | Cash Session Domain | cash_sessions, cash_movements |
| 00009 | POS Sales Domain | sales, sale_items, sale_item_batches, payments, discount_authorizations |
| 00010 | Credit Payments | customer_balances, customer_payments |
| 00011 | Returns Domain | returns, return_items |
| 00012 | Dashboard & Reports | Vistas/materializadas para reportes |
| 00013 | Fixes | Corrección FK sales → customers |
| 00014 | Credit Limits | Validación de límite de crédito |

### Roadmap Sugerido para AI-First

1. **Estabilizar**: Tests E2E, documentación API, observabilidad básica.
2. **Datos**: Agregar pgvector, embeddings de productos, histórico de sesiones.
3. **AI Features**: Predicción de demanda → recomendaciones en POS → detección de anomalías → asistente conversacional.
4. **Automatización**: Generación automática de órdenes de compra basada en predicciones.
