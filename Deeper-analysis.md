# Comparativo Detallado — 3 Propuestas Finalistas

---

## 🔵 Propuesta A: App de Servicios Hiperlocal
*"La sección amarilla moderna para tu ciudad"*

**Qué resuelve:** El electricista, plomero o mesero de tu ciudad no tiene presencia digital accesible. El cliente que los necesita no sabe
dónde buscarlos fuera de referencias de boca en boca.

**MVP en 5 días:**
- Registro de prestadores (oficio, zona, foto, contacto, precio referencial)
- Búsqueda por categoría y zona dentro de una ciudad
- Perfil público limpio con botón de contacto directo (WhatsApp/tel)
- Panel admin para aprobar registros
- 10 prestadores reales como datos de prueba

**Arquitectura AWS + Kiro:**

Frontend (React/Next.js) → API Gateway → Lambda
                                           ↓
                                       DynamoDB (prestadores, categorías, zonas)
                                           ↓
                                        S3 (fotos de perfil)
Kiro → genera CRUD, búsqueda, UI minimalista

**Roadmap post-bootcamp (capa IA):**
- Recomendaciones por historial y zona
- Chatbot: "necesito un plomero disponible hoy en zona norte"
- Rating + análisis de sentimiento de reseñas

| Aspecto | Evaluación |
|---|---|
| Impacto tecnológico (30%) | ✅ Alto — necesidad real, audiencia no tech-savvy desatendida |
| Innovación (30%) | ⚠️ Media — directorio existe, diferenciador es hiperlocal + precio justo + UI simple |
| Software funcional (30%) | ✅ Alto — MVP totalmente demostrable en 5 días con datos reales |
| Uso AWS/Kiro (10%) | ✅ S3 + DynamoDB + Lambda + API Gateway + Kiro para generar código |
| Viabilidad 5 días | ✅✅ La más viable de las 3 |
| Demo day | ✅ Buscar "electricista zona centro" → resultados reales → perfil → WhatsApp |
| Riesgo | ⚠️ Bajo impacto técnico percibido vs proyectos con IA visible |
| Monetización real | ✅ $50 MXN/mes por prestador — modelo claro y sostenible |

---

## 🟡 Propuesta B: POS AI-First
*"El punto de venta que entiende lenguaje natural"*

**Qué resuelve:** Los POS tradicionales requieren capacitación y son rígidos. Un dueño de negocio debería poder preguntarle a su sistema
"¿qué vendí esta semana?" o "¿qué producto tiene más rotación?" en lenguaje natural.

**MVP en 5 días:**
- CRUD de productos e inventario
- Registro de ventas
- Interfaz conversacional: preguntas en lenguaje natural → SQL generado → respuesta
- Dashboard básico de métricas
- Al menos 5 consultas demostrables predefinidas

**Arquitectura AWS + Kiro:**

Frontend (React) → API Gateway → Lambda
                                    ↓
                              RDS/DynamoDB (productos, ventas, inventario)
                                    ↓
                              Bedrock (Claude/Titan) → NL→SQL → respuesta
Kiro → genera CRUD, integración Bedrock, UI del chat

**Riesgo técnico principal:**
El LLM debe generar SQL correcto sobre datos reales. Si falla la consulta o da datos incorrectos en el demo, el impacto se pierde.
Requiere guardrails y consultas de respaldo predefinidas.

| Aspecto | Evaluación |
|---|---|
| Impacto tecnológico (30%) | ✅ Alto — mercado POS enorme, AI-first es diferenciador claro |
| Innovación (30%) | ✅✅ Alto — conversacional desde el diseño, no como add-on |
| Software funcional (30%) | ✅ Alto — CRUD + chat funcional es demostrable |
| Uso AWS/Kiro (10%) | ✅✅ Bedrock es el core, no un añadido |
| Viabilidad 5 días | ✅ Alta — scope acotable si se limita a consultas y no a facturación completa |
| Demo day | ✅✅ "¿Cuánto vendí ayer?" hablado en vivo → respuesta inmediata. Muy impactante |
| Riesgo | ⚠️ NLU puede fallar en demo, requiere pruebas exhaustivas previas |
| Monetización real | ⚠️ Media — competencia fuerte (Bind, Aspel, Square) pero nicho PYME sin tech |

