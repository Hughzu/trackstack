- login WIP :
    - finish refresh
        2. I can sketch the exact Go refresh contract so backend work is dead simple
        3. I can wire automatic pre-expiry refresh using expiresAt instead of waiting for the first 401
    - audit de gpt sur ce qui est fait
    - fix design AFFREUX de la page de login
    - loading entre login et redirect
    - button logout un peu dégueu, à voir ce que je peux faire en termes de design pour le coin en haut à droite
    - fix design douteux de chacune des pages

    - première fois avec le lazy loading j'ai vraiment un screen flickering chiant quand je change de page, je peux pas réduire ça ? 

- front
    - Code deep dive
        - understand what the fuck is going on in the astro side
            - update docs & focus on DDD and dataflow
            - understand tests
            - design system ? (implement skills to create new pages ?)
            - UI flickering, ... 

- Infra deep dive
    - make a cli (???) to understand what the fuck is going on in bootstrap 
        - idea is to have an AWS cli set up (with root access)
        - github cli set up
        - turso cli set up
        - then create mine that uses them in order to setup everything how I want (maybe create another AWS account for the rest, init the pipelines for terraform and things like that ?)
    
    - modify agents.md with docs about how many clis is available + the doc for my cli

    - AWS / terraform code deep dive



- SRL Module ? (vacancy days, checklist per month with reminders, ...)

- ECS Fargate

- Blog ? 
    - Sell platform engineering profile ? 

- EKS

- Observability

- Other cloud platform ? VPS ? 
