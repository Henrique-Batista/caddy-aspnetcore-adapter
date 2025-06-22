# aspnetadapter

Módulo para o Caddy Server que adiciona compatibilidade com aplicações ASP.NET e ASP.NET Core. Este módulo oferece:

* **Adapter de web.config** para importar regras de rewrite, redirect e documentos padrões;
* **Middleware de compatibilidade** para tratar extensões como `.aspx`, `.ashx`, `.axd`, etc.;
* **Gerenciador de processo** para iniciar o runtime do ASP.NET Core automaticamente sem modificar o `web.config`;
* **Integração com reverse\_proxy** para roteamento de requisições ao backend ASP.NET Core.

## Recursos Suportados

* Rewrite rules (`<rewrite>`)
* Redirecionamentos (`<httpRedirect>`)
* Documentos padrões (`<defaultDocument>`)
* Início automático de processos `.NET` com argumentos personalizados

## Requisitos

* Go 1.20 ou superior
* Caddy v2.7.0 ou superior
* Runtime do .NET instalado no sistema (`dotnet`)

---

## Instalação

Você pode compilar o módulo com `xcaddy`:

```bash
xcaddy build --with github.com/seuusuario/aspnetadapter
```

---

## Uso Básico

### Adaptar `web.config`

Se quiser converter `web.config` para um `Caddyfile`, use:

```bash
caddy adapt --config web.config --adapter webconfig
```

Isso gerará um Caddyfile contendo diretivas `rewrite`, `redir`, `try_files` e `file_server` equivalentes.

---

### Iniciar ASP.NET Core + Caddy

Você pode usar `process_manager` diretamente no `Caddyfile`:

```caddyfile
{
  order middleware.aspnet_compat before file_server
}

:8080 {
  process_manager {
    process_path "dotnet"
    process_args "MyApp.dll --urls http://127.0.0.1:5000"
    restart_delay 5s
  }

  route {
    middleware.aspnet_compat
    reverse_proxy 127.0.0.1:5000
    file_server
  }
}
```

Ou iniciar manualmente via linha de comando:

```bash
caddy run \
  --process_path "/usr/bin/dotnet" \
  --process_args "MyApp.dll --urls http://127.0.0.1:5000" \
  --restart_delay 5s
```

---

## Exemplo Completo

### web.config

```xml
<configuration>
  <system.webServer>
    <rewrite>
      <rules>
        <rule name="CleanUrls">
          <match url="^([_0-9a-z-]+)$" />
          <action type="Rewrite" url="/{R:1}.html" />
        </rule>
      </rules>
    </rewrite>
    <httpRedirect enabled="true" destination="https://example.com/" />
    <defaultDocument>
      <files>
        <add value="index.html" />
        <add value="default.html" />
      </files>
    </defaultDocument>
  </system.webServer>
</configuration>
```

### Adaptando para Caddy

```bash
caddy adapt --config web.config --adapter webconfig
```

---

## Testes

Inclua testes para:

* Parsing de `web.config`
* Execução controlada do processo
* Middleware de compatibilidade

Exemplo de teste simples:

```go
func TestStartProcess(t *testing.T) {
  cfg := Config{
    ProcessPath: "sleep",
    Arguments:   "1",
    RestartDelay: 1 * time.Second,
  }
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()
  err := StartProcess(ctx, cfg)
  if err != nil {
    t.Fatal(err)
  }
}
```

---

## Licença

MIT

---

## Contribuindo

Pull requests e sugestões são bem-vindos! Abra uma issue se quiser propor uma nova funcionalidade ou relatar bugs.