---

## 🔴 Propuesta C: Agente Telefónico Virtual
*"El empleado de call center que nunca duerme"*

**Qué resuelve:** Las PYMES necesitan atender llamadas pero no pueden costear una persona dedicada. Los agentes de texto (chatbots) no
funcionan para clientes que prefieren llamar. El agente atiende, transcribe, analiza sentimiento y aprende.

**MVP en 5 días:**
- Número telefónico que recibe llamadas (Twilio o Amazon Connect)
- STT en español (Amazon Transcribe)
- LLM responde como agente de negocio (Bedrock)
- TTS devuelve la respuesta en voz (Amazon Polly)
- Dashboard: transcripciones + análisis de sentimiento (Amazon Comprehend)
- Fine-tuning y RAG quedan como roadmap post-bootcamp

**Arquitectura AWS + Kiro:**

Llamada entrante → Amazon Connect / Twilio
                        ↓
               Amazon Transcribe (STT, español)
                        ↓
               Bedrock (Claude) → respuesta
                        ↓
               Amazon Polly (TTS, voz natural)
                        ↓
               DynamoDB (log de llamadas)
                        ↓
               Comprehend (análisis de sentimiento)
                        ↓
               Dashboard web (React) ← Kiro genera UI

**Riesgo técnico principal:**
La **latencia** es el enemigo. STT → LLM → TTS puede tardar 3-8 segundos por turno, lo que hace la conversación antinatural. Mitigaciones:
respuestas pregeneradas para intents comunes, streaming de audio.

| Aspecto | Evaluación |
|---|---|
| Impacto tecnológico (30%) | ✅✅ Muy alto — nicho desatendido, voz en español para PYMES |
| Innovación (30%) | ✅✅ Muy alto — la mayoría va a texto, voz es diferenciador real |
| Software funcional (30%) | ⚠️ Media — funcional pero latencia puede arruinar el demo |
| Uso AWS/Kiro (10%) | ✅✅✅ El que más servicios AWS usa (Connect/Transcribe/Bedrock/Polly/Comprehend) |
| Viabilidad 5 días | ⚠️ Media — la más compleja de las 3, muchos servicios que integrar |
| Demo day | ✅✅✅ El más impactante si funciona: llamar en vivo durante la presentación |
| Riesgo | 🔴 Alto — latencia, STT en español impreciso, múltiples APIs que pueden fallar |
| Monetización real | ✅✅ Muy clara — SaaS por minuto o suscripción mensual por negocio |

---

## Resumen comparativo final

| Criterio | A. App Servicios | B. POS AI-First | C. Agente Telefónico |
|---|---|---|---|
| Impacto tecnológico (30%) | ✅ Alto | ✅ Alto | ✅✅ Muy alto |
| Innovación (30%) | ⚠️ Media | ✅✅ Alta | ✅✅ Muy alta |
| Software funcional (30%) | ✅✅ Más seguro | ✅ Alto | ⚠️ Riesgoso |
| AWS/Kiro (10%) | ✅ Básico-Medio | ✅✅ Bedrock core | ✅✅✅ Máximo uso |
| Viabilidad 5 días | ✅✅ La más segura | ✅ Alta | ⚠️ La más riesgosa |
| Demo day | ✅ Sólido | ✅✅ Impactante | ✅✅✅ WOW factor |
| Riesgo de fallo | 🟢 Bajo | 🟡 Medio | 🔴 Alto |
| Potencial post-bootcamp | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

## La decisión clave es tolerancia al riesgo

- Si priorizas **entregar algo impecable y funcional** → **App de Servicios**
- Si quieres **balance entre impacto técnico y viabilidad** → **POS AI-First**
- Si quieres **máximo impacto y puntaje, con riesgo alto** → **Agente Telefónico**
