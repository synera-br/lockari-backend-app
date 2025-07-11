# Exemplo de Modelo de Autorização para Lockari

Este arquivo contém exemplos de configuração do OpenFGA para o sistema Lockari.

## Modelo Base

```json
{
  "schema_version": "1.1",
  "type_definitions": [
    {
      "type": "user",
      "relations": {},
      "metadata": {
        "relations": {}
      }
    },
    {
      "type": "tenant",
      "relations": {
        "owner": {
          "this": {}
        },
        "admin": {
          "this": {}
        },
        "member": {
          "this": {}
        },
        "viewer": {
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
          "member": {
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
          }
        }
      }
    },
    {
      "type": "group",
      "relations": {
        "member": {
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
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "owner"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "admin"
                  }
                }
              }
            ]
          }
        },
        "editor": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "admin"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "editor"
                  }
                }
              }
            ]
          }
        },
        "viewer": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "editor"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "viewer"
                  }
                }
              }
            ]
          }
        },
        "parent": {
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
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "editor": {
            "directly_related_user_types": [
              {
                "type": "user"
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "viewer": {
            "directly_related_user_types": [
              {
                "type": "user"
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "parent": {
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
      "type": "token",
      "relations": {
        "owner": {
          "this": {}
        },
        "vault": {
          "this": {}
        },
        "can_read_secrets": {
          "this": {}
        },
        "can_write_secrets": {
          "this": {}
        },
        "can_manage_vault": {
          "this": {}
        },
        "revoked": {
          "this": {}
        },
        "expires": {
          "this": {}
        },
        "production": {
          "this": {}
        },
        "read_only": {
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
          "vault": {
            "directly_related_user_types": [
              {
                "type": "vault"
              }
            ]
          },
          "can_read_secrets": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "can_write_secrets": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "can_manage_vault": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "revoked": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "expires": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "production": {
            "directly_related_user_types": [
              {
                "type": "user"
              }
            ]
          },
          "read_only": {
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
      "type": "secret",
      "relations": {
        "owner": {
          "this": {}
        },
        "editor": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "owner"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "editor"
                  }
                }
              }
            ]
          }
        },
        "viewer": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "editor"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "viewer"
                  }
                }
              }
            ]
          }
        },
        "parent": {
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
          "editor": {
            "directly_related_user_types": [
              {
                "type": "user"
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "viewer": {
            "directly_related_user_types": [
              {
                "type": "user"
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "parent": {
            "directly_related_user_types": [
              {
                "type": "vault"
              }
            ]
          }
        }
      }
    },
    {
      "type": "certificate",
      "relations": {
        "owner": {
          "this": {}
        },
        "editor": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "owner"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "editor"
                  }
                }
              }
            ]
          }
        },
        "viewer": {
          "union": {
            "child": [
              {
                "this": {}
              },
              {
                "tupleToUserset": {
                  "tupleset": {
                    "object": "",
                    "relation": "editor"
                  },
                  "computedUserset": {
                    "object": "",
                    "relation": "viewer"
                  }
                }
              }
            ]
          }
        },
        "parent": {
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
          "editor": {
            "directly_related_user_types": [
              {
                "type": "user"
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "viewer": {
            "directly_related_user_types": [
              {
                "type": "user"
              },
              {
                "type": "group",
                "relation": "member"
              }
            ]
          },
          "parent": {
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

## Exemplos de Tuplas

### Configuração de Tenant
```json
{
  "writes": [
    {
      "user": "user:alice",
      "relation": "owner",
      "object": "tenant:acme-corp"
    },
    {
      "user": "user:bob",
      "relation": "admin",
      "object": "tenant:acme-corp"
    },
    {
      "user": "user:charlie",
      "relation": "member",
      "object": "tenant:acme-corp"
    }
  ]
}
```

### Configuração de Vault
```json
{
  "writes": [
    {
      "user": "user:alice",
      "relation": "owner",
      "object": "vault:prod-secrets"
    },
    {
      "user": "tenant:acme-corp",
      "relation": "parent",
      "object": "vault:prod-secrets"
    },
    {
      "user": "group:devops",
      "relation": "admin",
      "object": "vault:prod-secrets"
    }
  ]
}
```

### Configuração de Grupos
```json
{
  "writes": [
    {
      "user": "user:bob",
      "relation": "member",
      "object": "group:devops"
    },
    {
      "user": "user:charlie",
      "relation": "member",
      "object": "group:devops"
    }
  ]
}
```

### Configuração de Secrets
```json
{
  "writes": [
    {
      "user": "user:alice",
      "relation": "owner",
      "object": "secret:database-password"
    },
    {
      "user": "vault:prod-secrets",
      "relation": "parent",
      "object": "secret:database-password"
    },
    {
      "user": "group:devops",
      "relation": "viewer",
      "object": "secret:database-password"
    }
  ]
}
```

### Configuração de Tokens de API
```json
{
  "writes": [
    {
      "user": "user:alice",
      "relation": "owner",
      "object": "token:ci-cd-token-123"
    },
    {
      "user": "vault:prod-secrets",
      "relation": "vault",
      "object": "token:ci-cd-token-123"
    },
    {
      "user": "user:system",
      "relation": "can_read_secrets",
      "object": "token:ci-cd-token-123"
    },
    {
      "user": "user:system",
      "relation": "read_only",
      "object": "token:ci-cd-token-123"
    },
    {
      "user": "user:system",
      "relation": "production",
      "object": "token:ci-cd-token-123"
    }
  ]
}
```

### Configuração de Token de Backup
```json
{
  "writes": [
    {
      "user": "user:admin",
      "relation": "owner",
      "object": "token:backup-token-456"
    },
    {
      "user": "vault:all-vaults",
      "relation": "vault",
      "object": "token:backup-token-456"
    },
    {
      "user": "user:system",
      "relation": "can_read_secrets",
      "object": "token:backup-token-456"
    },
    {
      "user": "user:system",
      "relation": "can_backup",
      "object": "token:backup-token-456"
    }
  ]
}
```

### Configuração de Token de Desenvolvimento
```json
{
  "writes": [
    {
      "user": "user:dev-team",
      "relation": "owner",
      "object": "token:dev-token-789"
    },
    {
      "user": "vault:dev-secrets",
      "relation": "vault",
      "object": "token:dev-token-789"
    },
    {
      "user": "user:system",
      "relation": "can_read_secrets",
      "object": "token:dev-token-789"
    },
    {
      "user": "user:system",
      "relation": "can_write_secrets",
      "object": "token:dev-token-789"
    },
    {
      "user": "user:system",
      "relation": "rate_limited",
      "object": "token:dev-token-789"
    }
  ]
}
```

## Exemplos de Verificação

### Verificar se usuário pode visualizar um secret
```bash
curl -X POST http://localhost:8080/stores/{store_id}/check \
  -H "Content-Type: application/json" \
  -d '{
    "tuple_key": {
      "user": "user:bob",
      "relation": "viewer",
      "object": "secret:database-password"
    }
  }'
