# Plano de Integração OpenFGA - Lockari Backend

## 1. Estrutura de Arquivos a Serem Criados

### 1.1. Configuração OpenFGA
```
pkg/authorization/
├── openfga/
│   ├── client.go              # Cliente OpenFGA
│   ├── config.go              # Configuração OpenFGA
│   ├── models.go              # Modelos DSL
│   ├── tuples.go              # Gerenciamento de tuplas
│   └── middleware.go          # Middleware de autorização
```

### 1.2. Entidades e Repositórios
```
internal/core/entity/authorization/
├── types.go                   # Tipos base de autorização
├── vault_permissions.go       # Permissões de vault
└── permission_service.go      # Serviço de permissões

internal/core/repository/authorization/
├── openfga_repository.go      # Repository para OpenFGA
└── permission_repository.go   # Repository de permissões
```

### 1.3. Handlers e Middleware
```
internal/handler/middleware/
├── authorization.go           # Middleware de autorização OpenFGA
└── tenant_isolation.go        # Middleware de isolamento multi-tenant
```

## 2. Modelo de Autorização OpenFGA (DSL)

### 2.1. Modelo Inicial para Vaults e Secrets

```json
{
  "schema_version": "1.1",
  "type_definitions": [
    {
      "type": "tenant",
      "relations": {
        "member": {
          "this": {}
        },
        "admin": {
          "this": {}
        }
      },
      "metadata": {
        "relations": {
          "member": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "admin": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          }
        }
      }
    },
    {
      "type": "vault",
      "relations": {
        "owner": {
          "this": {}
        },
        "admin": {
          "this": {}
        },
        "writer": {
          "this": {}
        },
        "viewer": {
          "this": {}
        },
        "tenant": {
          "this": {}
        }
      },
      "metadata": {
        "relations": {
          "owner": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "admin": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "writer": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "viewer": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "tenant": {
            "directly_related_user_types": [
              {
                "type": "tenant"
              }
            ]
          }
        }
      }
    },
    {
      "type": "secret",
      "relations": {
        "vault": {
          "this": {}
        },
        "read": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "viewer"
                  }
                }
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "writer"
                  }
                }
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "admin"
                  }
                }
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "owner"
                  }
                }
              }
            ]
          }
        },
        "write": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "writer"
                  }
                }
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "admin"
                  }
                }
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "owner"
                  }
                }
              }
            ]
          }
        },
        "delete": {
          "union": {
            "child": [
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "admin"
                  }
                }
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "relation": "vault"
                  },
                  "computedUserset": {
                    "relation": "owner"
                  }
                }
              }
            ]
          }
        }
      },
      "metadata": {
        "relations": {
          "vault": {
            "directly_related_user_types": [
              {
                "type": "vault"
              }
            ]
          }
        }
      }
    }
  ]
}
```

## 3. Integração com Firebase/Firestore Multi-Tenant

### 3.1. Estrutura de Collections Multi-Tenant

```
/tenants/{tenantId}/
├── users/
│   └── {userId}/
├── vaults/
│   └── {vaultId}/
│       ├── metadata (nome, descrição, tags)
│       ├── secrets/
│       │   └── {secretId}/
│       ├── certificates/
│       │   └── {certId}/
│       ├── ssh_keys/
│       │   └── {keyId}/
│       └── access_permissions/
│           └── {permissionId}/
├── audit_logs/
│   └── {logId}/
└── subscription/
    └── metadata
```

### 3.2. Custom Claims JWT Structure

```json
{
  "sub": "user123",
  "email": "user@example.com",
  "tenantId": "tenant_abc123",
  "roles": ["vault_owner", "tenant_member"],
  "permissions": {
    "vaults": ["create", "read"],
    "admin": ["read_users"]
  }
}
```

## 4. Implementação Prática

### 4.1. Configuração OpenFGA no config.yaml

```yaml
openfga:
  api_url: "http://localhost:8080"
  store_id: "01HXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  authorization_model_id: "01HYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  shared_secret: "${OPENFGA_SHARED_SECRET}"
```

### 4.2. Fluxo de Verificação de Permissões

