# PlugaWAPI — Roadmap

## Fase 1 — Base (mensagens + voz 1:1) — estimativa: 2 a 4 semanas
- Rebrand do WaCalls (MIT) para PlugaWAPI ✓ (build OK)
- Adicionar endpoints de mensagens (texto, mídia, contatos, grupos) baseados no GOWA
- Webhooks de mensagens (call.offer, message, etc.)
- Integração com Pluga Base (configuração + chamadas)

## Fase 2 — Vídeo, grupo, SIP/PBX, gravação, hold/transferência, avançado — estimativa: 2 a 3 meses
- Port de vídeo do AstraCalls (AGPL) — requer aceitar AGPL ou reimplementar
- Chamadas em grupo (áudio/vídeo)
- SIP/PBX gateway
- Gravação, hold, transferência
- Formulários, pix, carousel, etc.

## Licença
- Base: MIT (WaCalls)
- Vídeo/grupo: AGPL (AstraCalls) se portado

## Status
- Fase 1 (base): rebrand WaCalls OK, build OK, endpoints de envio de texto/mídia, contatos/grupos, reações, presença, check de usuário, webhook de mensagens, histórico local (SQLite) adicionados. Sem dependência do GOWA.
- Fase 2 (parcial): Hold/Resume/Mute, Gravação (WAV) e Transcrição+Análise IA implementados. Grupo, SIP/PBX, vídeo pendentes. Transferência adiada para Fase 3 (multiusuário).
- Próximo: integrar Pluga Base com PlugaWAPI (config + chamadas).

## Fase 3 — Multiusuário e transferência — estimativa: 2 a 3 meses
- Transferência cega entre atendentes
- Multiusuário, permissões, SaaS