```

### Verificar se usuário pode editar um vault
```bash
curl -X POST http://localhost:8080/stores/{store_id}/check \
  -H "Content-Type: application/json" \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "editor",
      "object": "vault:prod-secrets"
    }
  }'
```

### Listar todos os vaults que um usuário pode visualizar
```bash
curl -X POST http://localhost:8080/stores/{store_id}/list-objects \
  -H "Content-Type: application/json" \
  -d '{
    "user": "user:bob",
    "relation": "viewer",
    "type": "vault"
  }'
```

### Verificar se token pode ler um secret
```bash
curl -X POST http://localhost:8080/stores/{store_id}/check \
  -H "Content-Type: application/json" \
  -d '{
    "tuple_key": {
      "user": "token:ci-cd-token-123",
      "relation": "can_read_via_token",
      "object": "vault:prod-secrets"
    }
  }'
```

### Verificar se token pode escrever
```bash
curl -X POST http://localhost:8080/stores/{store_id}/check \
  -H "Content-Type: application/json" \
  -d '{
    "tuple_key": {
      "user": "token:dev-token-789",
      "relation": "can_write_via_token",
      "object": "vault:dev-secrets"
    }
  }'
```

### Verificar se token ainda é válido
```bash
curl -X POST http://localhost:8080/stores/{store_id}/check \
  -H "Content-Type: application/json" \
  -d '{
    "tuple_key": {
      "user": "token:ci-cd-token-123",
      "relation": "is_valid",
      "object": "token:ci-cd-token-123"
    }
  }'
