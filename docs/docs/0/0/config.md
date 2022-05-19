<h1 align="center" style="font-weight: bold">
    Configurations
</h1>

## **Config File**

### Config File Lookup Order of Precedence (CFLOP)

MangaG is a cross-compatible project, which means that it could be ran in different OS. There is however a lack of unity in the standardization on the location of config files in this OSes. And such, I have devised a precedence order for MangaG's config file in different platforms.

The following are the CFLOP for different OSes:

```mermaid
flowchart TD
    A([CFLOP]) --> L[--config argument]
        L --> B{OS?}
        B --> |*nix| C[./MangaG.yml]
            subgraph <br>
                C --> D{"XDG<br>CONFIG<br>HOME<br>(XCH)?"}
                D --> |true| E["${XCH}/MangaG/config.yml"] --> F
                D --> |false| F["~/.config/MangaG/config.yml"]
                F --> G["~/.hyk"]
                G --> H["/etc/xdg/MangaG/config.yml"]
                H --> I["/etc/MangaG/config.yml"]
            end
        B --> |Windows| J[.\MangaG.yml]
            subgraph <br><br>
                J --> K["${boot drive}:\\<br>Users\${username}\<br>AppData\Roaming\MangaG\<br>config.yml"]
            end
```
