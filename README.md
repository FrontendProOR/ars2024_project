# API Dokumentacija

Ova dokumentacija opisuje rute i operacije dostupne u našem API-ju za manipulaciju konfiguracijama i konfiguracionim grupama.
## Opšti podaci

- **Swagger dokumentacija:** [http://localhost:8000/swagger](http://localhost:8000/swagger)
- **Consul port:** [http://localhost:8500](http://localhost:8500)
- **Port aplikacije:** [http://localhost:8000](http://localhost:8000)

## Konfiguracije

### Dodavanje konfiguracije

**Metoda:** POST  
**Endpoint:** `/configs`

Dodaje novu konfiguraciju. Potrebno je poslati JSON objekat koji predstavlja konfiguraciju.

### Dobavljanje konfiguracije

**Metoda:** GET  
**Endpoint:** `/configs/{name}/{version}`

Dohvata konfiguraciju po imenu i verziji.

### Brisanje konfiguracije

**Metoda:** DELETE  
**Endpoint:** `/configs/{name}/{version}`

Briše konfiguraciju po imenu i verziji.

## Konfiguracione grupe

### Dodavanje konfiguracione grupe

**Metoda:** POST  
**Endpoint:** `/config-groups`

Dodaje novu konfiguracionu grupu. Potrebno je poslati JSON objekat koji predstavlja grupu.

### Dobavljanje konfiguracione grupe

**Metoda:** GET  
**Endpoint:** `/config-groups/{name}/{version}`

Dohvata konfiguracionu grupu po imenu i verziji.

### Brisanje konfiguracione grupe

**Metoda:** DELETE  
**Endpoint:** `/config-groups/{name}/{version}`

Briše konfiguracionu grupu po imenu i verziji.

### Dodavanje konfiguracije u grupu

**Metoda:** POST  
**Endpoint:** `/config-groups/{name}/{version}/{configName}/{configVersion}`

Dodaje konfiguraciju u određenu grupu.

### Pretraga konfiguracija sa labelama u grupi

**Metoda:** GET  
**Endpoint:** `/config-groups/{name}/{version}/config/search`

Pretražuje konfiguracije u grupi na osnovu labela.

### Dodavanje konfiguracije sa labelom u grupu

**Metoda:** POST  
**Endpoint:** `/config-groups/{name}/{version}/config`

Dodaje konfiguraciju u grupu sa određenim labelom.

### Brisanje konfiguracija sa labelama iz grupe

**Metoda:** DELETE  
**Endpoint:** `/config-groups/{name}/{version}/config/delete`

Briše konfiguracije sa određenim labelama iz grupe.

### Brisanje konfiguracije iz grupe

**Metoda:** DELETE  
**Endpoint:** `/config-groups/{name}/{version}/{configName}/{configVersion}`

Briše konfiguraciju iz grupe.