```

### Listar todos os vaults que um token pode acessar
```bash
curl -X POST http://localhost:8080/stores/{store_id}/list-objects \
  -H "Content-Type: application/json" \
  -d '{
    "user": "token:ci-cd-token-123",
    "relation": "can_read_via_token",
    "type": "vault"
  }'
```

## Comandos Úteis

### Criar Store
```bash
curl -X POST http://localhost:8080/stores \
  -H "Content-Type: application/json" \
  -d '{
    "name": "lockari-store"
  }'
```

### Definir Modelo de Autorização
```bash
curl -X POST http://localhost:8080/stores/{store_id}/authorization-models \
  -H "Content-Type: application/json" \
  -d @authorization-model.json
```

### Escrever Tuplas
```bash
curl -X POST http://localhost:8080/stores/{store_id}/write \
  -H "Content-Type: application/json" \
  -d @tuples.json
```

### Verificar Permissões
```bash
curl -X POST http://localhost:8080/stores/{store_id}/check \
  -H "Content-Type: application/json" \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "owner",
      "object": "vault:prod-secrets"
    }
  }'
```

## Integração com Golang

### Exemplo de cliente Go
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/openfga/go-sdk/client"
)

func main() {
    configuration := client.Configuration{
        ApiHost:      "localhost:8080",
        StoreId:      "your-store-id",
        AuthorizationModelId: "your-model-id",
    }
    
    apiClient := client.NewAPIClient(&configuration)
    
    // Verificar permissão
    resp, err := apiClient.OpenFgaApi.Check(context.Background()).Body(client.CheckRequest{
        TupleKey: client.TupleKey{
            User:     "user:alice",
            Relation: "owner",
            Object:   "vault:prod-secrets",
        },
    }).Execute()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Permission allowed: %t\n", resp.Allowed)
}
```

### Exemplo de cliente Go com Token
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/openfga/go-sdk/client"
)

func main() {
    configuration := client.Configuration{
        ApiHost:      "localhost:8080",
        StoreId:      "your-store-id",
        AuthorizationModelId: "your-model-id",
    }
    
    apiClient := client.NewAPIClient(&configuration)
    
    // Verificar se token pode acessar vault
    resp, err := apiClient.OpenFgaApi.Check(context.Background()).Body(client.CheckRequest{
        TupleKey: client.TupleKey{
            User:     "token:ci-cd-token-123",
            Relation: "can_read_via_token",
            Object:   "vault:prod-secrets",
        },
    }).Execute()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Token can access vault: %t\n", resp.Allowed)
    
    // Listar vaults acessíveis pelo token
    listResp, err := apiClient.OpenFgaApi.ListObjects(context.Background()).Body(client.ListObjectsRequest{
        User:     "token:ci-cd-token-123",
        Relation: "can_read_via_token",
        Type:     "vault",
    }).Execute()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Accessible vaults: %v\n", listResp.Objects)
}
```

## Monitoramento

### Métricas do OpenFGA
- `openfga_request_duration_seconds`: Duração das requisições
- `openfga_request_total`: Total de requisições
- `openfga_datastore_query_duration_seconds`: Duração das queries no banco

### Logs Importantes
- Logs de autorização
- Logs de performance
- Logs de erro

Este exemplo serve como base para implementar o sistema de autorização multi-tenant do Lockari usando OpenFGA.