```
1. Request chega ao endpoint → Middleware de Autenticação Firebase
2. Extrai tenantId do JWT → Middleware de Autorização OpenFGA  
3. Verifica permissão específica → Executa ação se autorizado
4. Registra evento de auditoria → Retorna resposta
```

### 4.3. Exemplos de Tuplas de Relacionamento

```
# Usuário é membro de um tenant
user:alice#member@tenant:company-abc

# Usuário é owner de um vault
user:alice#owner@vault:marketing-secrets

# Vault pertence a um tenant
vault:marketing-secrets#tenant@tenant:company-abc

# Secret pertence a um vault
secret:api-key-prod#vault@vault:marketing-secrets
```

## 5. Tarefas Detalhadas por Etapa

### Etapa 1: Setup Básico (Semana 1)
- [ ] Adicionar dependências OpenFGA ao go.mod
- [ ] Criar estrutura base de configuração
- [ ] Implementar cliente OpenFGA básico
- [ ] Configurar modelo DSL inicial
- [ ] Testes unitários básicos

### Etapa 2: Integração Firebase (Semana 1-2)
- [ ] Modificar processo de signup para criar tenant
- [ ] Implementar adição de tenantId ao JWT
- [ ] Criar estrutura multi-tenant no Firestore
- [ ] Implementar regras de segurança Firestore
- [ ] Testes de isolamento de tenant

### Etapa 3: Middleware e Autorização (Semana 2)
- [ ] Criar middleware de autorização OpenFGA
- [ ] Integrar verificações nos endpoints existentes
- [ ] Implementar gerenciamento de tuplas
- [ ] Criar endpoints de teste para verificação
- [ ] Testes de autorização

### Etapa 4: Endpoints Vault/Secret (Semana 2-3)
- [ ] Criar endpoints CRUD para vaults
- [ ] Criar endpoints CRUD para secrets
- [ ] Integrar verificações OpenFGA em todos endpoints
- [ ] Implementar auditoria de ações
- [ ] Testes de integração completos

### Etapa 5: Validação e Otimização (Semana 3)
- [ ] Testes de performance
- [ ] Validação de regras de negócio
- [ ] Documentação Swagger
- [ ] Deploy em ambiente de staging
- [ ] Testes end-to-end

## 6. Considerações de Segurança

### 6.1. Isolamento Multi-Tenant
- Todas as queries devem incluir tenantId
- Middleware deve validar tenantId do JWT
- OpenFGA deve verificar tenant ownership

### 6.2. Princípio do Menor Privilégio
- Usuários só têm acesso ao que precisam
- Verificações granulares por operação
- Logs detalhados de todas as ações

### 6.3. Backup e Recuperação
- Backup das tuplas OpenFGA
- Sincronização entre Firestore e OpenFGA
- Procedimentos de rollback

## 7. Métricas e Monitoramento

### 7.1. Métricas OpenFGA
- Latência das verificações
- Taxa de sucesso/falha
- Volume de tuplas por tenant

### 7.2. Métricas Multi-Tenancy
- Isolamento de dados
- Performance por tenant
- Uso de recursos

## 8. Próximos Passos

1. **Discussão Arquitetural**: Revisar modelo DSL proposto
2. **Prototipo Inicial**: Implementar cliente OpenFGA básico
3. **Testes de Conceito**: Validar integração Firebase + OpenFGA
4. **Refinamento**: Ajustar modelo baseado nos testes
5. **Implementação Gradual**: Começar com funcionalidades básicas

## 9. Dependências Externas

### 9.1. OpenFGA Server
- Necessário deploy do OpenFGA Server (Cloud Run/GKE)
- Banco PostgreSQL para persistência
- Configuração de rede e segurança

### 9.2. Ferramentas de Desenvolvimento
- OpenFGA CLI para testes
- Playground OpenFGA para debug
- Ferramentas de migração de modelo

## 10. Riscos e Mitigações

### 10.1. Riscos Técnicos
- **Latência**: Cache de verificações frequentes
- **Complexidade**: Documentação detalhada e testes
- **Sincronização**: Transações atômicas entre sistemas

### 10.2. Riscos de Negócio
- **Learning Curve**: Treinamento da equipe
- **Performance**: Testes de carga extensivos
- **Backup**: Estratégia de DR robusta
