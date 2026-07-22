  Aquí va la comparativa completa de 6 propuestas:

  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

  Comparativa de 6 Propuestas para el Bootcamp Código Facilito + Kiro

  ┌────────────┬─────────────────────────┬──────────────┬───────────────┬───────────────────┬───────────────────┬─────────────────────┐
  │ Criterio   │ 1. Observabilidad LLMs  │ 2. App de    │ 3. POS        │ 4. Agente         │ 5. Code Review    │ 6. AI Expense       │
  │            │ (Go)                    │ Servicios/Of │ AI-First      │ Telefónico        │ Agent (MCP)       │ Splitter (Web App)  │
  │            │                         │ icios        │               │ Virtual           │                   │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Categoría  │ Productividad devs      │ Aplicaciones │ Aplicaciones  │ Agentes           │ Agentes           │ Aplicaciones web    │
  │ bootcamp   │                         │ web          │ web           │ especializados    │ especializados    │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Viabilidad │ ⚠️ Media-Alta           │ ❌ Baja      │ ✅ Alta       │ ⚠️ Media          │ ✅ Alta (GitHub   │ ✅ Alta (CRUD +     │
  │ 5 días     │                         │              │               │ (requiere         │ API + LLM)        │ OCR/LLM)            │
  │            │                         │              │               │ telephony APIs,   │                   │                     │
  │            │                         │              │               │ STT/TTS)          │                   │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Impacto    │ Alto — nicho dev/MLOps  │ Medio —      │ Alto —        │ Alto — nicho      │ Alto — dolor real │ Medio-Alto —        │
  │ tecnológic │                         │ competencia  │ diferenciador │ desatendido (voz  │ de equipos dev    │ problema cotidiano  │
  │ o (30%)    │                         │ fuerte       │ vs POS        │ vs texto)         │                   │ universal           │
  │            │                         │              │ tradicional   │                   │                   │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Innovación │ Media —                 │ Baja —       │ Alta — POS    │ Alta — pocos      │ Alta — MCP +      │ Media — Splitwise   │
  │ (30%)      │ Langfuse/LangSmith      │ Workana,     │ conversaciona │ agentes de voz en │ agente            │ existe, pero        │
  │            │ existen, diferenciador  │ TaskRabbit   │ l es nuevo    │ español para      │ especializado es  │ AI-first con OCR es │
  │            │ es Go+self-hosted       │ dominan      │               │ PYMES             │ bleeding edge     │ nuevo               │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Software   │ Dashboard + ingesta de  │ Landing +    │ CRUD + chat   │ Llamada demo +    │ PR review         │ Foto de ticket →    │
  │ funcional  │ traces                  │ matching     │ funcional     │ transcripción +   │ automático +      │ split automático    │
  │ (30%)      │                         │ básico       │               │ sentimiento       │ comentarios       │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Uso        │ Lambda + DynamoDB/S3    │ Amplify +    │ Bedrock +     │ Amazon Connect +  │ Bedrock + Lambda  │ Textract/Rekognitio │
  │ AWS/Kiro   │ para traces             │ API Gateway  │ DynamoDB      │ Transcribe +      │ +                 │ n + Bedrock + S3    │
  │ (10%)      │                         │ + DB         │               │ Bedrock +         │ CodeCommit/GitHub │                     │
  │            │                         │              │               │ Comprehend        │ webhook           │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Demostrabl │ ✅ Dashboard impresiona │ ⚠️ Difícil   │ ✅ Muy        │ ✅✅ MUY          │ ✅ Abrir PR → ver │ ✅ Foto de ticket → │
  │ e (demo    │                         │ sin usuarios │ visual:       │ impactante:       │ review automático │ resultado           │
  │ day)       │                         │              │ hablarle al   │ llamar en vivo    │                   │ instantáneo         │
  │            │                         │              │ POS           │                   │                   │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Complejida │ Manejo de secretos se   │ Pagos,       │ NLU preciso   │ Latencia voz,     │ Contexto de       │ OCR impreciso,      │
  │ d oculta   │ infla                   │ verificación │ sobre datos   │ STT/TTS en        │ código grande,    │ divisiones ambiguas │
  │            │                         │ , UX mobile  │ reales        │ español, manejo   │ falsos positivos  │                     │
  │            │                         │              │               │ de turnos         │                   │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Diferencia │ Open source +           │ Precio justo │ Conversaciona │ Voz en español +  │ MCP-native +      │ AI-first: no tipeas │
  │ dor clave  │ self-hosted + Go        │ + sin ads    │ l desde       │ análisis de       │ integración Kiro  │ nada, solo foto     │
  │            │                         │              │ diseño        │ sentimiento + RAG │ nativa            │                     │
  ├────────────┼─────────────────────────┼──────────────┼───────────────┼───────────────────┼───────────────────┼─────────────────────┤
  │ Riesgo     │ Scope creep en features │ Marketplace  │ LLM           │ Telephony         │ Calidad del       │ OCR falla con       │
  │ principal  │                         │ de 2 lados   │ malinterpreta │ compleja +        │ review depende    │ tickets mal         │
  │            │                         │              │ queries       │ latencia          │ del contexto      │ impresos            │
  └────────────┴─────────────────────────┴──────────────┴───────────────┴───────────────────┴───────────────────┴─────────────────────┘

  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

  Detalle de las 2 propuestas nuevas:

  Propuesta 5: Code Review Agent (MCP-based)

  Categoría: Agentes especializados

  Un agente que se conecta a GitHub vía webhooks, y al detectar un nuevo PR automáticamente:

  - Analiza el diff usando un LLM (Bedrock)
  - Identifica bugs potenciales, malas prácticas, problemas de seguridad
  - Comenta directamente en el PR con sugerencias accionables
  - Usa MCPs (Model Context Protocol) para acceder al contexto del repo

  ¿Por qué es buena opción? Es exactamente la temática "Agentes especializados" del bootcamp. Usa Kiro de forma nativa (Kiro usa MCPs). El
  dolor es real (todo equipo quiere reviews más rápidos). Y es muy demostrable: abres un PR, en 30 segundos aparecen comentarios
  inteligentes.

  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

  Propuesta 6: AI Expense Splitter (Web App)

  Categoría: Aplicaciones web

  Una app web donde tomas una foto de un ticket de restaurante/compra y automáticamente:

  - Extrae los ítems y precios con OCR (AWS Textract)
  - Asigna quién pidió qué mediante una interfaz simple
  - Calcula propinas, impuestos y división
  - Genera links de pago o resúmenes compartibles

  ¿Por qué es buena opción? Resuelve un problema cotidiano (criterio del bootcamp para apps web). Usa APIs del navegador (cámara) +
  servicios AWS (Textract, S3, Bedrock para interpretar tickets ambiguos). Es instantáneamente demostrable.

  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

  Mi ranking actualizado (considerando todos los criterios de evaluación):

  ┌─────┬───────────────────────────┬─────────────────────────────────────────────────────────┐
  │ #   │ Propuesta                 │ Score estimado                                          │
  ├─────┼───────────────────────────┼─────────────────────────────────────────────────────────┤
  │ 1   │ POS AI-First              │ ⭐⭐⭐⭐⭐ — Viabilidad + impacto + demo killer         │
  ├─────┼───────────────────────────┼─────────────────────────────────────────────────────────┤
  │ 2   │ Code Review Agent (MCP)   │ ⭐⭐⭐⭐⭐ — Innovación máxima + uso nativo de Kiro/AWS │
  ├─────┼───────────────────────────┼─────────────────────────────────────────────────────────┤
  │ 3   │ Agente Telefónico Virtual │ ⭐⭐⭐⭐ — Demo impactante pero complejidad alta        │
  ├─────┼───────────────────────────┼─────────────────────────────────────────────────────────┤
  │ 4   │ Observabilidad LLMs       │ ⭐⭐⭐⭐ — Sólido pero competencia existente            │
  ├─────┼───────────────────────────┼─────────────────────────────────────────────────────────┤
  │ 5   │ AI Expense Splitter       │ ⭐⭐⭐½ — Fácil de hacer pero innovación moderada       │
  ├─────┼───────────────────────────┼─────────────────────────────────────────────────────────┤
  │ 6   │ App de Servicios/Oficios  │ ⭐⭐½ — Inviable en 5 días                              │
  └─────┴───────────────────────────┴─────────────────────────────────────────────────────────┘
