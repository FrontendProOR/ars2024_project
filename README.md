# API Dokumentacija

Ova dokumentacija opisuje rute i operacije dostupne u našem API-ju za manipulaciju konfiguracijama i konfiguracionim grupama.

## Opšti podaci

- **Swagger dokumentacija:** [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html)
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

### Metrike:

1. **Ukupan broj zahteva za prethodna 24 sata**
   - Metrika: `ars2024_project_total_requests`

2. **Broj uspešnih zahteva (status kodova odgovora 2xx, 3xx) za prethodna 24 sata**
   - Metrika: `ars2024_project_successful_requests`

3. **Broj neuspešnih zahteva (status kodova odgovora 4xx, 5xx) za prethodna 24 sata**
   - Metrika: `ars2024_project_failed_requests`

4. **Prosečno vreme izvršavanja zahteva za svaki endpoint**
   - Izraz:
     ```plaintext
     sum(rate(ars2024_project_request_duration_seconds_sum[24h])) / sum(rate(ars2024_project_request_duration_seconds_count[24h]))
     ```
     Ovaj izraz računa prosečno vreme izvršavanja zahteva za svaki endpoint tokom poslednjih 24 sata.

5. **Broj zahteva u jedinici vremena (minut) za svaki endpoint za prethodna 24 sata**
   - Izraz:
     ```plaintext
     sum(rate(ars2024_project_total_requests{endpoint="/config-groups/{name}/{version}"}[24h])) / 24 / 60
     ```
     Ovaj izraz računa broj zahteva po minutu za određeni endpoint `/config-groups/{name}/{version}` tokom prethodnih 24 sata.
